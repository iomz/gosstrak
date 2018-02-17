// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"encoding/json"
	"sort"
	"strings"
)

// LocalityMap contains PatriciaTrie node by prefix
// and its reference counts
type LocalityMap map[string]int

// ToJSON returns LocalityData in JSON format
func (lm LocalityMap) ToJSON() []byte {
	head := new(LocalityData)
	head.name = "Entry Node"
	head.locality = 100
	head.children = []*LocalityData{}

	total := 0
	if _, ok := lm[""]; ok {
		total = lm[""]
	} else {
		head.name = lm.first().NotificationURI
		total = lm.first().Locality
	}

	for node, count := range lm {
		locality := 100 * float32(count) / float32(total)
		path := strings.Split(node, ",")
		// Root node
		if len(path) == 1 {
			continue
		}
		head.InsertLocality(path, locality)
	}

	// Construct ld
	res, _ := json.Marshal(head)
	return res
}

// LocalityPair contains the NotificationURI and the Locality for the node
type LocalityPair struct {
	NotificationURI string
	Locality        int
}

func (lm LocalityMap) first() LocalityPair {
	nl := make([]LocalityPair, len(lm))
	i := 0
	for n, l := range lm {
		nl[i] = LocalityPair{n, l}
		i++
	}
	sort.Slice(nl, func(i, j int) bool {
		return nl[i].Locality > nl[j].Locality
	})
	return nl[0]
}

// LocalityData contains usage locality
// for a specific group of IDs
type LocalityData struct {
	name     string
	locality float32
	children []*LocalityData
}

// LocalityDataJSON defines result JSON structure
type LocalityDataJSON struct {
	Name     string             `json:"name"`
	Value    float32            `json:"value"`
	Children []LocalityDataJSON `json:"children"`
}

// MarshalJSON overwrites marshaller for *LocalityData
func (ld *LocalityData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&[]LocalityDataJSON{ld.JSON()})
}

// JSON returns LocalityDataJSON struct for *LocalityData
func (ld *LocalityData) JSON() LocalityDataJSON {
	if len(ld.children) == 0 {
		return LocalityDataJSON{
			Name:  ld.name,
			Value: ld.locality,
		}
	}
	children := []LocalityDataJSON{}
	for _, child := range ld.children {
		children = append(children, child.JSON())
	}
	return LocalityDataJSON{
		Name:     ld.name,
		Value:    ld.locality,
		Children: children,
	}
}

// InsertLocality recursively generate LocalityData
// to the node with locality
func (ld *LocalityData) InsertLocality(path []string, locality float32) {
	if len(ld.name) == 0 {
		ld.name = path[0]
	}
	path = path[1:]

	// This node is the leaf
	if len(path) == 0 {
		ld.locality = locality
		return
	}

	// If this node has any child
	if len(ld.children) != 0 {
		for _, child := range ld.children {
			// If found
			if child.name == path[0] {
				child.InsertLocality(path, locality)
				return
			}
		}
	} else {
		ld.children = []*LocalityData{}
	}

	// Append a new child
	child := &LocalityData{}
	child.InsertLocality(path, locality)
	ld.children = append(ld.children, child)
	if len(ld.children) == 2 {
		sort.Slice(ld.children, func(i, j int) bool {
			return ld.children[j].name == "Match" || ld.children[j].name == "Mismatch"
		})
	}
	return
}
