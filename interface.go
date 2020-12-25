package trie

import "context"

// Tree implemnets ternary trie-tree.
type Tree struct {
	// Root is root of the tree. Only Child is valid.
	root *Node

	// nc means node counts
	nc int
}

// Node implemnets node of ternary trie-tree.
type Node struct {
	label Label
	Value interface{}
	Low   *Node
	High  *Node

	Child *Node
	cc    int // count of children.
}

// Match is matched data.
type Match struct {
	Value interface{}
}

// MatchTree compares a string with multiple strings using Aho-Corasick
// algorithm.
type MatchTree struct {
	root *Node
}

type matchData struct {
	value interface{}
	fail  *Node
}

// Key is a sequence of Labels that comprises an input to a Tree.
type Key interface {
	Iterate(context.Context) <-chan Label
}

type Label interface {
	Compare(Label) int
}
