package filter

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"strings"
)

type PatriciaTrie struct {
	Notify string
	Filter *Filter
	One    *PatriciaTrie
	Zero   *PatriciaTrie
}

func (pt *PatriciaTrie) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Notify
	enc.Encode(pt.Notify)

	// Filter
	hasFilter := pt.Filter != nil
	enc.Encode(hasFilter)
	if hasFilter {
		err = enc.Encode(pt.Filter)
	}

	// One
	hasOne := pt.One != nil
	enc.Encode(hasOne)
	if hasOne {
		err = enc.Encode(pt.One)
	}

	// Zero
	hasZero := pt.Zero != nil
	enc.Encode(hasZero)
	if hasZero {
		err = enc.Encode(pt.Zero)
	}

	//buf.Encode
	return buf.Bytes(), err
}

func (pt *PatriciaTrie) UnmarshalBinary(data []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))

	// Notify
	if err = dec.Decode(&pt.Notify); err != nil {
		return
	}

	// Filter
	var hasFilter bool
	if err = dec.Decode(&hasFilter); err != nil {
		return
	}
	if hasFilter {
		err = dec.Decode(&pt.Filter)
	} else {
		pt.Filter = nil
	}

	// One
	var hasOne bool
	if err = dec.Decode(&hasOne); err != nil {
		return
	}
	if hasOne {
		err = dec.Decode(&pt.One)
	} else {
		pt.One = nil
	}

	// Zero
	var hasZero bool
	if err = dec.Decode(&hasZero); err != nil {
		return
	}
	if hasOne {
		err = dec.Decode(&pt.Zero)
	} else {
		pt.Zero = nil
	}

	return
}

func (pt *PatriciaTrie) Match(id []byte, matched *[]string) {
	// if not match, return empty string immediately
	if !pt.Filter.match(id) {
		return
	}

	// if the id matched with this node, return notify
	if len(pt.Notify) != 0 {
		*matched = append(*matched, pt.Notify)
	}

	// Determine next Filter
	next_bit_offset := pt.Filter.Offset + pt.Filter.Size
	nb, err := get_next_bit(id, next_bit_offset)
	if err != nil {
		panic(err)
	}
	if nb == '1' && pt.One != nil {
		pt.One.Match(id, matched)
	} else if nb == '0' && pt.Zero != nil {
		pt.Zero.Match(id, matched)
	}
}

func (pt *PatriciaTrie) constructTrie(prefix string, fm FilterMap) {
	onePrefixBranch := ""
	zeroPrefixBranch := ""
	fks := fm.keys()
	for i := 0; i < len(fks); i++ {
		// if the prefix is already longer than the testee
		if len(fks[i]) < len(prefix) {
			continue
		}
		// ignore the testee without the prefix
		if !strings.HasPrefix(fks[i], prefix) {
			continue
		}
		p := fks[i][len(prefix):]
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
		pt.One = &PatriciaTrie{}
		pt.One.Filter = NewFilter(onePrefixBranch, len(prefix))
		cumulativePrefix = prefix + onePrefixBranch
		// check if the prefix matches whole filter
		if n, ok := fm[cumulativePrefix]; ok {
			pt.One.Notify = n
		}
		pt.One.constructTrie(cumulativePrefix, fm)
	}
	// if there's a branch starts with 0
	if len(zeroPrefixBranch) != 0 {
		pt.Zero = &PatriciaTrie{}
		pt.Zero.Filter = NewFilter(zeroPrefixBranch, len(prefix))
		cumulativePrefix = prefix + zeroPrefixBranch
		// check if the prefix matches whole filter
		if n, ok := fm[cumulativePrefix]; ok {
			pt.Zero.Notify = n
		}
		pt.Zero.constructTrie(cumulativePrefix, fm)
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
	if len(pt.Notify) != 0 {
		n = "-> " + pt.Notify
	}
	fmt.Fprintf(writer, "%s--%s %s\n", strings.Repeat(" ", indent), pt.Filter.ToString(), n)
	if pt.One != nil {
		pt.One.print(writer, indent+2)
	}
	if pt.Zero != nil {
		pt.Zero.print(writer, indent+2)
	}
}

func BuildPatriciaTrie(fm FilterMap) *PatriciaTrie {
	p1 := lcp(fm.keys())
	if len(p1) == 0 {
		// do something if there's no common prefix
	}
	head := &PatriciaTrie{}
	head.Filter = NewFilter(p1, 0)
	head.constructTrie(p1, fm)

	return head
}

func get_next_bit(id []byte, nbo int) (rune, error) {
	o := nbo / ByteLength
	// No more bit in the ID
	if len(id) == o {
		return 'x', nil
	}
	if len(id) < o {
		return '?', errors.New("get_next_bit error")
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
