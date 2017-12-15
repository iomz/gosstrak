package filter

import (
	"encoding/json"
)

// PatriciaTrieLocalityData contains usage locality
// for a specific group of IDs
type PatriciaTrieLocalityData struct {
	name     string
	locality float32
	parent   *PatriciaTrieLocalityData
	children []*PatriciaTrieLocalityData
}

// PTLDJSON defines result JSON structure
type PTLDJSON struct {
	Name     string     `json:"name"`
	Parent   string     `json:"parent"`
	Value    float32    `json:"value"`
	Children []PTLDJSON `json:"children"`
}

// MarshalJSON overwrites marshaller for *PatriciaTrieLocalityData
func (ptld *PatriciaTrieLocalityData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&[]PTLDJSON{ptld.JSON()})
}

// JSON returns PTLDJSON struct for *PatriciaTrieLocalityData
func (ptld *PatriciaTrieLocalityData) JSON() PTLDJSON {
	var parentName string
	if ptld.parent == nil {
		parentName = "null"
	} else {
		parentName = ptld.parent.name
	}

	if len(ptld.children) == 0 {
		return PTLDJSON{
			Name:   ptld.name,
			Parent: parentName,
			Value:  ptld.locality,
		}
	}
	children := []PTLDJSON{}
	for i := 0; i < len(ptld.children); i++ {
		children = append(children, ptld.children[i].JSON())
	}
	return PTLDJSON{
		Name:     ptld.name,
		Parent:   parentName,
		Value:    ptld.locality,
		Children: children,
	}
}

// InsertLocality recursively generate PatriciaTrieLocalityData
// to the node with locality
func (ptld *PatriciaTrieLocalityData) InsertLocality(path []string, locality float32) {
	if len(ptld.name) == 0 {
		ptld.name = path[0]
	}
	path = path[1:]

	// This node is the leaf
	if len(path) == 0 {
		ptld.locality = locality
		return
	}

	// If this node has any child
	if len(ptld.children) != 0 {
		for i := 0; i < len(ptld.children); i++ {
			// If found
			if ptld.children[i].name == path[0] {
				ptld.children[i].InsertLocality(path, locality)
				return
			}
		}
	} else {
		ptld.children = []*PatriciaTrieLocalityData{}
	}

	// Append a new child
	child := &PatriciaTrieLocalityData{}
	child.parent = ptld
	child.InsertLocality(path, locality)
	ptld.children = append(ptld.children, child)
	return
}
