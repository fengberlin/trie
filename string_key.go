package trie

import (
	"context"
)

type StringKey string

func stringKeyIterate(ctx context.Context, s string, ch chan Label) {
	defer close(ch)
	for _, r := range s {
		select {
		case <-ctx.Done():
			return
		case ch <- r:
		}
	}
}

func (sl StringKey) Iterate(ctx context.Context) <-chan Label {
	ch := make(chan Label)
	go stringKeyIterate(ctx, string(sl), ch)
	return ch
}
