// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"encoding/json"
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
	head.parent = nil
	head.children = []*LocalityData{}

	total := lm[""]
	for node, count := range lm {
		locality := 100 * float32(count) / float32(total)
		path := strings.Split(node, "-")
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
