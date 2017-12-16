package filter

import (
	"strings"

	"github.com/iomz/go-llrp/binutil"
)

// Composition contains list of composite filters
// and its origins
type Composition struct {
	Filter   *Filter
	Children ChildFilters
}

// ChildFilters contains processed filters in map
// with its original filter as a key
type ChildFilters map[string]*Filter

// NewComposition constructs Composition
func NewComposition(filters []*Filter) *Composition {
	// 0.) Find the edges
	// This is going to be the range of the composite filter
	// of size: tail-head
	head := -1 // ByteOffset
	tail := -1 // ByteOffset
	for _, f := range filters {
		if head == -1 || f.ByteOffset > head {
			head = f.ByteOffset
		}
		if tail == -1 || f.ByteOffset+f.ByteSize < tail {
			tail = f.ByteOffset + f.ByteSize
		}
	}

	// 1.) Create FilterDisjunction and PMASK
	// head-f.ByteOffset+i indicates ith byte of ByteFilter
	// where head: i == 0 and tail: i == tail-head-1
	disjunct := make([]byte, tail-head)
	pmask := make([]byte, tail-head)
	for _, f := range filters {
		for i := 0; i < tail-head; i++ {
			pmask[i] = pmask[i] | f.ByteMask[head-f.ByteOffset+i]
			disjunct[i] = disjunct[i] | f.ByteFilter[head-f.ByteOffset+i]
		}
	}

	// 2.) Create XMASK
	// capture all exclusive bits by XORing filters
	xmask := make([]byte, tail-head)
	f := filters[0]                     // first filter
	for i := 1; i < len(filters); i++ { // second+ filters
		ff := filters[i]
		for i := 0; i < tail-head; i++ {
			xm := f.ByteFilter[head-f.ByteOffset+i] ^ ff.ByteFilter[head-ff.ByteOffset+i]
			xmask[i] |= xm
		}
	}

	// 3.) Create Composition
	// check all the bytes in the compositeFilter
	// replace wildcard bits in the filter with x
	var cfString string
	for i := 0; i < tail-head; i++ {
		mask := pmask[i] | xmask[i]
		smask := []rune(binutil.ParseByteSliceToBinString([]byte{mask}))
		sdisjunct := []rune(binutil.ParseByteSliceToBinString([]byte{disjunct[i]}))
		for i, m := range smask {
			if m == '1' {
				sdisjunct[i] = 'x'
			}
		}
		cfString += string(sdisjunct)
	}
	compositeFilter := NewFilter(cfString, head*ByteLength)

	// 4.) Create ChildFilters
	// check all the filters with the compositeFilter
	// remove duplicates from the original filters
	// a duplicates byte is
	// cf.ByteFilter[i]^f.ByteFilter[head-f.ByteOffset+i] == 0 &&
	// cf.ByteMask[i]^f.ByteMask[head-f.ByteOffset+i] == 0
	children := ChildFilters{}
	for _, f := range filters {
		var nfString string
		var nfOffset int
		for i, b := range f.ByteFilter {
			o := f.ByteOffset + i
			if o < head { // f starts before head
				if len(nfString) == 0 { // nf starts from f.ByteOffset
					nfOffset = f.ByteOffset * ByteLength
				}
				nfString += processWildcard(b, f.ByteMask[i])
			} else if head <= o && o < tail { // overwrap range
				if len(nfString) == 0 { // nf starts from here, within cf
					nfOffset = o * ByteLength
				}
				if compositeFilter.ByteFilter[o-head]^b == 0 &&
					compositeFilter.ByteMask[o-head]^f.ByteMask[i] == 0 { // duplicate with cf
					nfString += strings.Repeat("x", ByteLength) // ignore this byte
				} else { // there is at least one exclusive bit; cannot ignore this byte
					nfString += processWildcard(b, f.ByteMask[i])
				}
			} else { // the byte is after cf, just add the remainders
				nfString += processWildcard(b, f.ByteMask[i])
			}
		}
		children[f.String] = NewFilter(nfString, nfOffset)
	}

	return &Composition{
		Filter:   compositeFilter,
		Children: children,
	}
}

// processWildcard takes a byte from both filter and mask
// replace wildcard bit with x
func processWildcard(b byte, m byte) string {
	sb := []rune(binutil.ParseByteSliceToBinString([]byte{b}))
	sm := binutil.ParseByteSliceToBinString([]byte{b & m})
	for ii, bb := range sm {
		if bb == '1' {
			sb[ii] = 'x'
		}
	}
	return string(sb)
}

// groupSequentialSlice group the sequential []int
// return the groups as a slice of slices
func groupSequentialSlice(rs []int) [][]int {
	var gs [][]int
	var g []int
	for _, r := range rs {
		if len(g) == 0 {
			g = []int{r}
		} else if g[len(g)-1]+1 == r {
			g = append(g, r)
		} else {
			gs = append(gs, g)
			g = []int{r}
		}
	}
	gs = append(gs, g)
	return gs
}
