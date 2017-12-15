package filter

import (
	"encoding/json"
	"strings"
)

// PatriciaTrieLocalityMap contains PatriciaTrie node by prefix
// and its reference counts
type PatriciaTrieLocalityMap map[string]int

// ToJSON returns PatriciaTrieLocalityData in JSON format
func (ptlm PatriciaTrieLocalityMap) ToJSON() []byte {
	entry := new(PatriciaTrieLocalityData)

	total := ptlm[""]
	entry.name = "Entry"
	entry.locality = 100
	entry.parent = nil
	entry.children = []*PatriciaTrieLocalityData{}

	for node, count := range ptlm {
		locality := 100 * float32(count) / float32(total)
		path := strings.Split(node, "-")
		// Root node
		if len(path) == 1 {
			continue
		}
		entry.InsertLocality(path, locality)
	}

	// Construct ptld
	res, _ := json.Marshal(entry)
	return res
}
