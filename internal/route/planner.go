package route

import (
	"context"
	"errors"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"time"
)

type Hub struct {
	Code            string
	Timezone        string
	OpenAt          time.Duration
	CloseAt         time.Duration
	HandlingMinutes int
}
type Segment struct {
	From        string
	To          string
	Depart      time.Time
	Arrive      time.Time
	Flight      string
	AvailableKg int64
}
type Plan struct {
	ShipmentID   string
	Segments     []Segment
	TotalMinutes int
	TotalStops   int
}
type Planner struct {
	hubs     map[string]Hub
	segments []Segment
}

func NewPlanner() *Planner { return &Planner{hubs: map[string]Hub{}, segments: []Segment{}} }
func (p *Planner) AddHub(h Hub) error {
	if h.Code == "" || h.HandlingMinutes < 0 || h.CloseAt <= h.OpenAt {
		return domain.ErrInvalid
	}
	if _, ok := p.hubs[h.Code]; ok {
		return domain.ErrConflict
	}
	p.hubs[h.Code] = h
	return nil
}
func (p *Planner) AddSegment(s Segment) error {
	if s.From == "" || s.To == "" || s.From == s.To || !s.Arrive.After(s.Depart) || s.AvailableKg < 0 {
		return domain.ErrInvalid
	}
	if _, ok := p.hubs[s.From]; !ok {
		return domain.ErrNotFound
	}
	if _, ok := p.hubs[s.To]; !ok {
		return domain.ErrNotFound
	}
	p.segments = append(p.segments, s)
	return nil
}
func (p *Planner) Plan(ctx context.Context, shipment string, from, to string, weight int64, earliest time.Time) (Plan, error) {
	if weight <= 0 {
		return Plan{}, domain.ErrInvalid
	}
	queue := []Segment{}
	for _, s := range p.segments {
		if s.From == from && s.AvailableKg >= weight && !s.Depart.Before(earliest) {
			queue = append(queue, s)
		}
	}
	sort.Slice(queue, func(i, j int) bool { return queue[i].Arrive.Before(queue[j].Arrive) })
	if len(queue) == 0 {
		return Plan{}, domain.ErrNotFound
	}
	for _, first := range queue {
		if first.To == to {
			return Plan{ShipmentID: shipment, Segments: []Segment{first}, TotalMinutes: int(first.Arrive.Sub(first.Depart).Minutes())}, nil
		}
		for _, second := range p.segments {
			select {
			case <-ctx.Done():
				return Plan{}, ctx.Err()
			default:
			}
			if second.From == first.To && second.To == to && second.Depart.After(first.Arrive.Add(time.Duration(p.hubs[first.To].HandlingMinutes)*time.Minute)) && second.AvailableKg >= weight {
				total := int(second.Arrive.Sub(first.Depart).Minutes())
				return Plan{ShipmentID: shipment, Segments: []Segment{first, second}, TotalMinutes: total, TotalStops: 1}, nil
			}
		}
	}
	return Plan{}, errors.New("no feasible route")
}
