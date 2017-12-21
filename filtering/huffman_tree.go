// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	//"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
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

// MarshalBinary overwrites the marshaller in gob encoding *HuffmanTree
func (ht *HuffmanTree) MarshalBinary() (_ []byte, err error) {
	return []byte{}, errors.New("")
	//return buf.Bytes(), err
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (ht *HuffmanTree) UnmarshalBinary(data []byte) (err error) {
	return
}

// AnalyzeLocality increments the locality per node for the specific id
func (ht *HuffmanTree) AnalyzeLocality(id []byte, prefix string, lm *LocalityMap) {
	return
}

// Search returns a slice of notificationURI
func (ht *HuffmanTree) Search(id []byte) (matches []string) {
	return
}

/*
func (ht *HuffmanTree) build(sub *Subscriptions, hc *HuffmanCodes, compLimit int) *HuffmanTree {
	if (*hc)[0].compLevel == 0 { // hc[0] comes match
	} else if (*hc)[1].compLevel == 0 { // hc[1] comes match
		if (*hc)[1].p < max((*hc)[0].group[0].p, (*hc)[0].group[1].p) {
		}
	}






	ent := (*hc)[0] // higher p
	log.Print("(*hc)[0]", ent)
	if ent.compLevel == 0 { // this entry is a filter
		ht.notificationURI = (*sub)[ent.filter].NotificationURI
		ht.filterObject = NewFilter(ent.filter, 0)
		ht.matchNext = nil
		if (*sub)[ent.filter].Subset != nil { // if has any subset
			subset := &HuffmanTree{}
			ht.matchNext = subset.buildSubset(*(*sub)[ent.filter].Subset)
		}
	} else if ent.compLevel > compLimit { // unpack the group
		ht = ht.build(sub, &ent.group, compLimit)
	} else { // ok to use composition
		f0 := NewFilter(ent.group[0].filter, 0)
		f1 := NewFilter(ent.group[1].filter, 0)
		composition := NewComposition([]*FilterObject{f0, f1})
		ht.filterObject = composition.filterObject
		next := &HuffmanTree{}
		newGroup := ent.group.reconstruct(composition.children)
		ht.matchNext = next.build(sub, newGroup, compLimit)
	}

	ent = (*hc)[1] // lower p
	log.Print("(*hc)[1]", ent)
	if ent.compLevel == 0 { // this entry is a filter
		next := &HuffmanTree{}
		next.notificationURI = (*sub)[ent.filter].NotificationURI
		next.filterObject = NewFilter(ent.filter, 0)
		ht.mismatchNext = next
		if (*sub)[ent.filter].Subset != nil { // if has any subset
			subset := &HuffmanTree{}
			ht.mismatchNext.matchNext = subset.buildSubset(*(*sub)[ent.filter].Subset)
		}
	} else if ent.compLevel > compLimit { // this entry becomes composite
		next := &HuffmanTree{}
		next = next.build(sub, &ent.group, compLimit)
		ht.mismatchNext = next
	} else {

	}

	return ht
}

func (ht *HuffmanTree) buildSubset(subset Subscriptions) *HuffmanTree {
	return ht
}

*/

func (ht *HuffmanTree) equal(want *HuffmanTree) (ok bool, got *HuffmanTree, wanted *HuffmanTree) {
	if ht.notificationURI != want.notificationURI ||
		!reflect.DeepEqual(ht.filterObject, want.filterObject) ||
		!reflect.DeepEqual(ht.matchNext, want.matchNext) ||
		!reflect.DeepEqual(ht.mismatchNext, want.mismatchNext) {
		return false, ht, want
	}
	if want.matchNext != nil {
		res, cgot, cwanted := ht.matchNext.equal(want.matchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	if want.mismatchNext != nil {
		res, cgot, cwanted := ht.mismatchNext.equal(want.mismatchNext)
		if !res {
			return res, cgot, cwanted
		}
	}
	return true, nil, nil
}

// Dump returs a string representation of the PatriciaTrie
func (ht *HuffmanTree) Dump() string {
	writer := &bytes.Buffer{}
	ht.print(writer, 0)
	return writer.String()
}

func makeBranch(sub *Subscriptions, isComposition bool, first *entry, second *entry) *HuffmanTree {
	ht := &HuffmanTree{}
	if isComposition { // make a composition from first and second
		if first.compLevel == 0 && second.compLevel == 0 {
			f0 := NewFilter(first.filter, first.offset)
			f1 := NewFilter(second.filter, second.offset)
			composition := NewComposition([]*FilterObject{f0, f1})
			nf := NewFilter(composition.filter, composition.offset)

			if nf.IsTransparent() { // if the composition is meaningless
				ht.filterObject = f0
				ht.notificationURI = first.notificationURI
				ht.mismatchNext = &HuffmanTree{}
				ht.mismatchNext.filterObject = f1
				ht.mismatchNext.notificationURI = second.notificationURI
			} else {
				// composition
				ht.filterObject = nf
				// handle first
				ht.matchNext = &HuffmanTree{}
				ht.matchNext.filterObject = composition.children[first.filter]
				ht.matchNext.notificationURI = first.notificationURI
				// handle second
				ht.matchNext.mismatchNext = &HuffmanTree{}
				ht.matchNext.mismatchNext.filterObject = composition.children[second.filter]
				ht.matchNext.mismatchNext.notificationURI = second.notificationURI
			}
		} else if first.compLevel != 0 && second.compLevel == 0 {
			f0 := NewFilter(first.branch.filterObject.String, first.branch.filterObject.Offset)
			f1 := NewFilter(second.filter, second.offset)
			composition := NewComposition([]*FilterObject{f0, f1})
			nf := NewFilter(composition.filter, composition.offset)

			log.Print(nf.ToString())
			if nf.IsTransparent() { // if the composition is meaningless
				// handle first
				ht = first.branch
				// handle second
				ht.mismatchNext = &HuffmanTree{}
				ht.mismatchNext.filterObject = NewFilter(second.filter, second.offset)
				ht.mismatchNext.notificationURI = second.notificationURI
			} else {
				// composition
				ht.filterObject = nf
				// handle first
				ht.matchNext = first.branch
				ht.matchNext.filterObject = composition.children[first.branch.filterObject.String]
				// handle second
				ht.matchNext.mismatchNext = &HuffmanTree{}
				ht.matchNext.mismatchNext.filterObject = composition.children[second.filter]
				ht.matchNext.mismatchNext.notificationURI = second.notificationURI
			}
		} else if first.compLevel == 0 && second.compLevel != 0 {
			f0 := NewFilter(first.filter, first.offset)
			f1 := NewFilter(second.branch.filterObject.String, second.branch.filterObject.Offset)
			composition := NewComposition([]*FilterObject{f0, f1})
			nf := NewFilter(composition.filter, composition.offset)

			if nf.IsTransparent() { // if the composition is meaningless
				// handle first
				ht.filterObject = NewFilter(first.filter, first.offset)
				ht.notificationURI = first.notificationURI
				// handle second
				ht.mismatchNext = second.branch
			} else {
				// handle first
				ht.filterObject = nf
				ht.matchNext = &HuffmanTree{}
				ht.matchNext.filterObject = composition.children[first.filter]
				ht.matchNext.notificationURI = first.notificationURI
				// handle second
				ht.mismatchNext = second.branch
				ht.mismatchNext.filterObject = composition.children[second.branch.filterObject.String]
			}
		} else if first.compLevel != 0 && second.compLevel != 0 {
			f0 := NewFilter(first.branch.filterObject.String, first.branch.filterObject.Offset)
			f1 := NewFilter(second.branch.filterObject.String, second.branch.filterObject.Offset)
			composition := NewComposition([]*FilterObject{f0, f1})
			nf := NewFilter(composition.filter, composition.offset)

			// handle first
			ht.matchNext = first.branch
			if !nf.IsTransparent() {
				ht.matchNext.filterObject = composition.children[first.branch.filterObject.String]
			}

			// handle second
			ht.mismatchNext = second.branch
			if !nf.IsTransparent() {
				ht.mismatchNext.filterObject = composition.children[second.branch.filterObject.String]
			}
		}
	} else { // not making a composition, just a branch of first, then second
		if first.compLevel == 0 {
			ht.filterObject = NewFilter(first.filter, first.offset)
			ht.notificationURI = first.notificationURI
		} else {
			ht = first.branch
		}
		if second.compLevel == 0 {
			ht.mismatchNext = &HuffmanTree{}
			ht.mismatchNext.filterObject = NewFilter(second.filter, second.offset)
			ht.mismatchNext.notificationURI = second.notificationURI
			/*
				if (*sub)[second.filter].Subset != nil {
					ht.mismatchNext.matchNext = &HuffmanTree{}
					ht.mismatchNext.matchNext = ht.matchNext.makeSubsetBranch((*sub)[second.filter].Subset)
				} else {
					ht.mismatchNext.matchNext = nil
				}
			*/
		} else {
			ht.mismatchNext = second.branch
		}
	}
	return ht
}

func (ht *HuffmanTree) makeSubsetBranch(sub *Subscriptions) *HuffmanTree {
	return ht
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
func BuildHuffmanTree(sub Subscriptions, compLimit int) *HuffmanTree {
	hc := NewHuffmanCodes(sub)
	sub.linkSubset()
	hc = hc.autoencode(compLimit)
	ht := &HuffmanTree{}
	return ht
}
