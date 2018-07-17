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

	"github.com/iomz/gosstrak/tdt"
)

// SplayTree struct
type SplayTree struct {
	root    *SplayTreeNode
	tdtCore *tdt.Core
}

type SplayTreeNode struct {
	reportURI    string
	filterObject *FilterObject
	matchNext    *SplayTreeNode
	mismatchNext *SplayTreeNode
}

// AddSubscription adds a set of subscriptions if not exists yet
func (spt *SplayTree) AddSubscription(sub ByteSubscriptions) {
	for _, fs := range sub.Keys() {
		spt.root.add(fs, sub[fs].ReportURI)
	}
}

// DeleteSubscription deletes a set of subscriptions if already exist
func (spt *SplayTree) DeleteSubscription(sub ByteSubscriptions) {
	for _, fs := range sub.Keys() {
		spt.root.delete(fs, sub[fs].ReportURI)
	}
}

// Dump returs a string representation of the PatriciaTrie
func (spt *SplayTree) Dump() string {
	writer := &bytes.Buffer{}
	spt.root.print(writer, 0)
	return writer.String()
}

// MarshalBinary overwrites the marshaller in gob encoding *SplayTree
func (spt *SplayTree) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.SplayTree")

	// Encode SplayTreeNode
	enc.Encode(spt.root)

	return buf.Bytes(), err
}

func (spt *SplayTree) Name() string {
	return "SplayTree"
}

// Search returns a slice of reportURI
func (spt *SplayTree) Search(id []byte) []string {
	return spt.root.splaySearch(spt, nil, id)
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (spt *SplayTree) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.SplayTree" {
		return errors.New("Wrong Filtering Engine: " + typeOfEngine)
	}

	// Decode SplayTreeNode
	err = dec.Decode(&spt.root)

	spt.tdtCore = tdt.NewCore()

	return
}

// MarshalBinary overwrites the marshaller in gob encoding *SplayTreeNode
func (sptn *SplayTreeNode) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// ReportURI
	enc.Encode(sptn.reportURI)

	// Filter
	hasFilter := sptn.filterObject != nil
	enc.Encode(hasFilter)
	if hasFilter {
		err = enc.Encode(sptn.filterObject)
	}

	// matchNext
	hasMatchNext := sptn.matchNext != nil
	enc.Encode(hasMatchNext)
	if hasMatchNext {
		err = enc.Encode(sptn.matchNext)
	}

	// mismatchNext
	hasMismatchNext := sptn.mismatchNext != nil
	enc.Encode(hasMismatchNext)
	if hasMismatchNext {
		err = enc.Encode(sptn.mismatchNext)
	}

	return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *SplayTreeNode
func (sptn *SplayTreeNode) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// reportURI
	if err = dec.Decode(&sptn.reportURI); err != nil {
		return
	}

	// FilterObject
	var hasFilterObject bool
	if err = dec.Decode(&hasFilterObject); err != nil {
		return
	}
	if hasFilterObject {
		err = dec.Decode(&sptn.filterObject)
	}

	// matchNext
	var hasMatchNext bool
	if err = dec.Decode(&hasMatchNext); err != nil {
		return
	}
	if hasMatchNext {
		err = dec.Decode(&sptn.matchNext)
	} else {
		sptn.matchNext = nil
	}

	// mismatchNext
	var hasMismatchNext bool
	if err = dec.Decode(&hasMismatchNext); err != nil {
		return
	}
	if hasMismatchNext {
		err = dec.Decode(&sptn.mismatchNext)
	} else {
		sptn.mismatchNext = nil
	}

	return
}

// Internal helper methods -----------------------------------------------------

// add a set of subscriptions if not exists yet
func (sptn *SplayTreeNode) add(fs string, reportURI string) {
	if strings.HasPrefix(fs, sptn.filterObject.String) { // fs \in sptn.FilterObject.String
		if fs == sptn.filterObject.String { // the identical filter
			// if the reportURI is different, update it
			if sptn.reportURI != reportURI {
				sptn.reportURI = reportURI
			}
		} else { // it's a subset (matching branch) of the current node, and there's no matchNext node
			if sptn.matchNext == nil {
				sptn.matchNext = &SplayTreeNode{}
				sptn.matchNext.filterObject = NewFilter(fs[sptn.filterObject.Size:], sptn.filterObject.Offset+sptn.filterObject.Size)
				sptn.matchNext.reportURI = reportURI
			} else { // if there's already matchNext node
				sptn.matchNext.add(fs[sptn.filterObject.Size:], reportURI)
			}
		}
	} else { // doesn't match with the current node, traverse the mismatchNext node
		if sptn.mismatchNext == nil { // there's no mismatchNext node
			sptn.mismatchNext = &SplayTreeNode{}
			sptn.mismatchNext.filterObject = NewFilter(fs, sptn.filterObject.Offset)
			sptn.mismatchNext.reportURI = reportURI
		} else { // if there's already mismatchNext node
			sptn.mismatchNext.add(fs, reportURI)
		}
	}
	return
}

func (sptn *SplayTreeNode) build(sub ByteSubscriptions) *SplayTreeNode {
	current := sptn
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
	return sptn
}

// delete a set of subscriptions if not exists yet
func (sptn *SplayTreeNode) delete(fs string, reportURI string) {
	if strings.HasPrefix(fs, sptn.filterObject.String) { // fs \in sptn.FilterObject.String
		if fs == sptn.filterObject.String { // this node is to delete
			if sptn.matchNext == nil && sptn.mismatchNext == nil { // something wrong
			} else if sptn.matchNext != nil { // if there is subset, keep the node as an aggregation node
				if sptn.matchNext.mismatchNext != nil {
					sptn.reportURI = ""
				} else { // if none other mismatch branch, concatenate the matchNext with to-be-deleted node
					sptn.filterObject = NewFilter(fs+sptn.matchNext.filterObject.String, sptn.filterObject.Offset)
					sptn.reportURI = sptn.matchNext.reportURI
					sptn.matchNext = sptn.matchNext.matchNext
				}
			} else if sptn.mismatchNext != nil { // replace this node with mismatchNext
				sptn.filterObject = sptn.mismatchNext.filterObject
				sptn.reportURI = sptn.mismatchNext.reportURI
				sptn.matchNext = sptn.mismatchNext.matchNext
				sptn.mismatchNext = sptn.mismatchNext.mismatchNext
			}
		} else { // it's a subset (matching branch) of the current node, and there's no matchNext node
			if sptn.matchNext != nil { // if there's a matchNext node
				if fs[sptn.filterObject.Size:] == sptn.matchNext.filterObject.String &&
					sptn.matchNext.matchNext == nil && sptn.matchNext.mismatchNext == nil { // the matchNext is to delete
					sptn.matchNext = nil
				} else {
					sptn.matchNext.delete(fs[sptn.filterObject.Size:], reportURI)
				}
			}
		}
	} else { // doesn't match with the current node, traverse the mismatchNext node
		if sptn.mismatchNext != nil { // there's a mismatchNext node
			if fs == sptn.mismatchNext.filterObject.String &&
				sptn.mismatchNext.matchNext == nil && sptn.mismatchNext.mismatchNext == nil {
				sptn.mismatchNext = nil
			} else {
				sptn.mismatchNext.delete(fs, reportURI)
			}
		}
	}
	return
}

func (sptn *SplayTreeNode) equal(want *SplayTreeNode) (ok bool, got *SplayTreeNode, wanted *SplayTreeNode) {
	if sptn.reportURI != want.reportURI ||
		!reflect.DeepEqual(sptn.filterObject, want.filterObject) {
		return false, sptn, want
	}
	if want.matchNext != nil {
		if sptn.matchNext == nil {
			return false, nil, want.matchNext
		}
		res, cgot, cwanted := sptn.matchNext.equal(want.matchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	if want.mismatchNext != nil {
		if sptn.mismatchNext == nil {
			return false, nil, want.mismatchNext
		}
		res, cgot, cwanted := sptn.mismatchNext.equal(want.mismatchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	return true, nil, nil
}

func (sptn *SplayTreeNode) print(writer io.Writer, indent int) {
	var n string
	if len(sptn.reportURI) != 0 {
		n = "-> " + sptn.reportURI
	}
	fmt.Fprintf(writer, "--%s %s\n", sptn.filterObject.ToString(), n)
	if sptn.matchNext != nil {
		fmt.Fprintf(writer, "%sok", strings.Repeat(" ", indent+2))
		sptn.matchNext.print(writer, indent+2)
	}
	if sptn.mismatchNext != nil {
		fmt.Fprintf(writer, "%sng", strings.Repeat(" ", indent+2))
		sptn.mismatchNext.print(writer, indent+2)
	}
}

func (sptn *SplayTreeNode) splaySearch(spt *SplayTree, parent *SplayTreeNode, id []byte) []string {
	matches := []string{}
	if sptn.filterObject.Match(id) {
		matches = append(matches, sptn.reportURI)
		if sptn.matchNext != nil {
			// Do Search & Splay in the subsets
			matches = append(matches, sptn.matchNext.splaySearch(spt, sptn, id)...)
		}
		// Do Splay
		// 0. Check if this is the root node, do nothing if so
		if parent != nil {
			// 1. Remove this node by connecting parent to the next mismatchNext node
			parent.mismatchNext = sptn.mismatchNext
			// 2. Insert self to root
			sptn.mismatchNext = spt.root
			spt.root = sptn
		}
		return matches
	}
	if sptn.mismatchNext != nil {
		return sptn.mismatchNext.splaySearch(spt, sptn, id)
	}
	return matches
}

// NewSplayTree builds SplayTree from ByteSubscriptions
// returns the pointer to the node node
func NewSplayTree(sub ByteSubscriptions) Engine {
	// make subsets to the child subscriptions of the corresponding parents
	sub.linkSubset()

	spt := &SplayTree{}
	spt.root = &SplayTreeNode{}
	spt.root = spt.root.build(sub)
	spt.tdtCore = tdt.NewCore()

	return spt
}
