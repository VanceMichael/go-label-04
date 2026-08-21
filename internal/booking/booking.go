package booking

import (
	"context"
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"sync"
	"time"
)

type Request struct {
	TenantID    string
	ShipmentID  string
	LegID       string
	WeightKg    int64
	Priority    int
	RequestedAt time.Time
}
type Allocation struct {
	ID        string
	Request   Request
	Accepted  bool
	Reason    string
	CreatedAt time.Time
}
type Ledger struct {
	mu          sync.RWMutex
	capacity    map[string]int64
	reserved    map[string]int64
	allocations map[string]Allocation
}

func NewLedger() *Ledger {
	return &Ledger{capacity: map[string]int64{}, reserved: map[string]int64{}, allocations: map[string]Allocation{}}
}
func (l *Ledger) DefineLeg(id string, capacity int64) error {
	if id == "" || capacity <= 0 {
		return domain.ErrInvalid
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.capacity[id]; ok {
		return domain.ErrConflict
	}
	l.capacity[id] = capacity
	l.reserved[id] = 0
	return nil
}
func (l *Ledger) Reserve(ctx context.Context, r Request) (Allocation, error) {
	select {
	case <-ctx.Done():
		return Allocation{}, ctx.Err()
	default:
	}
	if r.LegID == "" || r.WeightKg <= 0 {
		return Allocation{}, domain.ErrInvalid
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.capacity[r.LegID]; !ok {
		return Allocation{}, domain.ErrNotFound
	}
	if l.reserved[r.LegID]+r.WeightKg > l.capacity[r.LegID] {
		return Allocation{Request: r, Reason: "capacity", CreatedAt: time.Now().UTC()}, domain.ErrCapacity
	}
	id := fmt.Sprintf("%s-%d", r.ShipmentID, time.Now().UnixNano())
	a := Allocation{ID: id, Request: r, Accepted: true, CreatedAt: time.Now().UTC()}
	l.reserved[r.LegID] += r.WeightKg
	l.allocations[id] = a
	return a, nil
}
func (l *Ledger) Release(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	a, ok := l.allocations[id]
	if !ok {
		return domain.ErrNotFound
	}
	if !a.Accepted {
		return nil
	}
	l.reserved[a.Request.LegID] -= a.Request.WeightKg
	a.Accepted = false
	a.Reason = "released"
	l.allocations[id] = a
	return nil
}
func (l *Ledger) Available(id string) (int64, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	cap, ok := l.capacity[id]
	if !ok {
		return 0, domain.ErrNotFound
	}
	return cap - l.reserved[id], nil
}
func (l *Ledger) List(tenant string) []Allocation {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Allocation, 0)
	for _, a := range l.allocations {
		if a.Request.TenantID == tenant {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
