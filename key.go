package trie

import "context"

type StringKey string
type RuneLabel rune

func stringKeyIterate(ctx context.Context, s string, ch chan Label) {
	defer close(ch)
	for _, r := range s {
		select {
		case <-ctx.Done():
			return
		case ch <- RuneLabel(r):
		}
	}
}

func (sl StringKey) Iterate(ctx context.Context) <-chan Label {
	ch := make(chan Label)
	go stringKeyIterate(ctx, string(sl), ch)
	return ch
}

func (l1 RuneLabel) Compare(l Label) int {
	l2, ok := l.(RuneLabel)
	if !ok { // meh...
		return -1
	}

	switch {
	case l1 < l2:
		return -1
	case l1 > l2:
		return 1
	default:
		return 0
	}
}

func (l RuneLabel) Rune() rune {
	return rune(l)
}

func (l RuneLabel) String() string {
	return string([]rune{l.Rune()})
}
