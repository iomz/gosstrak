package filter

import (
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

func (f *Filter) match(id []byte) bool {
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

// makeFilter returns padded filter and mask in rune slices
func makeFilter(bs []rune, offset int) (int, []rune, []rune) {
	var f []rune
	var m []rune
	var bodyLength int

	leftPaddingLength := offset % ByteLength
	// pad left with 1 if necessary
	if offset != 0 {
		f = []rune(strings.Repeat("1", leftPaddingLength))
		m = []rune(strings.Repeat("1", leftPaddingLength))
	}

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

	// if there is remainder, continue making the filter
	if len(bs)+leftPaddingLength > ByteLength {
		_, fc, mc := makeFilter(bs[bodyLength:], 0)
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
