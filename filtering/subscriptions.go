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
	"sync"
)

// SubMap contains filter string as keys and &Info as values
type SubMap map[string]*Info

// Subscriptions contains filter string as key and Info as value
type Subscriptions struct {
	sync.RWMutex
	m SubMap
}

// Info contains notificationURI and pValue for a filter
type Info struct {
	Offset          int
	NotificationURI string
	EntropyValue    float64
	Subset          Subscriptions
}

// Clone retuns a new copy of subscriptions
func (sub Subscriptions) Clone() Subscriptions {
	return Subscriptions{}
}

// Delete takes a key and delete it
func (sub Subscriptions) Delete(s string) {
	sub.Lock()
	delete(sub.m, s)
	sub.Unlock()
}

// Dump retuns a string representation of the Subscriptions
func (sub Subscriptions) Dump() string {
	writer := &bytes.Buffer{}
	sub.print(writer, 0)
	return writer.String()
}

// Get takes a key and returns value
func (sub Subscriptions) Get(s string) *Info {
	sub.RLock()
	info := sub.m[s]
	sub.RUnlock()
	return info
}

// Keys return a slice of keys in m
func (sub Subscriptions) Keys() []string {
	ks := make([]string, len(sub.m))
	i := 0
	sub.RLock()
	for k := range sub.m {
		ks[i] = k
		i++
	}
	sub.RUnlock()
	sort.Strings(ks)
	return ks
}

// Length returns len(sub.m)
func (sub Subscriptions) Length() int {
	sub.RLock()
	len := len(sub.m)
	sub.RUnlock()
	return len
}

// Has takes a key and returns true if it's in the map
func (sub Subscriptions) Has(s string) bool {
	sub.RLock()
	_, ok := sub.m[s]
	sub.RUnlock()
	return ok
}

// Set takes a key and a value
func (sub Subscriptions) Set(s string, info *Info) {
	sub.Lock()
	if sub.m == nil {
		sub.m = make(SubMap)
	}
	sub.m[s] = info
	sub.Unlock()
}

// linkSubset finds subsets and nest them under the parents
func (sub Subscriptions) linkSubset() {
	nds := NewNodes(sub)
	for _, nd := range nds {
		for _, fs := range sub.Keys() {
			info := sub.Get(fs)
			linkCandidate := nd.filter
			// check if fs is a subset of the linkCandidate
			if strings.HasPrefix(fs, linkCandidate) &&
				fs != linkCandidate { // they shouldn't be the same
				// if there is no subset already
				if sub.Get(linkCandidate).Subset.Length() == 0 {
					sub.Get(linkCandidate).Subset.Set(fs[len(linkCandidate):], &Info{
						Offset:          info.Offset + len(linkCandidate),
						NotificationURI: info.NotificationURI,
						EntropyValue:    info.EntropyValue,
					})
				} else {
					sub.Get(linkCandidate).Subset.Set(fs[len(linkCandidate):], &Info{
						Offset:          info.Offset + len(linkCandidate),
						NotificationURI: info.NotificationURI,
						EntropyValue:    info.EntropyValue,
					})
				}
				// recursively link the subset
				sub.Get(linkCandidate).Subset.linkSubset()
				// finaly delete the filter from the upper Subscriptions
				sub.Delete(fs)
			}
		}
	}
	return
}

func recalculateEntropyValue(sub Subscriptions) float64 {
	ev := float64(0)
	for _, fs := range sub.Keys() {
		ss := sub.Get(fs).Subset
		if ss.Length() != 0 {
			sub.Lock()
			sub.m[fs].EntropyValue += recalculateEntropyValue(ss)
			sub.Unlock()
		}
		ev += sub.Get(fs).EntropyValue
	}
	return ev
}

func (sub Subscriptions) print(writer io.Writer, indent int) {
	for _, fs := range sub.Keys() {
		fmt.Fprintf(writer, "%s--%s %f\n", strings.Repeat(" ", indent), fs, sub.Get(fs).EntropyValue)
		ss := sub.Get(fs).Subset
		if ss.Length() != 0 {
			ss.print(writer, indent+2)
		}
	}
}
