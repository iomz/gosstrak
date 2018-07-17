// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
)

// List is a slice of pointers to ExactMatch
type List []*ExactMatch

// ExactMatch is a raw filter directly taken from ByteSubscriptions
type ExactMatch struct {
	reportURI string
	filter    *FilterObject
}

// AddSubscription adds a set of subscriptions if not exists yet
func (list *List) AddSubscription(sub ByteSubscriptions) {
	// store ExactMatch in sorted order from sub
	for _, fs := range sub.Keys() {
		em := &ExactMatch{sub[fs].ReportURI, NewFilter(fs, sub[fs].Offset)}
		if list.IndexOf(em) < 0 {
			*list = append(*list, em)
		}
	}
}

// DeleteSubscription deletes a set of subscriptions if already exist
func (list *List) DeleteSubscription(sub ByteSubscriptions) {
	// store ExactMatch in sorted order from sub
	for _, fs := range sub.Keys() {
		em := &ExactMatch{sub[fs].ReportURI, NewFilter(fs, sub[fs].Offset)}
		if i := list.IndexOf(em); i > -1 {
			*list = append((*list)[:i], (*list)[i+1:]...)
		}
	}
}

// IndexOf check the index of ExactMatch in the List
// returns -1 if not exist
func (list *List) IndexOf(em *ExactMatch) int {
	for i, a := range *list {
		if reflect.DeepEqual(a, em) {
			return i
		}
	}
	return -1
}

// Dump returs a string representation of the PatriciaTrie
func (list *List) Dump() string {
	writer := &bytes.Buffer{}
	for _, em := range *list {
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
	enc.Encode(len(*list))
	for _, em := range *list {
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

// Search returns a slice of reportURI
func (list *List) Search(id []byte) (matches []string) {
	for _, em := range *list {
		if em.filter.Match(id) {
			matches = append(matches, em.reportURI)
		}
	}
	return
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding List
func (list *List) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.List" {
		return errors.New("Wrong Filtering Engine: " + typeOfEngine)
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
		*list = append(*list, &em)
	}

	return
}

// NewList builds a simple list of filters from filter.ByteSubscriptions
// returns the pointer to the slice of ExactMatch struct
func NewList(sub ByteSubscriptions) Engine {
	list := List{}

	// store ExactMatch in sorted order from sub
	for _, fs := range sub.Keys() {
		list = append(list, &ExactMatch{sub[fs].ReportURI, NewFilter(fs, 0)})
	}

	return &list
}
