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
	code   int
	group  []*entry
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

func (hc *HuffmanCodes) sortByCode() {
	sort.Slice(*hc, func(i, j int) bool {
		return (*hc)[i].code < (*hc)[j].code
	})
	return
}

func (hc HuffmanCodes) autoencode() *HuffmanCodes {
	if len(hc) == 2 {
		return &hc
	}
	last := hc[0]
	penultimate := hc[1]
	ent := &entry{
		filter: "",
		p:      penultimate.p + last.p,
		bit:    0,
		code:   -1,
		group:  []*entry{penultimate, last},
	}
	nhc := append(hc[2:], ent)
	nhc.sortByP()
	return nhc.autoencode()
}

func NewHuffmanCodes(sub Subscriptions) *HuffmanCodes {
	hc := make(HuffmanCodes, len(sub))
	i := 0
	for fs, info := range sub {
		hc[i] = &entry{fs, info.EntropyValue, 0, -1, []*entry{}}
		i++
	}
	hc.sortByP()
	return &hc
}
