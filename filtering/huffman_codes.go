// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// HuffmanCodes store a slice of *entry to
// sort by p to build HuffmanTree
type HuffmanCodes []*entry

// Dump returns a string representation of HuffmanCodes
func (hc HuffmanCodes) Dump() string {
	writer := &bytes.Buffer{}
	hc[0].print(writer, 0)
	hc[1].print(writer, 0)
	return writer.String()
}

type entry struct {
	filter          string
	offset          int
	notificationURI string
	p               float64
	compLevel       int
	group           HuffmanCodes
	branch          *HuffmanTree
}

func (ent *entry) equal(want *entry) (ok bool, got *entry, wanted *entry) {
	if ent.filter != want.filter ||
		ent.offset != want.offset ||
		ent.p != want.p ||
		ent.compLevel != want.compLevel {
		return false, ent, want
	}
	for i, child := range want.group {
		if len(ent.group) != 2 {
			return false, nil, child
		}
		res, cgot, cwanted := ent.group[i].equal(child)
		if !res {
			return res, cgot, cwanted
		}
	}
	return true, nil, nil
}

func (ent *entry) print(writer io.Writer, indent int) {
	if len(ent.group) == 0 && len(ent.filter) != 0 {
		fmt.Fprintf(writer, "%s--%f %d %s %d\n",
			strings.Repeat(" ", indent), ent.p, ent.compLevel, ent.filter, ent.offset)
	} else if len(ent.group) == 2 { // has group
		ent.group[0].print(writer, indent+2)
		ent.group[1].print(writer, indent+2)
	}
}

func (hc *HuffmanCodes) sortByP() {
	sort.Slice(*hc, func(i, j int) bool {
		return (*hc)[i].p < (*hc)[j].p
	})
	return
}

func (hc HuffmanCodes) autoencode(compLimit int) *HuffmanCodes {
	// write bits to last and penultimate
	hc.sortByP()
	last := hc[0]
	penultimate := hc[1]

	if len(hc) == 2 {
		return &hc
	}

	// If the pair exceeds the compLimit
	var ent *entry
	if max(last.compLevel+1, penultimate.compLevel+1) > compLimit || compLimit == -1 {
	} else { // if not, make a new composition
		ent = &entry{
			filter:    "",
			offset:    0,
			p:         penultimate.p + last.p,
			compLevel: max(last.compLevel+1, penultimate.compLevel+1),
			group:     []*entry{penultimate, last},
			branch:    &HuffmanTree{},
		}
	}

	// Remove last and penultimate and add the composition
	nhc := append(hc[2:], ent)

	return nhc.autoencode(compLimit)
}

// NewHuffmanCodes returns a pointer to HuffmanCodes
// generated by the given Subscriptions
func NewHuffmanCodes(sub Subscriptions) *HuffmanCodes {
	hc := make(HuffmanCodes, len(sub))
	i := 0
	for fs, info := range sub {
		hc[i] = &entry{fs, 0, info.NotificationURI, info.EntropyValue, 0, []*entry{}, nil}
		i++
	}
	hc.sortByP()
	return &hc
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
