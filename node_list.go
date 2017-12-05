package main

import (
  "io"
  "sort"
)

const (
  // DefaultMaxChildrenPerNode defines a default number of children in a node
  DefaultMaxChildrenPerNode = 2
)

// nodeList defines a set of tries
type nodeList interface {
  add(child *Trie) nodeList
  head() *Trie
  length() int
  next(b byte) *Trie
  print(w io.Writer, indent int)
  total() int
}

type tries []*Trie

func (t tries) Len() int {
  return len(t)
}

func (t tries) Less(i, j int) bool {
  strings := sort.StringSlice{string(t[i].filter.stringFilter), string(t[j].filter.stringFilter)}
  return strings.Less(0, 1)
}

func (t tries) Swap(i, j int) {
  t[i], t[j] = t[j], t[i]
}

type nodeListStruct struct {
  children tries
}

// newNodeList creates a new set of nodeList
func newNodeList(maxChildren int) nodeList {
  if maxChildren == 0 {
    maxChildren = DefaultMaxChildrenPerNode
  }
  return &nodeListStruct{
    children: make(tries, 0, maxChildren),
  }
}

func (list *nodeListStruct) add(child *Trie) nodeList {
  // Search for an empty spot and insert the child if possible
  if len(list.children) != cap(list.children) {
    list.children = append(list.children, child)
    return list
  }
  //////////////////////////////////
  // the program cannot reach here
  //////////////////////////////////
  newList := newNodeList(cap(list.children) + 1)
  for _, n := range list.children {
    newList.add(n)
  }
  newList.add(child)
  return newList
}

func (list *nodeListStruct) head() *Trie {
  return list.children[0]
}

func (list *nodeListStruct) length() int {
  return len(list.children)
}

func (list *nodeListStruct) next(b byte) *Trie {
	for _, child := range list.children {
		if child.filter.stringFilter[0] == b {
			return child
		}
	}
	return nil
}

func (list *nodeListStruct) print(w io.Writer, indent int) {
  for _, child := range list.children {
    if child != nil {
      child.print(w, indent)
    }
  }
}

func (list *nodeListStruct) total() int {
  total := 0
  for _, child := range list.children {
    if child != nil {
      total += child.total()
    }
  }
  return total
}
