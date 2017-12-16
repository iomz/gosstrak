package filter

import (
	"github.com/iomz/go-llrp/binutil"
)

// CompositeFilter contains list of composite filters
// and its origins
type CompositeFilter struct {
	Filters  []*Filter
	Elements ElementFilters
}

// ElementFilters contains processed filters in map
// with its original filter as a key
type ElementFilters map[string][]*Filter

// Match matches with this composite filter
func (f *CompositeFilter) Match(id []byte) bool {
	/*
		for i := 0; i < f.ByteSize; i++ {
			if (id[f.ByteOffset+i]|f.ByteMask[i])^f.ByteFilter[i] != byte(0) {
				return false
			}
		}
	*/
	return true
}

// NewCompositeFilter constructs CompositeFilter
func NewCompositeFilter(filters []*Filter) *CompositeFilter {
	// Find the edges
	head := -1 // ByteOffset
	tail := -1 // ByteOffset
	for _, f := range filters {
		if head == -1 {
			head = f.ByteOffset
		} else if f.ByteOffset > head {
			head = f.ByteOffset
		}
		if tail == -1 {
			tail = f.ByteOffset + f.ByteSize
		} else if f.ByteOffset+f.ByteSize < tail {
			tail = f.ByteOffset + f.ByteSize
		}
	}

	// Create CF and PMASK
	cf := make([]byte, tail-head)
	pmask := make([]byte, tail-head)
	for _, f := range filters {
		for i := 0; i < tail-head; i++ {
			pmask[i] = pmask[i] | f.ByteMask[head-f.ByteOffset+i]
			cf[i] = cf[i] | f.ByteFilter[head-f.ByteOffset+i]
		}
	}

	// Create XMASK
	xmask := make([]byte, tail-head)
	for _, f := range filters {
		for _, ff := range filters {
			for i := 0; i < tail-head; i++ {
				xm := f.ByteFilter[head-f.ByteOffset+i] ^ ff.ByteFilter[head-f.ByteOffset+i]
				if xm == byte(0) {
					continue
				}
				xmask[i] = xmask[i] | xm
			}
		}
	}

	// Create CompositeFilters
	cfs := []*Filter{}
	cfByteOffsets := []int{}
	for i := 0; i < tail-head; i++ {
		if pmask[i]|xmask[i] == byte(255) { // if mask is 0xFF, discard cf
			continue
		}
		cfByteOffsets = append(cfByteOffsets, i)
	}
	cfByteGroup := groupSequentialSlice(cfByteOffsets)
	for _, group := range cfByteGroup {
		var ncfString string
		var ncfOffset int                                // BitOffset
		for positionInGroup, byteOffset := range group { // [0]
			if positionInGroup == 0 {
				ncfOffset = (head + byteOffset) * ByteLength
			}
			mask := pmask[byteOffset] | xmask[byteOffset]
			smask := []rune(binutil.ParseByteSliceToBinString([]byte{mask}))
			sncf := []rune(binutil.ParseByteSliceToBinString([]byte{cf[byteOffset]}))
			for i, m := range smask {
				if m == '1' {
					sncf[i] = 'x'
				}
			}
			ncfString += string(sncf)
		}
		ncf := NewFilter(ncfString, ncfOffset)
		cfs = append(cfs, ncf)
	}

	// Create ElementFilters
	efs := ElementFilters{}
	for _, f := range filters {
		// Create a new filters subtracing the CF overwrap from the original
		nfByteOffsets := []int{}
		for i, bf := range f.ByteFilter {
			for _, cf := range cfs {
				if cf.HasByteAt(f.ByteOffset + i) {
					cfbf, cfbm, _ := cf.GetByteAt(f.ByteOffset + i)
					if bf^cfbf != 0 || f.ByteMask[i]^cfbm != 0 {
						nfByteOffsets = append(nfByteOffsets, i)
					}
				} else {
					nfByteOffsets = append(nfByteOffsets, i)
				}
			}
		}
		nfByteGroup := groupSequentialSlice(nfByteOffsets)
		// + a new filters
		nfs := []*Filter{}
		for _, nfGroup := range nfByteGroup { // [[1]]
			var nfOffset int
			var nfString string
			for positionInGroup, nfByteOffset := range nfGroup { //[1]
				if positionInGroup == 0 { // first in the group
					if nfByteOffset == 0 { // first byte in the original filter
						nfOffset = f.Offset
						nfString = f.String[:ByteLength-(f.Offset%ByteLength)]
					} else { // not first byte in the original filter
						nfOffset = (f.ByteOffset + nfByteOffset) * ByteLength
						nfBitOffset := nfByteOffset * ByteLength
						if f.Size-nfBitOffset < ByteLength {
							nfString = f.String[nfBitOffset : nfBitOffset+(f.Size-nfBitOffset)%ByteLength]
						} else {
							nfString = f.String[nfBitOffset : nfBitOffset+ByteLength]
						}
					}
					/*
						if f.ByteOffset+nfByteOffset < f.ByteSize { // filter is just within this byte
							nfString = f.String
						} else { // has next byte, until the end of this byte
							nfString = f.String[:ByteLength-(f.Offset%ByteLength)]
						}
					*/
					// test comes until here ------------------------------

				} else { // not first in the group
					bitOffset := (f.ByteOffset + nfByteOffset) * ByteLength
					if nfByteOffset+1 == f.ByteSize { // last byte in the original filter
						nfString += f.String[bitOffset:]
					} else {
						nfString += f.String[bitOffset : bitOffset+ByteLength]
					}
				}
				// test ---------------------------------------------------
			}
			nf := NewFilter(nfString, nfOffset)
			nfs = append(nfs, nf)
		}
		efs[f.String] = nfs
	}

	return &CompositeFilter{
		Filters:  cfs,
		Elements: efs,
	}
}

// groupSequentialSlice group the sequential []int and return as slice of slice
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
