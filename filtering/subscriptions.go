// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"sort"
)

// Subscriptions contains filter string as key and Info as value
type Subscriptions map[string]*Info

// Info contains notificationURI and pValue for a filter
type Info struct {
	NotificationURI string
	EntropyValue    float64
}

func (sub Subscriptions) HuffmanTable() HuffmanTable {
	ht := make(HuffmanTable, len(sub))
	i := 0
	for fs, info := range sub {
		ht[i] = HuffmanCode{[]string{fs}, info.EntropyValue}
		i++
	}
	ht.sort()
	return ht
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
