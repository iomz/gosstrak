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
	Offset          int
	NotificationURI string
	EntropyValue    float64
	Subset          Subscriptions
}

// Clone retuns a new copy of subscriptions
func (sub Subscriptions) Clone() Subscriptions {
	clone := Subscriptions{}

	for _, fs := range sub.Keys() {
		info := Info{}
		info.Offset = sub[fs].Offset
		info.NotificationURI = sub[fs].NotificationURI
		info.EntropyValue = sub[fs].EntropyValue
		if len(sub[fs].Subset) == 0 {
			info.Subset = Subscriptions{}
		} else {
			info.Subset = sub[fs].Subset.Clone()
		}
		clone[fs] = &info
	}
	return clone
}

// Dump retuns a string representation of the Subscriptions
func (sub Subscriptions) Dump() string {
	writer := &bytes.Buffer{}
	sub.print(writer, 0)
	return writer.String()
}

// Keys return a slice of keys in m
func (sub Subscriptions) Keys() []string {
	ks := make([]string, len(sub))
	i := 0
	for k := range sub {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	return ks
}

// linkSubset finds subsets and nest them under the parents
func (sub Subscriptions) linkSubset() {
	nds := NewNodes(sub)
	for _, nd := range nds {
		for _, fs := range sub.Keys() {
			info := sub[fs]
			linkCandidate := nd.filter
			// check if fs is a subset of the linkCandidate
			if strings.HasPrefix(fs, linkCandidate) &&
				fs != linkCandidate { // they shouldn't be the same
				// if there is no subset already
				if len(sub[linkCandidate].Subset) == 0 {
					sub[linkCandidate].Subset = Subscriptions{}
					sub[linkCandidate].Subset[fs[len(linkCandidate):]] = &Info{
						Offset:          info.Offset + len(linkCandidate),
						NotificationURI: info.NotificationURI,
						EntropyValue:    info.EntropyValue,
					}
				} else {
					sub[linkCandidate].Subset[fs[len(linkCandidate):]] = &Info{
						Offset:          info.Offset + len(linkCandidate),
						NotificationURI: info.NotificationURI,
						EntropyValue:    info.EntropyValue,
					}
				}
				// recursively link the subset
				sub[linkCandidate].Subset.linkSubset()
				// finaly delete the filter from the upper Subscriptions
				delete(sub, fs)
			}
		}
	}
	return
}

func recalculateEntropyValue(sub Subscriptions) float64 {
	ev := float64(0)
	for _, fs := range sub.Keys() {
		ss := sub[fs].Subset
		if len(ss) != 0 {
			sub[fs].EntropyValue += recalculateEntropyValue(ss)
		}
		ev += sub[fs].EntropyValue
	}
	return ev
}

func (sub Subscriptions) print(writer io.Writer, indent int) {
	for _, fs := range sub.Keys() {
		fmt.Fprintf(writer, "%s--%s %f\n", strings.Repeat(" ", indent), fs, sub[fs].EntropyValue)
		ss := sub[fs].Subset
		if len(ss) != 0 {
			ss.print(writer, indent+2)
		}
	}
}
