package trie_test

import (
	"testing"

	"github.com/lestrrat-go/trie"
	"github.com/stretchr/testify/assert"
)

func toInts(matches []trie.Match) []int {
	if len(matches) == 0 {
		return nil
	}
	r := make([]int, len(matches))
	for i, m := range matches {
		r[i] = m.Value.(int)
	}
	return r
}

func TestMatch(t *testing.T) {
	// Build tree.
	tr := trie.New()
	tr.Put(trie.StringKey("ab"), 2)
	tr.Put(trie.StringKey("bc"), 4)
	tr.Put(trie.StringKey("bab"), 6)
	tr.Put(trie.StringKey("d"), 7)
	tr.Put(trie.StringKey("abcde"), 10)
	mt := trie.Compile(tr)

	// Check tree.
	f := func(key trie.Key, exp []int) {
		act := toInts(mt.MatchAll(key, nil))
		assert.Equal(t, act, exp, "not match for key=%q", key)
	}
	f(trie.StringKey("ab"), []int{2})
	f(trie.StringKey("bc"), []int{4})
	f(trie.StringKey("d"), []int{7})
	f(trie.StringKey("abcde"), []int{2, 4, 7, 10})
	f(trie.StringKey("babc"), []int{6, 2, 4})
}
