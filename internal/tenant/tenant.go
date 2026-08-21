package tenant

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strings"
	"sync"
)

type Registry struct {
	mu    sync.RWMutex
	items map[string]domain.Tenant
}

func New() *Registry { return &Registry{items: map[string]domain.Tenant{}} }
func (r *Registry) Add(t domain.Tenant) error {
	if strings.TrimSpace(t.ID) == "" || strings.TrimSpace(t.Name) == "" {
		return domain.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[t.ID]; ok {
		return domain.ErrConflict
	}
	r.items[t.ID] = t
	return nil
}
func (r *Registry) Deactivate(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	t.Active = false
	r.items[id] = t
	return nil
}
func (r *Registry) Get(id string) (domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	if !ok {
		return domain.Tenant{}, domain.ErrNotFound
	}
	return t, nil
}
func (r *Registry) CanAccess(id string) bool { t, err := r.Get(id); return err == nil && t.Active }
