// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"

	"github.com/iomz/go-llrp"
	"github.com/iomz/gosstrak/tdt"
)

// List is a slice of pointers to ExactMatch
type List struct {
	filters ListFilters
	tdtCore *tdt.Core
}

// ListFilters contains pointers to ExactMatch
type ListFilters []*ExactMatch

// ExactMatch is a raw filter directly taken from ByteSubscriptions
type ExactMatch struct {
	filter    *FilterObject
	reportURI string
}

// AddSubscription adds a set of subscriptions if not exists yet
func (list *List) AddSubscription(sub Subscriptions) {
	bsub := sub.ToByteSubscriptions()
	// store ExactMatch in sorted order from sub
	for _, fs := range bsub.Keys() {
		em := &ExactMatch{
			filter:    NewFilter(fs, bsub[fs].Offset),
			reportURI: bsub[fs].ReportURI,
		}
		if list.filters.IndexOf(em) < 0 {
			list.filters = append(list.filters, em)
		}
	}
}

// DeleteSubscription deletes a set of subscriptions if already exist
func (list *List) DeleteSubscription(sub Subscriptions) {
	bsub := sub.ToByteSubscriptions()
	// store ExactMatch in sorted order from sub
	for _, fs := range bsub.Keys() {
		em := &ExactMatch{
			filter:    NewFilter(fs, bsub[fs].Offset),
			reportURI: bsub[fs].ReportURI,
		}
		if i := list.filters.IndexOf(em); i > -1 {
			if i+1 == len(list.filters) {
				list.filters = list.filters[:i]
			} else {
				list.filters = append(list.filters[:i], list.filters[i+1:]...)
			}
		}
	}
}

// IndexOf check the index of ExactMatch in the List
// returns -1 if not exist
func (lf ListFilters) IndexOf(em *ExactMatch) int {
	for i, a := range lf {
		if reflect.DeepEqual(a, em) {
			return i
		}
	}
	return -1
}

// Dump returs a string representation of the PatriciaTrie
func (list *List) Dump() string {
	writer := &bytes.Buffer{}
	for _, em := range list.filters {
		fmt.Fprintf(writer, "--%s %s\n", em.filter.ToString(), em.reportURI)
	}
	return writer.String()
}

// MarshalBinary overwrites the marshaller in gob encoding *List
func (list *List) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.List")

	// Size of List
	enc.Encode(len(list.filters))
	for _, em := range list.filters {
		// Notify
		enc.Encode(em.reportURI)
		// Filter
		err = enc.Encode(em.filter)
	}

	return buf.Bytes(), err
}

// Name returs the name of this engine type
func (list *List) Name() string {
	return "List"
}

// Search returns a pureIdentity of the llrp.ReadEvent if found any subscription without err
func (list *List) Search(re llrp.ReadEvent) (pureIdentity string, reportURIs []string, err error) {
	for _, em := range list.filters {
		if em.filter.Match(re.ID) {
			reportURIs = append(reportURIs, em.reportURI)
		}
	}
	if len(reportURIs) == 0 {
		return pureIdentity, reportURIs, fmt.Errorf("no match found for %v", re.ID)
	}
	pureIdentity, err = list.tdtCore.Translate(re.PC, re.ID)
	return
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding List
func (list *List) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.List" {
		return fmt.Errorf("Wrong Filtering Engine: %s", typeOfEngine)
	}

	// Size of List
	var listSize int
	if err = dec.Decode(&listSize); err != nil {
		return
	}

	for i := 0; i < listSize; i++ {
		em := ExactMatch{}
		// Notify
		if err = dec.Decode(&em.reportURI); err != nil {
			return
		}
		// Filter
		err = dec.Decode(&em.filter)
		list.filters = append(list.filters, &em)
	}

	// tdt.Core
	list.tdtCore = tdt.NewCore()

	return
}

// NewList builds a simple list of filters from filter.ByteSubscriptions
// returns the pointer to the slice of ExactMatch struct
func NewList(sub Subscriptions) Engine {
	list := &List{}

	// preprocess the subscriptions
	bsub := sub.ToByteSubscriptions()

	// store ExactMatch in sorted order from sub
	for _, fs := range bsub.Keys() {
		list.filters = append(list.filters, &ExactMatch{
			filter:    NewFilter(fs, 0),
			reportURI: bsub[fs].ReportURI,
		})
	}

	// initialize the tdt.Core
	list.tdtCore = tdt.NewCore()
	return list
}
