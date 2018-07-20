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

// PatriciaTrie struct
type PatriciaTrie struct {
	root    *PatriciaTrieNode
	tdtCore *tdt.Core
}

// PatriciaTrieNode is a node for PatriciaTrie
type PatriciaTrieNode struct {
	reportURI    string
	filterObject *FilterObject
	one          *PatriciaTrieNode
	zero         *PatriciaTrieNode
}

// AddSubscription adds a set of subscriptions if not exists yet
func (pt *PatriciaTrie) AddSubscription(sub Subscriptions) {
	bsub := sub.ToByteSubscriptions()
	for _, fs := range bsub.Keys() {
		pt.root.add(fs, bsub[fs].ReportURI)
	}
}

// DeleteSubscription deletes a set of subscriptions if already exist
func (pt *PatriciaTrie) DeleteSubscription(sub Subscriptions) {
	bsub := sub.ToByteSubscriptions()
	for _, fs := range bsub.Keys() {
		pt.root.delete(fs, bsub[fs].ReportURI)
	}
}

// Dump returs a string representation of the PatriciaTrie
func (pt *PatriciaTrie) Dump() string {
	writer := &bytes.Buffer{}
	pt.root.print(writer, 0)
	return writer.String()
}

// MarshalBinary overwrites the marshaller in gob encoding *PatriciaTrie
func (pt *PatriciaTrie) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.PatriciaTrie")

	// Encode PatriciaTrieNode
	enc.Encode(pt.root)

	return buf.Bytes(), err
}

// Name returs the name of this engine type
func (pt *PatriciaTrie) Name() string {
	return "PatriciaTrie"
}

// Search returns a pureIdentity of the llrp.ReadEvent if found any subscription without err
func (pt *PatriciaTrie) Search(re llrp.ReadEvent) (pureIdentity string, reportURIs []string, err error) {
	reportURIs = pt.root.search(re.ID)
	if len(reportURIs) == 0 {
		return pureIdentity, reportURIs, fmt.Errorf("no match found for %v", re.ID)
	}
	pureIdentity, err = pt.tdtCore.Translate(re.PC, re.ID)
	return
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (pt *PatriciaTrie) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.PatriciaTrie" {
		return errors.New("Wrong Filtering Engine: " + typeOfEngine)
	}

	// Decode PatriciaTrieNode
	err = dec.Decode(&pt.root)

	// tdt.Core
	pt.tdtCore = tdt.NewCore()

	return
}

// MarshalBinary overwrites the marshaller in gob encoding *PatriciaTrie
func (ptn *PatriciaTrieNode) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// reportURI
	enc.Encode(ptn.reportURI)

	// Filter
	hasFilter := ptn.filterObject != nil
	enc.Encode(hasFilter)
	if hasFilter {
		err = enc.Encode(ptn.filterObject)
	}

	// One
	hasOne := ptn.one != nil
	enc.Encode(hasOne)
	if hasOne {
		err = enc.Encode(ptn.one)
	}

	// Zero
	hasZero := ptn.zero != nil
	enc.Encode(hasZero)
	if hasZero {
		err = enc.Encode(ptn.zero)
	}

	return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (ptn *PatriciaTrieNode) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// reportURIs
	if err = dec.Decode(&ptn.reportURI); err != nil {
		return
	}

	// Filter
	var hasFilter bool
	if err = dec.Decode(&hasFilter); err != nil {
		return
	}
	if hasFilter {
		err = dec.Decode(&ptn.filterObject)
	} else {
		ptn.filterObject = nil
	}

	// One
	var hasOne bool
	if err = dec.Decode(&hasOne); err != nil {
		return
	}
	if hasOne {
		err = dec.Decode(&ptn.one)
	} else {
		ptn.one = nil
	}

	// Zero
	var hasZero bool
	if err = dec.Decode(&hasZero); err != nil {
		return
	}
	if hasOne {
		err = dec.Decode(&ptn.zero)
	} else {
		ptn.zero = nil
	}

	return
}

// Internal helper methods -----------------------------------------------------

// add a subscription if not exist yet
func (ptn *PatriciaTrieNode) add(fs string, reportURI string) {
	if strings.HasPrefix(fs, ptn.filterObject.String) { // fs \in pt.FilterObject.String
		if fs == ptn.filterObject.String { // the identical filter
			// if the reportURI is different, update it
			if ptn.reportURI != reportURI {
				ptn.reportURI = reportURI
			}
			return //end
		}
		//} else if len(fs) < pt.filterObject.Size { // Needs a reconstruction
	} else {
		newCommonPrefix := lcp([]string{fs, ptn.filterObject.String})
		ncpLength := len(newCommonPrefix)
		newNode := &PatriciaTrieNode{}
		newNode.filterObject = NewFilter(ptn.filterObject.String[ncpLength:], ptn.filterObject.Offset+ncpLength)
		newNode.one = ptn.one
		newNode.zero = ptn.zero
		newNode.reportURI = ptn.reportURI
		ptn.reportURI = ""
		currentOffset := ptn.filterObject.Offset
		ptn.filterObject = NewFilter(newCommonPrefix, currentOffset)
		switch fs[ncpLength] {
		case '1':
			ptn.zero = newNode
			ptn.one = &PatriciaTrieNode{}
			ptn.one.filterObject = NewFilter(fs[ncpLength:], currentOffset+ncpLength)
			ptn.one.reportURI = reportURI
		case '0':
			ptn.one = newNode
			ptn.zero = &PatriciaTrieNode{}
			ptn.zero.filterObject = NewFilter(fs[ncpLength:], currentOffset+ncpLength)
			ptn.zero.reportURI = reportURI
		}
		return //end
	}

	// If there's remainder
	if len(fs) > ptn.filterObject.Size {
		switch fs[ptn.filterObject.Size] {
		case '1':
			if ptn.one == nil {
				ptn.one = &PatriciaTrieNode{}
				ptn.one.filterObject = NewFilter(fs[ptn.filterObject.Size:], ptn.filterObject.Offset+ptn.filterObject.Size)
				ptn.one.reportURI = reportURI
				return //end
			}
			ptn.one.add(fs[ptn.filterObject.Size:], reportURI)
		case '0':
			if ptn.zero == nil {
				ptn.zero = &PatriciaTrieNode{}
				ptn.zero.filterObject = NewFilter(fs[ptn.filterObject.Size:], ptn.filterObject.Offset+ptn.filterObject.Size)
				ptn.zero.reportURI = reportURI
				return //end
			}
			ptn.zero.add(fs[ptn.filterObject.Size:], reportURI)
		}
	}
}

// build PatriciaTrieNode recursively
func (ptn *PatriciaTrieNode) build(prefix string, bsub ByteSubscriptions) {
	onePrefixBranch := ""
	zeroPrefixBranch := ""
	fks := bsub.Keys()
	for _, fk := range fks {
		// if the prefix is already longer than the testee
		if len(fk) < len(prefix) {
			continue
		}
		// ignore the testee without the prefix
		if !strings.HasPrefix(fk, prefix) {
			continue
		}
		p := fk[len(prefix):]
		// ignore if no remainder
		if len(p) == 0 {
			continue
		}
		// if the remainder starts with 1
		if strings.HasPrefix(p, "1") {
			if len(onePrefixBranch) == 0 {
				onePrefixBranch = p
			} else {
				onePrefixBranch = lcp([]string{p, onePrefixBranch})
			}
			// if the remainder starts with 0
		} else if strings.HasPrefix(p, "0") {
			if len(zeroPrefixBranch) == 0 {
				zeroPrefixBranch = p
			} else {
				zeroPrefixBranch = lcp([]string{p, zeroPrefixBranch})
			}
		}
	}
	cumulativePrefix := ""
	// if there's a branch starts with 1
	if len(onePrefixBranch) != 0 {
		ptn.one = &PatriciaTrieNode{}
		ptn.one.filterObject = NewFilter(onePrefixBranch, len(prefix))
		cumulativePrefix = prefix + onePrefixBranch
		// check if the prefix matches whole filter
		if _, ok := bsub[cumulativePrefix]; ok {
			ptn.one.reportURI = bsub[cumulativePrefix].ReportURI
		}
		ptn.one.build(cumulativePrefix, bsub)
	}
	// if there's a branch starts with 0
	if len(zeroPrefixBranch) != 0 {
		ptn.zero = &PatriciaTrieNode{}
		ptn.zero.filterObject = NewFilter(zeroPrefixBranch, len(prefix))
		cumulativePrefix = prefix + zeroPrefixBranch
		// check if the prefix matches whole filter
		if _, ok := bsub[cumulativePrefix]; ok {
			ptn.zero.reportURI = bsub[cumulativePrefix].ReportURI
		}
		ptn.zero.build(cumulativePrefix, bsub)
	}
}

// delete a subscription if already exists
func (ptn *PatriciaTrieNode) delete(fs string, reportURI string) {
	// No such filter exist in the trie
	if !strings.HasPrefix(fs, ptn.filterObject.String) {
		return
	}

	// This is the filter to delete
	if fs == ptn.filterObject.String {
		if ptn.one != nil && ptn.zero != nil { // node in the middle
			ptn.reportURI = ""
		} else if ptn.one != nil { // has only one node
			newFilter := NewFilter(ptn.filterObject.String+ptn.one.filterObject.String, ptn.filterObject.Offset)
			ptn.filterObject = newFilter
			ptn.reportURI = ptn.one.reportURI
			ptn.zero = ptn.one.zero
			ptn.one = ptn.one.one
		} else if ptn.zero != nil { // has only zero node
			newFilter := NewFilter(ptn.filterObject.String+ptn.zero.filterObject.String, ptn.filterObject.Offset)
			ptn.filterObject = newFilter
			ptn.reportURI = ptn.zero.reportURI
			ptn.one = ptn.zero.one
			ptn.zero = ptn.zero.zero
		}
		return //end
	}

	// If there's remainder
	if len(fs) > ptn.filterObject.Size {
		switch fs[ptn.filterObject.Size] {
		case '1':
			if ptn.one != nil {
				if fs[ptn.filterObject.Size:] == ptn.one.filterObject.String &&
					ptn.one.one == nil && ptn.one.zero == nil {
					ptn.one = nil
				} else {
					ptn.one.delete(fs[ptn.filterObject.Size:], reportURI)
				}
			}
		case '0':
			if ptn.zero != nil {
				if fs[ptn.filterObject.Size:] == ptn.zero.filterObject.String &&
					ptn.zero.one == nil && ptn.zero.zero == nil {
					ptn.zero = nil
				} else {
					ptn.zero.delete(fs[ptn.filterObject.Size:], reportURI)
				}
			}
		}
	}
	return
}

func (ptn *PatriciaTrieNode) equal(want *PatriciaTrieNode) (ok bool, got *PatriciaTrieNode, wanted *PatriciaTrieNode) {
	if ptn.reportURI != want.reportURI ||
		!reflect.DeepEqual(ptn.filterObject, want.filterObject) {
		return false, ptn, want
	}
	if want.one != nil {
		if ptn.one == nil {
			return false, nil, want.one
		}
		res, cgot, cwanted := ptn.one.equal(want.one)
		if !res {
			return res, cgot, cwanted
		}
	}
	if want.zero != nil {
		if ptn.zero == nil {
			return false, nil, want.zero
		}
		res, cgot, cwanted := ptn.zero.equal(want.zero)
		if !res {
			return res, cgot, cwanted
		}
	}
	return true, nil, nil
}

func (ptn *PatriciaTrieNode) print(writer io.Writer, indent int) {
	var n string
	if len(ptn.reportURI) != 0 {
		n = "-> " + ptn.reportURI
	}
	fmt.Fprintf(writer, "%s--%s %s\n", strings.Repeat(" ", indent), ptn.filterObject.ToString(), n)
	if ptn.one != nil {
		ptn.one.print(writer, indent+2)
	}
	if ptn.zero != nil {
		ptn.zero.print(writer, indent+2)
	}
}

func (ptn *PatriciaTrieNode) search(id []byte) (reportURIs []string) {
	// if not match, return empty slice immediately
	if !ptn.filterObject.Match(id) {
		return
	}

	// if the id matched with this node, return reportURI
	if len(ptn.reportURI) != 0 {
		reportURIs = append(reportURIs, ptn.reportURI)
	}

	// Determine next filter
	nextBitOffset := ptn.filterObject.Offset + ptn.filterObject.Size
	nb, err := getNextBit(id, nextBitOffset)
	if err != nil {
		panic(err)
	}
	if nb == '1' && ptn.one != nil {
		reportURIs = append(reportURIs, ptn.one.search(id)...)
	} else if nb == '0' && ptn.zero != nil {
		reportURIs = append(reportURIs, ptn.zero.search(id)...)
	}
	return
}

// NewPatriciaTrie builds PatriciaTrie from filter.ByteSubscriptions
// returns the pointer to the node
func NewPatriciaTrie(sub Subscriptions) Engine {
	pt := &PatriciaTrie{}

	// preprocess the subscriptions
	bsub := sub.ToByteSubscriptions()

	// build PatriciaTrie
	p1 := lcp(bsub.Keys())
	if len(p1) == 0 {
		// do something if there's no common prefix
	}
	pt.root = &PatriciaTrieNode{}
	pt.root.filterObject = NewFilter(p1, 0)
	pt.root.build(p1, bsub)

	// initialize the tdt.Core
	pt.tdtCore = tdt.NewCore()

	return pt
}
