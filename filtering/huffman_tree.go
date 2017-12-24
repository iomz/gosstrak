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
	//"log"
	"reflect"
	"strings"
)

// HuffmanTree struct
type HuffmanTree struct {
	notificationURI string
	filterObject    *FilterObject
	matchNext       *HuffmanTree
	mismatchNext    *HuffmanTree
}

// AnalyzeLocality increments the locality per node for the specific id
func (ht *HuffmanTree) AnalyzeLocality(id []byte, prefix string, lm *LocalityMap) {
	return
}

// Dump returs a string representation of the PatriciaTrie
func (ht *HuffmanTree) Dump() string {
	writer := &bytes.Buffer{}
	ht.print(writer, 0)
	return writer.String()
}

// MarshalBinary overwrites the marshaller in gob encoding *HuffmanTree
func (ht *HuffmanTree) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.HuffmanTree")

	// Notify
	enc.Encode(ht.notificationURI)

	// Filter
	hasFilter := ht.filterObject != nil
	enc.Encode(hasFilter)
	if hasFilter {
		err = enc.Encode(ht.filterObject)
	}

	// matchNext
	hasMatchNext := ht.matchNext != nil
	enc.Encode(hasMatchNext)
	if hasMatchNext {
		err = enc.Encode(ht.matchNext)
	}

	// mismatchNext
	hasMismatchNext := ht.mismatchNext != nil
	enc.Encode(hasMismatchNext)
	if hasMismatchNext {
		err = enc.Encode(ht.mismatchNext)
	}

	return buf.Bytes(), err
}

// Search returns a slice of notificationURI
func (ht *HuffmanTree) Search(id []byte) []string {
	matches := []string{}
	if ht.filterObject.Match(id) {
		matches = append(matches, ht.notificationURI)
		if ht.matchNext != nil {
			matches = append(matches, ht.matchNext.Search(id)...)
		}
		return matches
	}
	if ht.mismatchNext != nil {
		return ht.mismatchNext.Search(id)
	}
	return matches
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (ht *HuffmanTree) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.HuffmanTree" {
		return errors.New("Wrong Filtering Engine: " + typeOfEngine)
	}

	// Notify
	if err = dec.Decode(&ht.notificationURI); err != nil {
		return
	}

	// FilterObject
	var hasFilterObject bool
	if err = dec.Decode(&hasFilterObject); err != nil {
		return
	}
	if hasFilterObject {
		err = dec.Decode(&ht.filterObject)
	}

	// matchNext
	var hasMatchNext bool
	if err = dec.Decode(&hasMatchNext); err != nil {
		return
	}
	if hasMatchNext {
		err = dec.Decode(&ht.matchNext)
	} else {
		ht.matchNext = nil
	}

	// mismatchNext
	var hasMismatchNext bool
	if err = dec.Decode(&hasMismatchNext); err != nil {
		return
	}
	if hasMismatchNext {
		err = dec.Decode(&ht.mismatchNext)
	} else {
		ht.mismatchNext = nil
	}

	return
}

// Internal helper methods -----------------------------------------------------

func (ht *HuffmanTree) build(sub *Subscriptions, hc *HuffmanCodes) *HuffmanTree {
	node := ht
	for i, ent := range *hc {
		node.filterObject = NewFilter(ent.filter, ent.offset)
		node.notificationURI = ent.notificationURI
		if _, ok := (*sub)[ent.filter]; ok && (*sub)[ent.filter].Subset != nil {
			subset := (*sub)[ent.filter].Subset
			subhc := NewHuffmanCodes(subset)
			subht := &HuffmanTree{}
			node.matchNext = subht.build(subset, subhc)
		} else {
			node.matchNext = nil
		}
		if i+1 < len(*hc) {
			node.mismatchNext = &HuffmanTree{}
			node = ht.mismatchNext
		} else {
			node.mismatchNext = nil
		}
	}
	return ht
}

func (ht *HuffmanTree) equal(want *HuffmanTree) (ok bool, got *HuffmanTree, wanted *HuffmanTree) {
	if ht.notificationURI != want.notificationURI ||
		!reflect.DeepEqual(ht.filterObject, want.filterObject) {
		return false, ht, want
	}
	if want.matchNext != nil {
		if ht.matchNext == nil {
			return false, nil, want.matchNext
		}
		res, cgot, cwanted := ht.matchNext.equal(want.matchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	if want.mismatchNext != nil {
		if ht.mismatchNext == nil {
			return false, nil, want.mismatchNext
		}
		res, cgot, cwanted := ht.mismatchNext.equal(want.mismatchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	return true, nil, nil
}

func (ht *HuffmanTree) print(writer io.Writer, indent int) {
	var n string
	if len(ht.notificationURI) != 0 {
		n = "-> " + ht.notificationURI
	}
	fmt.Fprintf(writer, "%s--%s %s\n", strings.Repeat(" ", indent), ht.filterObject.ToString(), n)
	if ht.matchNext != nil {
		ht.matchNext.print(writer, indent+2)
	}
	if ht.mismatchNext != nil {
		ht.mismatchNext.print(writer, indent+2)
	}
}

// BuildHuffmanTree builds HuffmanTree from Subscriptions
// returns the pointer to the entry node
func BuildHuffmanTree(sub *Subscriptions) *HuffmanTree {
	// make subsets to the child subscriptions of the corresponding parents
	sub.linkSubset()

	// recalculate the cumulative evs if it has subset subscriptions
	for _, info := range *sub {
		if info.Subset != nil {
			info.EntropyValue += recalculateEntropyValue(info.Subset)
		}
	}

	hc := NewHuffmanCodes(sub)
	ht := &HuffmanTree{}
	return ht.build(sub, hc)
}
