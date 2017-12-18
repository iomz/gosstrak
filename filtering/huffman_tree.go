// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	//"encoding/gob"
	"errors"
	//"fmt"
	"io"
)

// HuffmanTree struct
type HuffmanTree struct {
	notificationURI string
	filter          *FilterObject
	matchNext       *HuffmanTree
	mismatchNext    *HuffmanTree
	subset          []*HuffmanTree
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

func (ht *HuffmanTree) build() {
}

// Dump returs a string representation of the PatriciaTrie
func (ht *HuffmanTree) Dump() string {
	writer := &bytes.Buffer{}
	ht.print(writer, 0)
	return writer.String()
}

func (ht *HuffmanTree) print(writer io.Writer, indent int) {
}

// BuildHuffmanTree builds HuffmanTree from Subscriptions
// returns the pointer to the entry node
func BuildHuffmanTree(sub Subscriptions) *HuffmanTree {
	hc := NewHuffmanCodes(sub)
	sub.linkSubset()
	hc = hc.autoencode()
	/*
		for _, hc := range *huffmanTable {
			fmt.Print(hc.Dump())
		}
	*/
	return build(sub, hc)
}

func build(sub Subscriptions, hc *HuffmanCodes) *HuffmanTree {
	ht := &HuffmanTree{}
	/*
		for i, hc := range *huffmanTable {
			if i == 0 { // MatchPath
				group := hc.FilterGroup.members
				if len(group) == 2 { // if this has members
					if len(group[0].members) == 0 && len(group[1].members) == 0 {
					f0 := NewFilter(group[0].filter, 0)
					f1 := NewFilter(group[1].filter, 0)
					NewComposition()
				} else {
					// append subset to match
					sub[hc.FilterGroup.filter]
				}
			} else if i == 1 { // MismatchPath
			}
		}
	*/
	return ht
}
