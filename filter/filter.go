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
	stringFilter string
	filterSize   int
	offset       int
	byteFilter   []byte
	byteMask     []byte
	paddedOffset int
	checkSize    int
}

func (f *Filter) match(id []byte) bool {
	for i := 0; i < f.checkSize; i++ {
		if (id[f.paddedOffset+i]|f.byteMask[i])^f.byteFilter[i] != byte(0) {
			return false
		}
	}
	return true
}

// ToString returns a string representation of Filter
func (f *Filter) ToString() string {
	return fmt.Sprintf("%s(%d %d)", f.stringFilter, f.offset, f.filterSize)
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
func NewFilter(sf string, o int) *Filter {
	po, f, m := makeFilter([]rune(sf), o)
	bf, _ := binutil.ParseBinRuneSliceToUint8Slice(f)
	bm, _ := binutil.ParseBinRuneSliceToUint8Slice(m)

	return &Filter{
		stringFilter: sf,
		filterSize:   len(sf),
		offset:       o,
		byteFilter:   bf,
		byteMask:     bm,
		paddedOffset: po,
		checkSize:    len(bf),
	}
}
