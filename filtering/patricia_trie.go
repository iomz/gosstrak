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
	"io"
	"reflect"
	"strings"
)

// PatriciaTrie struct
type PatriciaTrie struct {
	notificationURI string
	filterObject    *FilterObject
	one             *PatriciaTrie
	zero            *PatriciaTrie
}

// AddSubscription adds a set of subscriptions if not exists yet
func (pt *PatriciaTrie) AddSubscription(sub Subscriptions) {
	for _, fs := range sub.Keys() {
		pt.add(fs, sub.Get(fs).NotificationURI)
	}
}

// AnalyzeLocality increments the locality per node for the specific id
func (pt *PatriciaTrie) AnalyzeLocality(id []byte, prefix string, lm *LocalityMap) {
	// if not match, return empty string immediately
	if !pt.filterObject.Match(id) {
		return
	}

	// if the node is the first two node
	if len(pt.filterObject.String) != 0 {
		prefix += "," + pt.filterObject.String
	}

	if _, ok := (*lm)[prefix]; !ok {
		(*lm)[prefix] = 1
	} else {
		(*lm)[prefix]++
	}

	// Determine next filter
	nextBitOffset := pt.filterObject.Offset + pt.filterObject.Size
	nb, err := getNextBit(id, nextBitOffset)
	if err != nil {
		panic(err)
	}
	if nb == '1' && pt.one != nil {
		pt.one.AnalyzeLocality(id, prefix, lm)
	} else if nb == '0' && pt.zero != nil {
		pt.zero.AnalyzeLocality(id, prefix, lm)
	}
}

// DeleteSubscription deletes a set of subscriptions if already exist
func (pt *PatriciaTrie) DeleteSubscription(sub Subscriptions) {
	for _, fs := range sub.Keys() {
		pt.delete(fs, sub.Get(fs).NotificationURI)
	}
}

// Dump returs a string representation of the PatriciaTrie
func (pt *PatriciaTrie) Dump() string {
	writer := &bytes.Buffer{}
	pt.print(writer, 0)
	return writer.String()
}

// MarshalBinary overwrites the marshaller in gob encoding *PatriciaTrie
func (pt *PatriciaTrie) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Type of Engine
	enc.Encode("Engine:filtering.PatriciaTrie")

	// Notify
	enc.Encode(pt.notificationURI)

	// Filter
	hasFilter := pt.filterObject != nil
	enc.Encode(hasFilter)
	if hasFilter {
		err = enc.Encode(pt.filterObject)
	}

	// One
	hasOne := pt.one != nil
	enc.Encode(hasOne)
	if hasOne {
		err = enc.Encode(pt.one)
	}

	// Zero
	hasZero := pt.zero != nil
	enc.Encode(hasZero)
	if hasZero {
		err = enc.Encode(pt.zero)
	}

	//buf.Encode
	return buf.Bytes(), err
}

// Search returns a slice of notificationURI
func (pt *PatriciaTrie) Search(id []byte) (matches []string) {
	// if not match, return empty slice immediately
	if !pt.filterObject.Match(id) {
		return
	}

	// if the id matched with this node, return notificationURI
	if len(pt.notificationURI) != 0 {
		matches = append(matches, pt.notificationURI)
	}

	// Determine next filter
	nextBitOffset := pt.filterObject.Offset + pt.filterObject.Size
	nb, err := getNextBit(id, nextBitOffset)
	if err != nil {
		panic(err)
	}
	if nb == '1' && pt.one != nil {
		matches = append(matches, pt.one.Search(id)...)
	} else if nb == '0' && pt.zero != nil {
		matches = append(matches, pt.zero.Search(id)...)
	}
	return
}

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (pt *PatriciaTrie) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Type of Engine
	var typeOfEngine string
	if err = dec.Decode(&typeOfEngine); err != nil || typeOfEngine != "Engine:filtering.PatriciaTrie" {
		return errors.New("Wrong Filtering Engine: " + typeOfEngine)
	}

	// Notify
	if err = dec.Decode(&pt.notificationURI); err != nil {
		return
	}

	// Filter
	var hasFilter bool
	if err = dec.Decode(&hasFilter); err != nil {
		return
	}
	if hasFilter {
		err = dec.Decode(&pt.filterObject)
	} else {
		pt.filterObject = nil
	}

	// One
	var hasOne bool
	if err = dec.Decode(&hasOne); err != nil {
		return
	}
	if hasOne {
		err = dec.Decode(&pt.one)
	} else {
		pt.one = nil
	}

	// Zero
	var hasZero bool
	if err = dec.Decode(&hasZero); err != nil {
		return
	}
	if hasOne {
		err = dec.Decode(&pt.zero)
	} else {
		pt.zero = nil
	}

	return
}

// Internal helper methods -----------------------------------------------------

// add a subscription if not exist yet
func (pt *PatriciaTrie) add(fs string, notificationURI string) {
	if strings.HasPrefix(fs, pt.filterObject.String) { // fs \in pt.FilterObject.String
		if fs == pt.filterObject.String { // the identical filter
			// if the notificationURI is different, update it
			if pt.notificationURI != notificationURI {
				pt.notificationURI = notificationURI
			}
			return //end
		}
		//} else if len(fs) < pt.filterObject.Size { // Needs a reconstruction
	} else {
		newCommonPrefix := lcp([]string{fs, pt.filterObject.String})
		ncpLength := len(newCommonPrefix)
		newNode := &PatriciaTrie{}
		newNode.filterObject = NewFilter(pt.filterObject.String[ncpLength:], pt.filterObject.Offset+ncpLength)
		newNode.one = pt.one
		newNode.zero = pt.zero
		newNode.notificationURI = pt.notificationURI
		pt.notificationURI = ""
		currentOffset := pt.filterObject.Offset
		pt.filterObject = NewFilter(newCommonPrefix, currentOffset)
		switch fs[ncpLength] {
		case '1':
			pt.zero = newNode
			pt.one = &PatriciaTrie{}
			pt.one.filterObject = NewFilter(fs[ncpLength:], currentOffset+ncpLength)
			pt.one.notificationURI = notificationURI
		case '0':
			pt.one = newNode
			pt.zero = &PatriciaTrie{}
			pt.zero.filterObject = NewFilter(fs[ncpLength:], currentOffset+ncpLength)
			pt.zero.notificationURI = notificationURI
		}
		return //end
	}

	// If there's remainder
	if len(fs) > pt.filterObject.Size {
		switch fs[pt.filterObject.Size] {
		case '1':
			if pt.one == nil {
				pt.one = &PatriciaTrie{}
				pt.one.filterObject = NewFilter(fs[pt.filterObject.Size:], pt.filterObject.Offset+pt.filterObject.Size)
				pt.one.notificationURI = notificationURI
				return //end
			}
			pt.one.add(fs[pt.filterObject.Size:], notificationURI)
		case '0':
			if pt.zero == nil {
				pt.zero = &PatriciaTrie{}
				pt.zero.filterObject = NewFilter(fs[pt.filterObject.Size:], pt.filterObject.Offset+pt.filterObject.Size)
				pt.zero.notificationURI = notificationURI
				return //end
			}
			pt.zero.add(fs[pt.filterObject.Size:], notificationURI)
		}
	}
}

func (pt *PatriciaTrie) build(prefix string, sub Subscriptions) {
	onePrefixBranch := ""
	zeroPrefixBranch := ""
	fks := sub.Keys()
	for _, fk := range fks {
		// if the prefix is already longer than the testee
		if len(fk) < len(prefix) {
			continue
		}
		// ignore the testee without the prefix
		if !strings.HasPrefix(fk, prefix) {
			continue
		}
		p := fk[len(prefix):]
		// ignore if no remainder
		if len(p) == 0 {
			continue
		}
		// if the remainder starts with 1
		if strings.HasPrefix(p, "1") {
			if len(onePrefixBranch) == 0 {
				onePrefixBranch = p
			} else {
				onePrefixBranch = lcp([]string{p, onePrefixBranch})
			}
			// if the remainder starts with 0
		} else if strings.HasPrefix(p, "0") {
			if len(zeroPrefixBranch) == 0 {
				zeroPrefixBranch = p
			} else {
				zeroPrefixBranch = lcp([]string{p, zeroPrefixBranch})
			}
		}
	}
	cumulativePrefix := ""
	// if there's a branch starts with 1
	if len(onePrefixBranch) != 0 {
		pt.one = &PatriciaTrie{}
		pt.one.filterObject = NewFilter(onePrefixBranch, len(prefix))
		cumulativePrefix = prefix + onePrefixBranch
		// check if the prefix matches whole filter
		if sub.Has(cumulativePrefix) {
			pt.one.notificationURI = sub.Get(cumulativePrefix).NotificationURI
		}
		pt.one.build(cumulativePrefix, sub)
	}
	// if there's a branch starts with 0
	if len(zeroPrefixBranch) != 0 {
		pt.zero = &PatriciaTrie{}
		pt.zero.filterObject = NewFilter(zeroPrefixBranch, len(prefix))
		cumulativePrefix = prefix + zeroPrefixBranch
		// check if the prefix matches whole filter
		if sub.Has(cumulativePrefix) {
			pt.zero.notificationURI = sub.Get(cumulativePrefix).NotificationURI
		}
		pt.zero.build(cumulativePrefix, sub)
	}
}

// delete a subscriptions if already exist
func (pt *PatriciaTrie) delete(fs string, notificationURI string) {
	// No such filter exist in the trie
	if !strings.HasPrefix(fs, pt.filterObject.String) {
		return
	}

	// This is the filter to delete
	if fs == pt.filterObject.String {
		if pt.one != nil && pt.zero != nil { // node in the middle
			pt.notificationURI = ""
		} else if pt.one != nil { // has only one node
			newFilter := NewFilter(pt.filterObject.String+pt.one.filterObject.String, pt.filterObject.Offset)
			pt.filterObject = newFilter
			pt.notificationURI = pt.one.notificationURI
			pt.zero = pt.one.zero
			pt.one = pt.one.one
		} else if pt.zero != nil { // has only zero node
			newFilter := NewFilter(pt.filterObject.String+pt.zero.filterObject.String, pt.filterObject.Offset)
			pt.filterObject = newFilter
			pt.notificationURI = pt.zero.notificationURI
			pt.one = pt.zero.one
			pt.zero = pt.zero.zero
		}
		return //end
	}

	// If there's remainder
	if len(fs) > pt.filterObject.Size {
		switch fs[pt.filterObject.Size] {
		case '1':
			if pt.one != nil {
				if fs[pt.filterObject.Size:] == pt.one.filterObject.String &&
					pt.one.one == nil && pt.one.zero == nil {
					pt.one = nil
				} else {
					pt.one.delete(fs[pt.filterObject.Size:], notificationURI)
				}
			}
		case '0':
			if pt.zero != nil {
				if fs[pt.filterObject.Size:] == pt.zero.filterObject.String &&
					pt.zero.one == nil && pt.zero.zero == nil {
					pt.zero = nil
				} else {
					pt.zero.delete(fs[pt.filterObject.Size:], notificationURI)
				}
			}
		}
	}
	return
}

func (pt *PatriciaTrie) print(writer io.Writer, indent int) {
	var n string
	if len(pt.notificationURI) != 0 {
		n = "-> " + pt.notificationURI
	}
	fmt.Fprintf(writer, "%s--%s %s\n", strings.Repeat(" ", indent), pt.filterObject.ToString(), n)
	if pt.one != nil {
		pt.one.print(writer, indent+2)
	}
	if pt.zero != nil {
		pt.zero.print(writer, indent+2)
	}
}

func (pt *PatriciaTrie) equal(want *PatriciaTrie) (ok bool, got *PatriciaTrie, wanted *PatriciaTrie) {
	if pt.notificationURI != want.notificationURI ||
		!reflect.DeepEqual(pt.filterObject, want.filterObject) {
		return false, pt, want
	}
	if want.one != nil {
		if pt.one == nil {
			return false, nil, want.one
		}
		res, cgot, cwanted := pt.one.equal(want.one)
		if !res {
			return res, cgot, cwanted
		}
	}
	if want.zero != nil {
		if pt.zero == nil {
			return false, nil, want.zero
		}
		res, cgot, cwanted := pt.zero.equal(want.zero)
		if !res {
			return res, cgot, cwanted
		}
	}
	return true, nil, nil
}

func getNextBit(id []byte, nbo int) (rune, error) {
	o := nbo / ByteLength
	// No more bit in the ID
	if len(id) == o && nbo%ByteLength == 0 {
		return 'x', nil
	}
	if len(id) <= o {
		return '?', errors.New("getNextBit error")
	}
	if (uint8(id[o])>>uint8((7-(nbo%ByteLength))))%2 == 0 {
		return '0', nil
	}
	return '1', nil
}

// lcp finds the longest common prefix of the input strings.
// It compares by bytes instead of runes (Unicode code points).
// It's up to the caller to do Unicode normalization if desired
// (e.g. see golang.org/x/text/unicode/norm).
func lcp(l []string) string {
	// Special cases first
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0]
	}
	// LCP of min and max (lexigraphically)
	// is the LCP of the whole set.
	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}

// NewPatriciaTrie builds PatriciaTrie from filter.Subscriptions
// returns the pointer to the node
func NewPatriciaTrie(s Subscriptions) Engine {
	// copy subscription
	sub := s.Clone()
	p1 := lcp(sub.Keys())
	if len(p1) == 0 {
		// do something if there's no common prefix
	}
	head := &PatriciaTrie{}
	head.filterObject = NewFilter(p1, 0)
	head.build(p1, sub)

	return head
}
