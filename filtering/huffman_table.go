// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"sort"
)

type HuffmanCode struct {
	FilterString []string
	EntropyValue float64
}

type HuffmanTable []HuffmanCode

func (ht HuffmanTable) sort() {
	sort.Slice(ht, func(i, j int) bool {
		return ht[i].EntropyValue < ht[j].EntropyValue
	})
	return
}
