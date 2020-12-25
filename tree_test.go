package trie_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lestrrat-go/trie"
	"github.com/stretchr/testify/assert"
)

func TestEach(t *testing.T) {
	tr := trie.New()
	tr.Put(trie.StringKey("foo"), "123")
	tr.Put(trie.StringKey("bar"), "999")
	tr.Put(trie.StringKey("日本語"), "こんにちは")

	expected := []rune{
		0,
		'b',
		'f',
		'日',
		'a',
		'o',
		'本',
		'r',
		'o',
		'語',
	}

	i := 0
	for n := range tr.Iterate(context.Background()) {
		t.Logf("%c", n.Label())

		var r rune
		if l := n.Label(); l != nil {
			r = l.(trie.RuneLabel).Rune()
		}
		if !assert.Equal(t, expected[i], r, `labels should match for input %d`, i) {
			return
		}
		i++
	}
}

func TestPut(t *testing.T) {
	f := func(t *testing.T, tr *trie.Tree, key trie.Key, value interface{}) {
		t.Helper()

		n := tr.Get(key)
		if value == nil {
			assert.Equal(t, n, (*trie.Node)(nil), "no nodes for %q", key)
			return
		}
		assert.Equal(t, n.Value, value, "value for %q", key)
	}

	testcases := []struct {
		Key   trie.Key
		Value interface{}
	}{
		{Key: trie.StringKey("foo"), Value: "123"},
		{Key: trie.StringKey("bar"), Value: "999"},
		{Key: trie.StringKey("日本語"), Value: "こんにちは"},
		{Key: trie.StringKey("baz")},
		{Key: trie.StringKey("English")},
	}

	tr := trie.New()
	tr.Put(trie.StringKey("foo"), "123")
	tr.Put(trie.StringKey("bar"), "999")
	tr.Put(trie.StringKey("日本語"), "こんにちは")

	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%s", tc.Key), func(t *testing.T) {
			f(t, tr, tc.Key, tc.Value)
		})
	}
}

func TestTree_nc(t *testing.T) {
	tr := trie.New()
	tr.Put(trie.StringKey("foo"), "123")
	tr.Put(trie.StringKey("bar"), "999")
	tr.Put(trie.StringKey("日本語"), "こんにちは")
	if cc := tr.NodeCount(); cc != 9 {
		t.Errorf("nc mismatch: %d", cc)
	}
}

func TestNode_cc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := func(key trie.Key, cc int) {
		n := new(trie.Node)
		for l := range key.Iterate(ctx) {
			n.Dig(l)
		}
		if !assert.Equal(t, n.ChildCount(), cc, "runes: %q", key) {
			return
		}
	}
	f(trie.StringKey(""), 0)
	f(trie.StringKey("a"), 1)
	f(trie.StringKey("bac"), 3)
	f(trie.StringKey("aaa"), 1)
	f(trie.StringKey("bbbaaaccc"), 3)
	f(trie.StringKey("bacbacbac"), 3)
	f(trie.StringKey("日本語こんにちは"), 8)
	f(trie.StringKey("あめんぼあかいなあいうえお"), 10)
}

// collectRunes1 coolects label runes from sibling nodes.
func collectRunes1(n *trie.Node, max int) []rune {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var runes []rune
	for q := range n.Iterate(trie.WithBFS(ctx)) {
		runes = append(runes, q.Label().(trie.RuneLabel).Rune())
	}
	return runes
}

// collectRunes2 coolects label runes from sibling nodes in reverse order.
func collectRunes2(n *trie.Node, max int) []rune {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var runes []rune
	for q := range n.Iterate(trie.WithBFSReverse(ctx)) {
		runes = append(runes, q.Label().(trie.RuneLabel).Rune())
	}
	return runes
}

func TestNode_Balance(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := new(trie.Node)
	for l := range trie.StringKey("123456789ABCDEF").Iterate(ctx) {
		n.Dig(l)
	}
	n.Balance()

	if !assert.NotNil(t, n.Child, "Child shoud not be nil after balancing") {
		return
	}

	r1 := collectRunes1(n.Child, n.ChildCount())
	if !assert.Equal(t, "84C26AE13579BDF", string(r1), "should be balanced") {
		return
	}
	r2 := collectRunes2(n.Child, n.ChildCount())
	if !assert.Equal(t, "8C4EA62FDB97531", string(r2), "should be balanced") {
		return
	}
}
