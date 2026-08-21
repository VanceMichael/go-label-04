package alerts

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
	"time"
)

type Alert struct {
	ID       string
	TenantID string
	Severity string
	Code     string
	Message  string
	OpenedAt time.Time
	ClosedAt *time.Time
}
type Registry struct {
	mu    sync.Mutex
	items map[string]Alert
}

func New() *Registry { return &Registry{items: map[string]Alert{}} }
func (r *Registry) Open(a Alert) error {
	if a.ID == "" || a.TenantID == "" || a.Code == "" || a.Message == "" || a.OpenedAt.IsZero() {
		return domain.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if old, ok := r.items[a.ID]; ok && old.ClosedAt == nil {
		return domain.ErrConflict
	}
	r.items[a.ID] = a
	return nil
}
func (r *Registry) Close(id string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	if a.ClosedAt != nil || at.Before(a.OpenedAt) {
		return domain.ErrState
	}
	a.ClosedAt = &at
	r.items[id] = a
	return nil
}
func (r *Registry) Active(tenant string) []Alert {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Alert, 0)
	for _, a := range r.items {
		if a.TenantID == tenant && a.ClosedAt == nil {
			out = append(out, a)
		}
	}
	return out
}
