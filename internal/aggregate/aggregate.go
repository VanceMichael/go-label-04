package aggregate

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
)

type Shipment struct {
	ID      string
	Status  domain.ShipmentStatus
	Version int64
	Events  []string
}
type Store struct {
	mu    sync.RWMutex
	items map[string]Shipment
}

func New() *Store { return &Store{items: map[string]Shipment{}} }
func (s *Store) Create(v Shipment) error {
	if v.ID == "" || v.Version < 1 {
		return domain.ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[v.ID]; ok {
		return domain.ErrConflict
	}
	v.Events = append([]string(nil), v.Events...)
	s.items[v.ID] = v
	return nil
}
func (s *Store) Transition(id string, next domain.ShipmentStatus, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	if v.Version != version {
		return domain.ErrConflict
	}
	if !v.Status.CanTransition(next) {
		return domain.ErrState
	}
	v.Status = next
	v.Version++
	v.Events = append(v.Events, string(next))
	s.items[id] = v
	return nil
}
func (s *Store) Get(id string) (Shipment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.items[id]
	if !ok {
		return Shipment{}, domain.ErrNotFound
	}
	v.Events = append([]string(nil), v.Events...)
	return v, nil
}
