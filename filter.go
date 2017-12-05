package main

import (
	"io"
	"strings"
)

const (
	// DefaultMaxChildren defines a default number of children in a node
	DefaultMaxChildren = 10
	// ByteLength defines a bit length of one byte
	ByteLength = 8
)

// Filter type is a struct for filter element
type Filter struct {
	stringFilter string
	offset       int
	byteFilter   []byte
	paddedOffset int
}

// newFilter creates a new Filter struct
func newFilter(fs string, o int) *Filter {
	return &Filter{
		stringFilter: "",
		offset:       0,
		byteFilter:   []byte{},
		paddedOffset: 0,
	}
}

// makeFilter returns padded filter and mask in rune slices
func makeFilter(bs []rune, offset int) ([]rune, []rune) {
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
		fc, mc := makeFilter(bs[bodyLength:], 0)
		f = append(f, fc...)
		m = append(m, mc...)
	}

	return f, m
}

// FilterList defines a set of tries
type FilterList struct {
	children tries
}

// newFilterList creates a new set of FilterList
func newFilterList(maxChildren int) *FilterList {
	if maxChildren == 0 {
		maxChildren = DefaultMaxChildren
	}
	return &FilterList{
		children: make(tries, maxChildren),
	}
}

func (list *FilterList) length() int {
	return len(list.children)
}

func (list *FilterList) head() *Trie {
	return list.children[0]
}

func (list *FilterList) add(child *Trie) *FilterList {
	// Search for an empty spot and insert the child if possible
	if len(list.children) != cap(list.children) {
		list.children = append(list.children, child)
		return list
	}
	newList := newFilterList(cap(list.children) + 1)
	newList.children = append(list.children, child)
	return newList
}

func (list *FilterList) next(sf string, o int) *Trie {
	for _, child := range list.children {
		if child.filter.stringFilter == sf && child.filter.offset == o {
			return child
		}
	}
	return nil
}

func (list *FilterList) total() int {
	total := 0
	for _, child := range list.children {
		if child != nil {
			total += child.total()
		}
	}
	return total
}

func (list *FilterList) print(w io.Writer, indent int) {
	for _, child := range list.children {
		if child != nil {
			child.print(w, indent)
		}
	}
}
