package filter

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type PatriciaTrie struct {
	filter *Filter
	one    *PatriciaTrie
	zero   *PatriciaTrie
	notify string
}

func (pt *PatriciaTrie) Match(id []byte) {
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
		pt.one = &PatriciaTrie{}
		pt.one.filter = NewFilter(onePrefixBranch, len(prefix))
		cumulativePrefix = prefix + onePrefixBranch
		// check if the prefix matches whole filter
		if n, ok := fm[cumulativePrefix]; ok {
			pt.one.notify = n
		}
		pt.one.constructTrie(cumulativePrefix, fm)
	}
	// if there's a branch starts with 0
	if len(zeroPrefixBranch) != 0 {
		pt.zero = &PatriciaTrie{}
		pt.zero.filter = NewFilter(zeroPrefixBranch, len(prefix))
		cumulativePrefix = prefix + zeroPrefixBranch
		// check if the prefix matches whole filter
		if n, ok := fm[cumulativePrefix]; ok {
			pt.zero.notify = n
		}
		pt.zero.constructTrie(cumulativePrefix, fm)
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
	if len(pt.notify) != 0 {
		n = "-> " + pt.notify
	}
	fmt.Fprintf(writer, "%s--%s %s\n", strings.Repeat(" ", indent), pt.filter.stringFilter, n)
	if pt.one != nil {
		pt.one.print(writer, indent+2)
	}
	if pt.zero != nil {
		pt.zero.print(writer, indent+2)
	}
}

func BuildPatriciaTrie(fm FilterMap) *PatriciaTrie {
	p1 := lcp(fm.keys())
	if len(p1) == 0 {
		// do something if there's no common prefix
	}
	head := &PatriciaTrie{}
	head.filter = NewFilter(p1, 0)
	head.constructTrie(p1, fm)

	return head
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
