package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Trie defines a struct for Trie element
type Trie struct {
	filter   *Filter
	notify   string
	children nodeList
}

// NewTrie constructs Trie
func NewTrie() *Trie {
	trie := &Trie{}
	trie.children = newNodeList(0)
	return trie
}

func (trie *Trie) findNodeWithCommonPrefix(cp string) *Trie {
	// no children
	if trie.children.length() == 0 {
		return trie
	}

	// check if the cp already exists in children
	child := trie.children.next(byte(cp[0]))
	if child != nil {
		if strings.HasPrefix(cp, child.filter.stringFilter) {
			return child.findNodeWithCommonPrefix(cp[len(child.filter.stringFilter):])
		}
	}
	return trie
}

func (trie *Trie) print(writer io.Writer, indent int) {
	fmt.Fprintf(writer, "%s%s %v\n", strings.Repeat(" ", indent), string(trie.filter.stringFilter), trie.filter.offset)
	trie.children.print(writer, indent+2)
}

func (trie *Trie) total() int {
	return 1 + trie.children.total()
}

// NewNode constructs a new branch
func NewNode(f *Filter, n string) *Trie {
	return &Trie{
		filter:   f,
		children: newNodeList(0),
		notify:   n,
	}
}

// Insert inserts a new filter node into the trie using the given filterString and offset
func (trie *Trie) Insert(commonPrefix string, sf string, o int, n string) {
	node := trie

	// Find the right node for the new node's parent
	node = node.findNodeWithCommonPrefix(commonPrefix)

	f := NewFilter(sf, o)
	child := NewNode(f, n)
	node.children.add(child)
}

// Errors ----------------------------------------------------------------------

var (
	// ErrSkipSubtree is to notify skipping this subtree
	ErrSkipSubtree = errors.New("Skip this subtree")
	// ErrNilFilterString notifies nil fs passing
	ErrNilFilterString = errors.New("Nil filterString passed into a method call")
	// ErrNilOffset notifies nil o passing
	ErrNilOffset = errors.New("Nil offset passed into a method call")
)
