package locks

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
	"time"
)

type Keyed struct {
	mu   sync.Mutex
	held map[string]chan struct{}
}

func New() *Keyed { return &Keyed{held: map[string]chan struct{}{}} }
func (k *Keyed) Acquire(ctx context.Context, key string) (func(), error) {
	if key == "" {
		return nil, domain.ErrInvalid
	}
	k.mu.Lock()
	ch, ok := k.held[key]
	if !ok {
		ch = make(chan struct{}, 1)
		k.held[key] = ch
	}
	k.mu.Unlock()
	select {
	case ch <- struct{}{}:
		return func() { <-ch }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
func (k *Keyed) WithTimeout(ctx context.Context, key string, d time.Duration) (func(), error) {
	if d <= 0 {
		return nil, domain.ErrInvalid
	}
	child, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	return k.Acquire(child, key)
}
