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

// Subscriptions contains filter string as key and Info as value
type Subscriptions map[string]*Info

// Info contains notificationURI and pValue for a filter
type Info struct {
	NotificationURI string
	EntropyValue    float64
	Subset          *Subscriptions
}

// Dump retuns a string representation of the Subscriptions
func (sub Subscriptions) Dump() string {
	writer := &bytes.Buffer{}
	sub.print(writer, 0)
	return writer.String()
}

func (sub Subscriptions) keys() []string {
	ks := make([]string, len(sub))
	i := 0
	for k := range sub {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	return ks
}

func (sub Subscriptions) linkSubset() {
	hc := *NewHuffmanCodes(sub)
	for _, ent := range hc {
		// an entry must have two members in the group
		if len(ent.group) == 0 {
			for fs, info := range sub {
				linkCandidate := ent.filter
				// check if fs is a subset of the linkCandidate
				if strings.HasPrefix(fs, linkCandidate) &&
					fs != linkCandidate { // they shouldn't be the same
					if sub[linkCandidate].Subset == nil {
						sub[linkCandidate].Subset = &Subscriptions{fs: info}
					} else {
						(*sub[linkCandidate].Subset)[fs] = info
					}
					// recursively link the subset
					sub[linkCandidate].Subset.linkSubset()
					// also, add the subset's EntropyValue to the link's
					sub[linkCandidate].EntropyValue += info.EntropyValue
					// finaly delete the filter from the upper Subscriptions
					delete(sub, fs)
				}
			}
		}
	}
	return
}

func (sub Subscriptions) print(writer io.Writer, indent int) {
	for fs, info := range sub {
		fmt.Fprintf(writer, "%s--%s %f\n", strings.Repeat(" ", indent), fs, info.EntropyValue)
		if info.Subset != nil {
			info.Subset.print(writer, indent+2)
		}
	}
}
