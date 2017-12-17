package filtering

import (
	"encoding/json"
)

// LocalityData contains usage locality
// for a specific group of IDs
type LocalityData struct {
	name     string
	locality float32
	parent   *LocalityData
	children []*LocalityData
}

// LocalityDataJSON defines result JSON structure
type LocalityDataJSON struct {
	Name     string             `json:"name"`
	Parent   string             `json:"parent"`
	Value    float32            `json:"value"`
	Children []LocalityDataJSON `json:"children"`
}

// MarshalJSON overwrites marshaller for *LocalityData
func (ld *LocalityData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&[]LocalityDataJSON{ld.JSON()})
}

// JSON returns LocalityDataJSON struct for *LocalityData
func (ld *LocalityData) JSON() LocalityDataJSON {
	var parentName string
	if ld.parent == nil {
		parentName = "null"
	} else {
		parentName = ld.parent.name
	}

	if len(ld.children) == 0 {
		return LocalityDataJSON{
			Name:   ld.name,
			Parent: parentName,
			Value:  ld.locality,
		}
	}
	children := []LocalityDataJSON{}
	for _, child := range ld.children {
		children = append(children, child.JSON())
	}
	return LocalityDataJSON{
		Name:     ld.name,
		Parent:   parentName,
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
	child.parent = ld
	child.InsertLocality(path, locality)
	ld.children = append(ld.children, child)
	return
}
