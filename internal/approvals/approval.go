package approvals

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
	"time"
)

type Request struct {
	ID          string
	TenantID    string
	ObjectID    string
	RequestedBy string
	ApprovedBy  string
	Status      string
	CreatedAt   time.Time
	DecidedAt   *time.Time
}
type Queue struct {
	mu    sync.Mutex
	items map[string]Request
}

func New() *Queue { return &Queue{items: map[string]Request{}} }
func (q *Queue) Submit(v Request) error {
	if v.ID == "" || v.TenantID == "" || v.ObjectID == "" || v.RequestedBy == "" || v.CreatedAt.IsZero() {
		return domain.ErrInvalid
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.items[v.ID]; ok {
		return domain.ErrConflict
	}
	v.Status = "pending"
	q.items[v.ID] = v
	return nil
}
func (q *Queue) Decide(id, actor string, approve bool, at time.Time) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	v, ok := q.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	if v.Status != "pending" || actor == "" || actor == v.RequestedBy || at.Before(v.CreatedAt) {
		return domain.ErrState
	}
	if approve {
		v.Status = "approved"
	} else {
		v.Status = "rejected"
	}
	v.ApprovedBy = actor
	v.DecidedAt = &at
	q.items[id] = v
	return nil
}
func (q *Queue) Get(id string) (Request, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	v, ok := q.items[id]
	if !ok {
		return Request{}, domain.ErrNotFound
	}
	return v, nil
}
