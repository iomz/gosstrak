package filter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iomz/go-llrp/binutil"
)

const (
	// ByteLength defines a bit length of one byte
	ByteLength = 8
)

// Filter type is a struct for filter element
type Filter struct {
	String     string
	Size       int
	Offset     int
	ByteFilter []byte
	ByteMask   []byte
	ByteOffset int
	ByteSize   int
}

// GetByteAt returns a byte of ByteFilter and ByteMask
// at the given offset, returns error if HasByteAt(bo) is false
func (f *Filter) GetByteAt(bo int) (byte, byte, error) {
	if f.HasByteAt(bo) {
		for i := 0; i < f.ByteSize; i++ {
			if bo == f.ByteOffset+i {
				return f.ByteFilter[i], f.ByteMask[i], nil
			}
		}
	}
	return 0, 0, errors.New("this filter doesn't have the byte in the given offset")
}

// HasByteAt returns true if the ByteFilter covers
// a byte starting with the given offset
func (f *Filter) HasByteAt(bo int) bool {
	if f.ByteOffset > bo { // the filter is after the given bo
		return false
	} else if f.ByteOffset+f.ByteSize <= bo { // the filter is before the given offset
		return false
	}
	return true
}

// Match returns true if the id is captured by this filter
func (f *Filter) Match(id []byte) bool {
	for i := 0; i < f.ByteSize; i++ {
		if (id[f.ByteOffset+i]|f.ByteMask[i])^f.ByteFilter[i] != byte(0) {
			return false
		}
	}
	return true
}

// ToString returns a string representation of Filter
func (f *Filter) ToString() string {
	return fmt.Sprintf("%s(%d %d)", f.String, f.Offset, f.Size)
}

// makeFilter returns padded offset, filter and mask in rune slices
func makeFilter(bs []rune, offset int) (int, []rune, []rune) {
	var f []rune
	var m []rune
	var bodyLength int

	leftPaddingLength := offset % ByteLength
	// pad left with 1 if necessary
	f = []rune(strings.Repeat("1", leftPaddingLength))
	m = []rune(strings.Repeat("1", leftPaddingLength))

	// insert filter body
	if ByteLength-leftPaddingLength < len(bs) {
		bodyLength = ByteLength - leftPaddingLength
	} else {
		bodyLength = len(bs)
	}
	f = append(f, bs[0:bodyLength]...)
	m = append(m, []rune(strings.Repeat("0", bodyLength))...)

	// pad right with 1 if necessary
	if len(f) < ByteLength {
		rightPaddingLength := ByteLength - len(f)
		f = append(f, []rune(strings.Repeat("1", rightPaddingLength))...)
		m = append(m, []rune(strings.Repeat("1", rightPaddingLength))...)
	}

	// Apply wildcard x bits in the filter to mask
	for i := range f {
		if f[i] == 'x' {
			f[i] = '1'
			m[i] = '1'
		}
	}

	// if there is remainder, continue making the filter
	if len(bs)+leftPaddingLength > ByteLength {
		nextOffset := ((offset / ByteLength) + 1) * ByteLength
		_, fc, mc := makeFilter(bs[bodyLength:], nextOffset)
		f = append(f, fc...)
		m = append(m, mc...)
	}

	return offset / ByteLength, f, m
}

// NewFilter constructs Filter
func NewFilter(s string, o int) *Filter {
	bo, f, m := makeFilter([]rune(s), o)
	bf, _ := binutil.ParseBinRuneSliceToUint8Slice(f)
	bm, _ := binutil.ParseBinRuneSliceToUint8Slice(m)

	return &Filter{
		String:     s,
		Size:       len(s),
		Offset:     o,
		ByteFilter: bf,
		ByteMask:   bm,
		ByteOffset: bo,
		ByteSize:   len(bf),
	}
}
