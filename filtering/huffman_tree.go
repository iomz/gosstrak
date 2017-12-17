// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	//"encoding/gob"
	"errors"
	"io"
)

// HuffmanTree struct
type HuffmanTree struct {
	notificationURI string
	filter          *FilterObject
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

// BuildHuffmanTree builds HuffmanTree from filter.Map
// returns the pointer to the entry node
func BuildHuffmanTree(sub Subscriptions) *HuffmanTree {
	return &HuffmanTree{}
}
