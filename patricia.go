package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Errors ----------------------------------------------------------------------

var (
	SkipSubtree  = errors.New("Skip this subtree")
	ErrNilFilterString = errors.New("Nil filterString passed into a method call")
	ErrNilOffset = errors.New("Nil offset passed into a method call")
)

type Trie struct {
	filter *Filter
	children *FilterList
}

func (trie *Trie) print(writer io.Writer, indent int) {
	fmt.Fprintf(writer, "%s%s %v\n", strings.Repeat(" ", indent), string(trie.filter.stringFilter), trie.filter.offset)
	trie.children.print(writer, indent+2)
}

func (trie *Trie) total() int {
	return 1 + trie.children.total()
}

type tries []*Trie

func (t tries) Len() int {
	return len(t)
}

// NewTrie constructs Trie
func NewTrie(f *Filter) *Trie {
	return &Trie{
		filter: f,
		children: newFilterList(0),
	}
}

/*
// Insert inserts a new filter node into the trie using the given filterString and offset
func (trie *Trie) Insert(fs string, o int) (inserted bool) {
	var (
		common int
		node   *Trie = trie
		child  *Trie
	)

	if node.filter == nil {
		node.filter = newFilter(fs, o)
		//// from here ////
		key = key[trie.maxPrefixPerNode:]
		goto AppendChild
	}

	for {
		// Compute the longest common prefix length.
		common = node.longestCommonPrefixLength(key)
		key = key[common:]

		// Only a part matches, split.
		if common < len(node.prefix) {
			goto SplitPrefix
		}

		// common == len(node.prefix) since never (common > len(node.prefix))
		// common == len(former key) <-> 0 == len(key)
		// -> former key == node.prefix
		if len(key) == 0 {
			goto InsertItem
		}

		// Check children for matching prefix.
		child = node.children.next(key[0])
		if child == nil {
			goto AppendChild
		}
		node = child
	}

SplitPrefix:
	// Split the prefix if necessary.
	child = new(Trie)
	*child = *node
	*node = *NewTrie()
	node.prefix = child.prefix[:common]
	child.prefix = child.prefix[common:]
	child = child.compact()
	node.children = node.children.add(child)

AppendChild:
	// Keep appending children until whole prefix is inserted.
	// This loop starts with empty node.prefix that needs to be filled.
	for len(key) != 0 {
		child := NewTrie()
		if len(key) <= trie.maxPrefixPerNode {
			child.prefix = key
			node.children = node.children.add(child)
			node = child
			goto InsertItem
		} else {
			child.prefix = key[:trie.maxPrefixPerNode]
			key = key[trie.maxPrefixPerNode:]
			node.children = node.children.add(child)
			node = child
		}
	}

InsertItem:
	// Try to insert the item if possible.
	if replace || node.item == nil {
		node.item = item
		return true
	}
	return false
}
*/
