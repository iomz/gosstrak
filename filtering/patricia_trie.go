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
	"strings"
)

// PatriciaTrie struct
type PatriciaTrie struct {
	notificationURI string
	filter          *FilterObject
	one             *PatriciaTrie
	zero            *PatriciaTrie
}

// MarshalBinary overwrites the marshaller in gob encoding *PatriciaTrie
func (pt *PatriciaTrie) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Notify
	enc.Encode(pt.notificationURI)

	// Filter
	hasFilter := pt.filter != nil
	enc.Encode(hasFilter)
	if hasFilter {
		err = enc.Encode(pt.filter)
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

// UnmarshalBinary overwrites the unmarshaller in gob decoding *PatriciaTrie
func (pt *PatriciaTrie) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

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
		err = dec.Decode(&pt.filter)
	} else {
		pt.filter = nil
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

// AnalyzeLocality increments the locality per node for the specific id
func (pt *PatriciaTrie) AnalyzeLocality(id []byte, prefix string, lm *LocalityMap) {
	// if not match, return empty string immediately
	if !pt.filter.Match(id) {
		return
	}

	// if the node is the first two node
	if len(pt.filter.String) != 0 {
		prefix += "-" + pt.filter.String
	}

	if _, ok := (*lm)[prefix]; !ok {
		(*lm)[prefix] = 1
	} else {
		(*lm)[prefix]++
	}

	// Determine next filter
	nextBitOffset := pt.filter.Offset + pt.filter.Size
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

// Search returns a slice of notificationURI
func (pt *PatriciaTrie) Search(id []byte) (matches []string) {
	// if not match, return empty slice immediately
	if !pt.filter.Match(id) {
		return
	}

	// if the id matched with this node, return notificationURI
	if len(pt.notificationURI) != 0 {
		matches = append(matches, pt.notificationURI)
	}

	// Determine next filter
	nextBitOffset := pt.filter.Offset + pt.filter.Size
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

func (pt *PatriciaTrie) build(prefix string, sub Subscriptions) {
	onePrefixBranch := ""
	zeroPrefixBranch := ""
	fks := sub.keys()
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
		pt.one.filter = NewFilter(onePrefixBranch, len(prefix))
		cumulativePrefix = prefix + onePrefixBranch
		// check if the prefix matches whole filter
		if info, ok := sub[cumulativePrefix]; ok {
			pt.one.notificationURI = info.NotificationURI
		}
		pt.one.build(cumulativePrefix, sub)
	}
	// if there's a branch starts with 0
	if len(zeroPrefixBranch) != 0 {
		pt.zero = &PatriciaTrie{}
		pt.zero.filter = NewFilter(zeroPrefixBranch, len(prefix))
		cumulativePrefix = prefix + zeroPrefixBranch
		// check if the prefix matches whole filter
		if info, ok := sub[cumulativePrefix]; ok {
			pt.zero.notificationURI = info.NotificationURI
		}
		pt.zero.build(cumulativePrefix, sub)
	}
}

// Dump returs a string representation of the PatriciaTrie
func (pt *PatriciaTrie) Dump() string {
	writer := &bytes.Buffer{}
	pt.print(writer, 0)
	return writer.String()
}

func (pt *PatriciaTrie) print(writer io.Writer, indent int) {
	var n string
	if len(pt.notificationURI) != 0 {
		n = "-> " + pt.notificationURI
	}
	fmt.Fprintf(writer, "%s--%s %s\n", strings.Repeat(" ", indent), pt.filter.ToString(), n)
	if pt.one != nil {
		pt.one.print(writer, indent+2)
	}
	if pt.zero != nil {
		pt.zero.print(writer, indent+2)
	}
}

// BuildPatriciaTrie builds PatriciaTrie from filter.Subscriptions
// returns the pointer to the entry node
func BuildPatriciaTrie(sub Subscriptions) *PatriciaTrie {
	p1 := lcp(sub.keys())
	if len(p1) == 0 {
		// do something if there's no common prefix
	}
	head := &PatriciaTrie{}
	head.filter = NewFilter(p1, 0)
	head.build(p1, sub)

	return head
}

func getNextBit(id []byte, nbo int) (rune, error) {
	o := nbo / ByteLength
	// No more bit in the ID
	if len(id) == o {
		return 'x', nil
	}
	if len(id) < o {
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
