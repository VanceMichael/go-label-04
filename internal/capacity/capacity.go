package capacity

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
)

type Bucket struct {
	Limit   int64
	Used    int64
	Version int64
}
type Book struct {
	mu      sync.Mutex
	buckets map[string]Bucket
}

func New() *Book { return &Book{buckets: make(map[string]Bucket)} }
func (b *Book) Define(key string, limit int64) error {
	if key == "" || limit <= 0 {
		return domain.ErrInvalid
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.buckets[key]; ok {
		return domain.ErrConflict
	}
	b.buckets[key] = Bucket{Limit: limit, Version: 1}
	return nil
}
func (b *Book) Hold(key string, amount int64, version int64) error {
	if amount <= 0 {
		return domain.ErrInvalid
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	v, ok := b.buckets[key]
	if !ok {
		return domain.ErrNotFound
	}
	if v.Version != version {
		return domain.ErrConflict
	}
	if v.Used+amount > v.Limit {
		return domain.ErrCapacity
	}
	v.Used += amount
	v.Version++
	b.buckets[key] = v
	return nil
}
func (b *Book) Release(key string, amount int64) error {
	if amount <= 0 {
		return domain.ErrInvalid
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	v, ok := b.buckets[key]
	if !ok {
		return domain.ErrNotFound
	}
	if amount > v.Used {
		return domain.ErrInvalid
	}
	v.Used -= amount
	v.Version++
	b.buckets[key] = v
	return nil
}
func (b *Book) Get(key string) (Bucket, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	v, ok := b.buckets[key]
	if !ok {
		return Bucket{}, domain.ErrNotFound
	}
	return v, nil
}
