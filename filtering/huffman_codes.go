// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"sort"
)

// HuffmanCodes store a slice of *entry to
// sort by p to build HuffmanTree
type HuffmanCodes []*entry

type entry struct {
	filter          string
	offset          int
	notificationURI string
	p               float64
}

func (ent *entry) equal(want *entry) (ok bool, got *entry, wanted *entry) {
	if ent.filter != want.filter ||
		ent.offset != want.offset ||
		ent.p != want.p {
		return false, ent, want
	}
	return true, nil, nil
}

func (hc *HuffmanCodes) sortByP() {
	sort.Slice(*hc, func(i, j int) bool {
		return (*hc)[i].p > (*hc)[j].p
	})
	return
}

// NewHuffmanCodes returns a pointer to HuffmanCodes
// generated by the given Subscriptions
func NewHuffmanCodes(sub *Subscriptions) *HuffmanCodes {
	hc := make(HuffmanCodes, len(*sub))
	i := 0
	for fs, info := range *sub {
		hc[i] = &entry{fs, info.Offset, info.NotificationURI, info.EntropyValue}
		i++
	}
	hc.sortByP()
	return &hc
}