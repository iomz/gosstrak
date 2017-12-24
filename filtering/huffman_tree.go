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
	return []byte{}, errors.New("")
	//return buf.Bytes(), err
}

// Search returns a slice of notificationURI
func (ht *HuffmanTree) Search(id []byte) (matches []string) {
	return
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (ht *HuffmanTree) UnmarshalBinary(data []byte) (err error) {
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

/*
// makeBranch makes a branch from given 2 entries
// this function is called from HuffmanTree.build() to
// autoencode the HuffmanCodes by p values
// sub: pointer of the subscriptions
// isComposition: if this branch is a composition
// first: the penultimate entry in HuffmanCodes after sorted by p
// second: the last entry in HuffmanCodes after sorted by p
func makeBranch(sub *Subscriptions, isComposition bool, first *entry, second *entry) *HuffmanTree {
	ht := &HuffmanTree{}
	// makeBranch() is making a branch with a composite filter
	// consisting of first and second
	if isComposition {
		// the process of constructing a composite branch has
		// 4*2*4 different cases
		// 1. both first and second are not composite branch
		// 1-1. if the resulting composite filter is meaningless (transparent)
		// 1-1-1. both first and second are not branch
		// 1-1-2. first is a branch but not second
		// 1-1-3. first is not a branch but second is
		// 1-1-4. both first and second are branches
		// 1-2. if the resulting composite filter is NOT meaningless
		// <abbr.>
		// 2. first is a composite branch but second is not
		// <abbr.>
		// 3. first is not a composite branch but second is
		// <abbr.>
		// 4. both first and second are composite branches
		// <abbr.>
		if first.compLevel == 0 && second.compLevel == 0 { // 1. both not composite branch
			f0 := NewFilter(first.filter, first.offset)
			f1 := NewFilter(second.filter, second.offset)
			composition := NewComposition([]*FilterObject{f0, f1})
			nf := NewFilter(composition.filter, composition.offset)

			if nf.IsTransparent() { // 1-1.
				// first
				if first.branch == nil { // 1-1-[1,3].
					ht.filterObject = f0
					ht.notificationURI = first.notificationURI
				} else { // 1-1-[2,4].
					ht = first.branch
				}
				// second
				if ht.mismatchNext == nil { // 1-1-[1,2,3,4].
					if second.branch == nil { // 1-1-[1,2].
						ht.mismatchNext = &HuffmanTree{}
						ht.mismatchNext.filterObject = f1
						ht.mismatchNext.notificationURI = second.notificationURI
					} else { // 1-1-[3,4].
						ht.mismatchNext = second.branch
					}
				} else { // 1-1-[1,2,3,4].
					if second.branch == nil { // 1-1-[1,2].
						swapTarget := ht.mismatchNext
						newMismatchNext := &HuffmanTree{}
						newMismatchNext.mismatchNext = swapTarget
						newMismatchNext.filterObject = f1
						newMismatchNext.notificationURI = second.notificationURI
						ht.mismatchNext = newMismatchNext
					} else { // 1-1-[3,4]
					}
				}
			} else { // 1-2.
				// composition
				ht.filterObject = nf
				// first
				ht.matchNext = &HuffmanTree{}
				ht.matchNext.filterObject = composition.children[first.filter]
				ht.matchNext.notificationURI = first.notificationURI
				// second
				ht.matchNext.mismatchNext = &HuffmanTree{}
				ht.matchNext.mismatchNext.filterObject = composition.children[second.filter]
				ht.matchNext.mismatchNext.notificationURI = second.notificationURI
			}
		} else if first.compLevel != 0 && second.compLevel == 0 { // 2. first is composite but second is not
			f0 := NewFilter(first.branch.filterObject.String, first.branch.filterObject.Offset)
			f1 := NewFilter(second.filter, second.offset)
			composition := NewComposition([]*FilterObject{f0, f1})
			nf := NewFilter(composition.filter, composition.offset)

			if nf.IsTransparent() { // 2-1.
				// first
				ht = first.branch
				// second
				swapTarget := ht.mismatchNext
				if swapTarget == nil {
					ht.mismatchNext = &HuffmanTree{}
				} else {
					newMismatchNext := &HuffmanTree{}
					newMismatchNext.mismatchNext = swapTarget
					ht.mismatchNext = newMismatchNext
				}
				ht.mismatchNext.filterObject = NewFilter(second.filter, second.offset)
				ht.mismatchNext.notificationURI = second.notificationURI
			} else { // 2-2.
				// composition
				ht.filterObject = nf
				// first
				ht.matchNext = first.branch
				ht.matchNext.filterObject = composition.children[first.branch.filterObject.String]
				// second
				swapTarget := ht.matchNext.mismatchNext
				if swapTarget == nil {
					ht.matchNext.mismatchNext = &HuffmanTree{}
				} else {
					newMismatchNext := &HuffmanTree{}
					newMismatchNext.mismatchNext = swapTarget
					ht.matchNext.mismatchNext = newMismatchNext
				}
				ht.matchNext.mismatchNext = &HuffmanTree{}
				ht.matchNext.mismatchNext.filterObject = composition.children[second.filter]
				ht.matchNext.mismatchNext.notificationURI = second.notificationURI
			}
		} else if first.compLevel == 0 && second.compLevel != 0 { // 3. first is not composite but second is
			f0 := NewFilter(first.filter, first.offset)
			f1 := NewFilter(second.branch.filterObject.String, second.branch.filterObject.Offset)
			composition := NewComposition([]*FilterObject{f0, f1})
			nf := NewFilter(composition.filter, composition.offset)

			if nf.IsTransparent() { // 3-1.
				// first
				ht.filterObject = NewFilter(first.filter, first.offset)
				ht.notificationURI = first.notificationURI
				// second
				ht.mismatchNext = second.branch
			} else { // 3-2.
				// composition
				ht.filterObject = nf
				// first
				swapTarget := ht.matchNext
				if swapTarget == nil {
					ht.matchNext = &HuffmanTree{}
				} else {
					newMatchNext := &HuffmanTree{}
					newMatchNext.matchNext = swapTarget
					ht.matchNext = newMatchNext
				}
				ht.matchNext.filterObject = composition.children[first.filter]
				ht.matchNext.notificationURI = first.notificationURI
				// second
				ht.mismatchNext = second.branch
				ht.mismatchNext.filterObject = composition.children[second.branch.filterObject.String]
			}
		} else if first.compLevel != 0 && second.compLevel != 0 { // 4. both are composite branches
			f0 := NewFilter(first.branch.filterObject.String, first.branch.filterObject.Offset)
			f1 := NewFilter(second.branch.filterObject.String, second.branch.filterObject.Offset)
			composition := NewComposition([]*FilterObject{f0, f1})
			nf := NewFilter(composition.filter, composition.offset)

			if nf.IsTransparent() { // 4-1.
				// first
				ht = first.branch
				// second
				swapTarget := ht.mismatchNext
				if swapTarget == nil {
					ht.mismatchNext = second.branch
				} else {
					newMismatchNext := &HuffmanTree{}
					newMismatchNext.mismatchNext = swapTarget
					ht.mismatchNext = newMismatchNext
				}
			} else { // 4-2.
				// composition
				ht.filterObject = nf
				// first
				ht.matchNext = first.branch
				ht.matchNext.filterObject = composition.children[first.branch.filterObject.String]
				// second
				ht.mismatchNext = second.branch
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
		} else {
			ht.mismatchNext = second.branch
		}
	}
	return ht
}
*/
