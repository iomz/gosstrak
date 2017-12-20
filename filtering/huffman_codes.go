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

type HuffmanCodes []*entry

func (hc HuffmanCodes) Dump() string {
	writer := &bytes.Buffer{}
	hc[0].print(writer, 0)
	hc[1].print(writer, 0)
	return writer.String()
}

type entry struct {
	filter string
	p      float64
	bit    rune
	code   string
	group  HuffmanCodes
}

func (ent *entry) equal(want *entry) (ok bool, got *entry, wanted *entry) {
	if ent.filter != want.filter ||
		ent.p != want.p ||
		ent.bit != want.bit ||
		ent.code != want.code {
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
		fmt.Fprintf(writer, "%s--%v %f %d %s\n",
			strings.Repeat(" ", indent), ent.bit, ent.p, ent.code, ent.filter)
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

func (hc HuffmanCodes) autoencode() *HuffmanCodes {
	// write bits to last and penultimate
	hc.sortByP()
	last := hc[0]
	penultimate := hc[1]
	last.bit = '0'
	penultimate.bit = '1'

	if len(hc) == 2 {
		return &hc
	}

	// Make a new composition
	ent := &entry{
		filter: "",
		p:      penultimate.p + last.p,
		bit:    0,
		code:   "",
		group:  []*entry{penultimate, last},
	}
	// Remove last and penultimate and add the composition
	nhc := append(hc[2:], ent)

	return nhc.autoencode()
}

func (hc *HuffmanCodes) gencode(code string) {
	for i := 0; i < 2; i++ {
		ent := (*hc)[i]
		if len(ent.group) == 0 {
			ent.code = code + string(ent.bit)
		} else if len(ent.group) == 2 {
			ent.group.gencode(code + string(ent.bit))
		}
	}
}

func NewHuffmanCodes(sub Subscriptions) *HuffmanCodes {
	hc := make(HuffmanCodes, len(sub))
	i := 0
	for fs, info := range sub {
		hc[i] = &entry{fs, info.EntropyValue, 0, "", []*entry{}}
		i++
	}
	hc.sortByP()
	return &hc
}
