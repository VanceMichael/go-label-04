package priority

import (
	"container/heap"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
)

type Item struct {
	ID       string
	Priority int
	Value    any
	index    int
}
type items []*Item

func (i items) Len() int { return len(i) }
func (i items) Less(a, b int) bool {
	if i[a].Priority == i[b].Priority {
		return i[a].ID < i[b].ID
	}
	return i[a].Priority > i[b].Priority
}
func (i items) Swap(a, b int) { i[a], i[b] = i[b], i[a]; i[a].index = a; i[b].index = b }
func (i *items) Push(v any)   { x := v.(*Item); x.index = len(*i); *i = append(*i, x) }
func (i *items) Pop() any {
	old := *i
	n := len(old)
	x := old[n-1]
	*i = old[:n-1]
	x.index = -1
	return x
}

type Queue struct {
	mu    sync.Mutex
	items items
}

func New() *Queue { q := &Queue{items: items{}}; heap.Init(&q.items); return q }
func (q *Queue) Add(v Item) error {
	if v.ID == "" || v.Value == nil {
		return domain.ErrInvalid
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(&q.items, &v)
	return nil
}
func (q *Queue) Take() (Item, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return Item{}, domain.ErrNotFound
	}
	return *heap.Pop(&q.items).(*Item), nil
}
func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
