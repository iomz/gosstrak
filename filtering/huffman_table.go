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

type group struct {
	members []*group
	filter  string
}

func (g *group) print(writer io.Writer, indent int) {
	if len(g.members) == 0 && len(g.filter) != 0 {
		fmt.Fprintf(writer, "%s--%s\n", strings.Repeat(" ", indent), g.filter)
	} else {
		for _, gg := range g.members {
			gg.print(writer, indent+2)
		}
	}
}

func (g *group) unite(m *group) *group {
	return &group{append([]*group{g}, m), ""}
}

type HuffmanCode struct {
	FilterGroup  *group
	EntropyValue float64
}

func (hc HuffmanCode) Dump() string {
	writer := &bytes.Buffer{}
	hc.FilterGroup.print(writer, 0)
	return writer.String()
}

type HuffmanTable []HuffmanCode

func (ht HuffmanTable) sort() {
	sort.Slice(ht, func(i, j int) bool {
		return ht[i].EntropyValue < ht[j].EntropyValue
	})
	return
}

func (ht HuffmanTable) autoencode() *HuffmanTable {
	if len(ht) == 2 {
		return &ht
	}
	last := ht[0]
	penultimate := ht[1]
	hc := HuffmanCode{
		penultimate.FilterGroup.unite(last.FilterGroup),
		penultimate.EntropyValue + last.EntropyValue,
	}
	nht := append(ht[2:], hc)
	nht.sort()
	return nht.autoencode()
}
