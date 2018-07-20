// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/iomz/go-llrp"
	"github.com/iomz/gosstrak/tdt"
)

// SplayTree struct
type SplayTree struct {
	root    *SplayTreeNode
	tdtCore *tdt.Core
}

// SplayTreeNode is a node for SplayTree
type SplayTreeNode struct {
	reportURI    string
	filterObject *FilterObject
	matchNext    *SplayTreeNode
	mismatchNext *SplayTreeNode
}

// AddSubscription adds a set of subscriptions if not exists yet
func (st *SplayTree) AddSubscription(sub Subscriptions) {
	bsub := sub.ToByteSubscriptions()
	for _, fs := range bsub.Keys() {
		st.root.add(fs, bsub[fs].ReportURI)
	}
}

// DeleteSubscription deletes a set of subscriptions if already exist
func (st *SplayTree) DeleteSubscription(sub Subscriptions) {
	bsub := sub.ToByteSubscriptions()
	for _, fs := range bsub.Keys() {
		st.root.delete(fs, bsub[fs].ReportURI)
	}
}

// Dump returs a string representation of the PatriciaTrie
func (st *SplayTree) Dump() string {
	writer := &bytes.Buffer{}
	st.root.print(writer, 0)
	return writer.String()
}

// MarshalBinary overwrites the marshaller in gob encoding *SplayTree
func (st *SplayTree) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.SplayTree")

	// Encode SplayTreeNode
	enc.Encode(st.root)

	return buf.Bytes(), err
}

// Name returns the name of this engine
func (st *SplayTree) Name() string {
	return "SplayTree"
}

// Search returns a pureIdentity of the llrp.ReadEvent if found any subscription without err
func (st *SplayTree) Search(re llrp.ReadEvent) (pureIdentity string, reportURIs []string, err error) {
	reportURIs = st.root.splaySearch(st, nil, re.ID)
	if len(reportURIs) == 0 {
		return pureIdentity, reportURIs, fmt.Errorf("no match found for %v", re.ID)
	}
	pureIdentity, err = st.tdtCore.Translate(re.PC, re.ID)
	return
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (st *SplayTree) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.SplayTree" {
		return errors.New("Wrong Filtering Engine: " + typeOfEngine)
	}

	// Decode SplayTreeNode
	err = dec.Decode(&st.root)

	// tdt.Core
	st.tdtCore = tdt.NewCore()

	return
}

// MarshalBinary overwrites the marshaller in gob encoding *SplayTreeNode
func (stn *SplayTreeNode) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// ReportURI
	enc.Encode(stn.reportURI)

	// Filter
	hasFilter := stn.filterObject != nil
	enc.Encode(hasFilter)
	if hasFilter {
		err = enc.Encode(stn.filterObject)
	}

	// matchNext
	hasMatchNext := stn.matchNext != nil
	enc.Encode(hasMatchNext)
	if hasMatchNext {
		err = enc.Encode(stn.matchNext)
	}

	// mismatchNext
	hasMismatchNext := stn.mismatchNext != nil
	enc.Encode(hasMismatchNext)
	if hasMismatchNext {
		err = enc.Encode(stn.mismatchNext)
	}

	return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *SplayTreeNode
func (stn *SplayTreeNode) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// reportURI
	if err = dec.Decode(&stn.reportURI); err != nil {
		return
	}

	// FilterObject
	var hasFilterObject bool
	if err = dec.Decode(&hasFilterObject); err != nil {
		return
	}
	if hasFilterObject {
		err = dec.Decode(&stn.filterObject)
	}

	// matchNext
	var hasMatchNext bool
	if err = dec.Decode(&hasMatchNext); err != nil {
		return
	}
	if hasMatchNext {
		err = dec.Decode(&stn.matchNext)
	} else {
		stn.matchNext = nil
	}

	// mismatchNext
	var hasMismatchNext bool
	if err = dec.Decode(&hasMismatchNext); err != nil {
		return
	}
	if hasMismatchNext {
		err = dec.Decode(&stn.mismatchNext)
	} else {
		stn.mismatchNext = nil
	}

	return
}

// Internal helper methods -----------------------------------------------------

// add a set of subscriptions if not exists yet
func (stn *SplayTreeNode) add(fs string, reportURI string) {
	if strings.HasPrefix(fs, stn.filterObject.String) { // fs \in stn.FilterObject.String
		if fs == stn.filterObject.String { // the identical filter
			// if the reportURI is different, update it
			if stn.reportURI != reportURI {
				stn.reportURI = reportURI
			}
		} else { // it's a subset (matching branch) of the current node, and there's no matchNext node
			if stn.matchNext == nil {
				stn.matchNext = &SplayTreeNode{}
				stn.matchNext.filterObject = NewFilter(fs[stn.filterObject.Size:], stn.filterObject.Offset+stn.filterObject.Size)
				stn.matchNext.reportURI = reportURI
			} else { // if there's already matchNext node
				stn.matchNext.add(fs[stn.filterObject.Size:], reportURI)
			}
		}
	} else { // doesn't match with the current node, traverse the mismatchNext node
		if stn.mismatchNext == nil { // there's no mismatchNext node
			stn.mismatchNext = &SplayTreeNode{}
			stn.mismatchNext.filterObject = NewFilter(fs, stn.filterObject.Offset)
			stn.mismatchNext.reportURI = reportURI
		} else { // if there's already mismatchNext node
			stn.mismatchNext.add(fs, reportURI)
		}
	}
	return
}

func (stn *SplayTreeNode) build(sub ByteSubscriptions) *SplayTreeNode {
	current := stn
	subscriptionSize := len(sub.Keys())
	for i, fs := range sub.Keys() {
		current.filterObject = NewFilter(fs, sub[fs].Offset)
		current.reportURI = sub[fs].ReportURI
		// if this node has subset
		if len(sub[fs].Subset) != 0 {
			matchNext := &SplayTreeNode{}
			current.matchNext = matchNext.build(sub[fs].Subset)
		} else {
			current.matchNext = nil
		}
		if i+1 < subscriptionSize {
			current.mismatchNext = &SplayTreeNode{}
			current = current.mismatchNext
		} else {
			current.mismatchNext = nil
		}
	}
	return stn
}

// delete a set of subscriptions if not exists yet
func (stn *SplayTreeNode) delete(fs string, reportURI string) {
	if strings.HasPrefix(fs, stn.filterObject.String) { // fs \in stn.FilterObject.String
		if fs == stn.filterObject.String { // this node is to delete
			if stn.matchNext == nil && stn.mismatchNext == nil { // something wrong
			} else if stn.matchNext != nil { // if there is subset, keep the node as an aggregation node
				if stn.matchNext.mismatchNext != nil {
					stn.reportURI = ""
				} else { // if none other mismatch branch, concatenate the matchNext with to-be-deleted node
					stn.filterObject = NewFilter(fs+stn.matchNext.filterObject.String, stn.filterObject.Offset)
					stn.reportURI = stn.matchNext.reportURI
					stn.matchNext = stn.matchNext.matchNext
				}
			} else if stn.mismatchNext != nil { // replace this node with mismatchNext
				stn.filterObject = stn.mismatchNext.filterObject
				stn.reportURI = stn.mismatchNext.reportURI
				stn.matchNext = stn.mismatchNext.matchNext
				stn.mismatchNext = stn.mismatchNext.mismatchNext
			}
		} else { // it's a subset (matching branch) of the current node, and there's no matchNext node
			if stn.matchNext != nil { // if there's a matchNext node
				if fs[stn.filterObject.Size:] == stn.matchNext.filterObject.String &&
					stn.matchNext.matchNext == nil && stn.matchNext.mismatchNext == nil { // the matchNext is to delete
					stn.matchNext = nil
				} else {
					stn.matchNext.delete(fs[stn.filterObject.Size:], reportURI)
				}
			}
		}
	} else { // doesn't match with the current node, traverse the mismatchNext node
		if stn.mismatchNext != nil { // there's a mismatchNext node
			if fs == stn.mismatchNext.filterObject.String &&
				stn.mismatchNext.matchNext == nil && stn.mismatchNext.mismatchNext == nil {
				stn.mismatchNext = nil
			} else {
				stn.mismatchNext.delete(fs, reportURI)
			}
		}
	}
	return
}

func (stn *SplayTreeNode) equal(want *SplayTreeNode) (ok bool, got *SplayTreeNode, wanted *SplayTreeNode) {
	if stn.reportURI != want.reportURI ||
		!reflect.DeepEqual(stn.filterObject, want.filterObject) {
		return false, stn, want
	}
	if want.matchNext != nil {
		if stn.matchNext == nil {
			return false, nil, want.matchNext
		}
		res, cgot, cwanted := stn.matchNext.equal(want.matchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	if want.mismatchNext != nil {
		if stn.mismatchNext == nil {
			return false, nil, want.mismatchNext
		}
		res, cgot, cwanted := stn.mismatchNext.equal(want.mismatchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	return true, nil, nil
}

func (stn *SplayTreeNode) print(writer io.Writer, indent int) {
	var n string
	if len(stn.reportURI) != 0 {
		n = "-> " + stn.reportURI
	}
	fmt.Fprintf(writer, "--%s %s\n", stn.filterObject.ToString(), n)
	if stn.matchNext != nil {
		fmt.Fprintf(writer, "%sok", strings.Repeat(" ", indent+2))
		stn.matchNext.print(writer, indent+2)
	}
	if stn.mismatchNext != nil {
		fmt.Fprintf(writer, "%sng", strings.Repeat(" ", indent+2))
		stn.mismatchNext.print(writer, indent+2)
	}
}

func (stn *SplayTreeNode) splaySearch(st *SplayTree, parent *SplayTreeNode, id []byte) []string {
	matches := []string{}
	if stn.filterObject.Match(id) {
		matches = append(matches, stn.reportURI)
		if stn.matchNext != nil {
			// Do Search & Splay in the subsets
			matches = append(matches, stn.matchNext.splaySearch(st, nil, id)...)
		}
		// Do Splay
		// 0. Check if this is the root node, do nothing if so
		if parent != nil {
			// 1. Remove this node by connecting parent to the next mismatchNext node
			parent.mismatchNext = stn.mismatchNext
			// 2. Insert self to root
			stn.mismatchNext = st.root
			st.root = stn
		}
		return matches
	}
	if stn.mismatchNext != nil {
		return stn.mismatchNext.splaySearch(st, stn, id)
	}
	return matches
}

// NewSplayTree builds SplayTree from ByteSubscriptions
// returns the pointer to the node node
func NewSplayTree(sub Subscriptions) Engine {
	st := &SplayTree{}

	// preprocess the subscriptions
	bsub := sub.ToByteSubscriptions()
	// make subsets to the child subscriptions of the corresponding parents
	bsub.linkSubset()

	// build SplayTree
	st.root = &SplayTreeNode{}
	st.root = st.root.build(bsub)

	// initialize the tdt.Core
	st.tdtCore = tdt.NewCore()

	return st
}
