package stream

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
)

type Bus[T any] struct {
	mu          sync.RWMutex
	subscribers map[int]chan T
	next        int
}

func New[T any]() *Bus[T] { return &Bus[T]{subscribers: map[int]chan T{}} }
func (b *Bus[T]) Subscribe(buffer int) (int, <-chan T, error) {
	if buffer < 1 {
		return 0, nil, domain.ErrInvalid
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.next
	b.next++
	ch := make(chan T, buffer)
	b.subscribers[id] = ch
	return id, ch, nil
}
func (b *Bus[T]) Unsubscribe(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if ch, ok := b.subscribers[id]; ok {
		delete(b.subscribers, id)
		close(ch)
	}
}
func (b *Bus[T]) Publish(ctx context.Context, v T) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subscribers {
		select {
		case ch <- v:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
