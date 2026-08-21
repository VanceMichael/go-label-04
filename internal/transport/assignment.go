package transport

import (
	"context"
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"sync"
	"time"
)

type Vehicle struct {
	ID            string
	Registration  string
	Carrier       string
	MaxKg         int64
	AvailableFrom time.Time
	AvailableTo   time.Time
}
type Assignment struct {
	ID         string
	ShipmentID string
	VehicleID  string
	WeightKg   int64
	Status     string
	AssignedAt time.Time
}
type Planner struct {
	mu          sync.RWMutex
	vehicles    map[string]Vehicle
	assignments map[string]Assignment
}

func New() *Planner {
	return &Planner{vehicles: map[string]Vehicle{}, assignments: map[string]Assignment{}}
}
func (p *Planner) AddVehicle(v Vehicle) error {
	if v.ID == "" || v.Registration == "" || v.MaxKg <= 0 || !v.AvailableTo.After(v.AvailableFrom) {
		return domain.ErrInvalid
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.vehicles[v.ID]; ok {
		return domain.ErrConflict
	}
	p.vehicles[v.ID] = v
	return nil
}
func (p *Planner) Assign(ctx context.Context, a Assignment, at time.Time) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if a.ID == "" || a.ShipmentID == "" || a.VehicleID == "" || a.WeightKg <= 0 {
		return domain.ErrInvalid
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	v, ok := p.vehicles[a.VehicleID]
	if !ok {
		return domain.ErrNotFound
	}
	if a.WeightKg > v.MaxKg || at.Before(v.AvailableFrom) || !at.Before(v.AvailableTo) {
		return domain.ErrCapacity
	}
	for _, old := range p.assignments {
		if old.VehicleID == a.VehicleID && old.Status != "cancelled" {
			return domain.ErrConflict
		}
	}
	a.Status = "assigned"
	a.AssignedAt = at
	p.assignments[a.ID] = a
	return nil
}
func (p *Planner) Complete(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	a, ok := p.assignments[id]
	if !ok {
		return domain.ErrNotFound
	}
	if a.Status != "assigned" {
		return domain.ErrState
	}
	a.Status = "completed"
	p.assignments[id] = a
	return nil
}
func (p *Planner) List(vehicle string) []Assignment {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Assignment, 0)
	for _, a := range p.assignments {
		if vehicle == "" || a.VehicleID == vehicle {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].AssignedAt.Before(out[j].AssignedAt) })
	return out
}
func (p *Planner) Describe(id string) (Assignment, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	a, ok := p.assignments[id]
	if !ok {
		return Assignment{}, fmt.Errorf("%w: assignment", domain.ErrNotFound)
	}
	return a, nil
}
