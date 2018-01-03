// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"errors"
)

// SplayTree struct
type SplayTree struct {
	root *OptimalBST
}

// AnalyzeLocality increments the locality per node for the specific id
func (spt *SplayTree) AnalyzeLocality(id []byte, path string, lm *LocalityMap) {
	spt.root.AnalyzeLocality(id, path, lm)
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

	// Encode OptimalBST
	enc.Encode(spt.root)

	return buf.Bytes(), err
}

// Search returns a slice of notificationURI
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

	// Decode OptimalBST
	err = dec.Decode(&spt.root)

	return
}

// Internal helper methods -----------------------------------------------------

func (obst *OptimalBST) splaySearch(spt *SplayTree, parent *OptimalBST, id []byte) []string {
	matches := []string{}
	if obst.filterObject.Match(id) {
		matches = append(matches, obst.notificationURI)
		if obst.matchNext != nil {
			// Do not splay the subsets
			matches = append(matches, obst.matchNext.Search(id)...)
		}
		// Splay
		// 0. Check if this is the root node, do nothing if so
		if parent != nil {
			// 1. Remove this node and connect parent to the next node
			parent.mismatchNext = obst.mismatchNext
			// 2. Insert self to root
			obst.mismatchNext = spt.root
			spt.root = obst
		}
		return matches
	}
	if obst.mismatchNext != nil {
		return obst.mismatchNext.splaySearch(spt, obst, id)
	}
	return matches
}

// BuildSplayTree builds SplayTree from Subscriptions
// returns the pointer to the node node
func BuildSplayTree(sub *Subscriptions) *SplayTree {
	// make subsets to the child subscriptions of the corresponding parents
	sub.linkSubset()

	nds := NewNodes(sub)
	spt := &SplayTree{}
	spt.root = &OptimalBST{}
	spt.root = spt.root.build(sub, nds)
	return spt
}
