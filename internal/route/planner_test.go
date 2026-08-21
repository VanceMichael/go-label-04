package route

import (
	"context"
	"testing"
	"time"
)

func TestTwoHopPlan(t *testing.T) {
	p := NewPlanner()
	_ = p.AddHub(Hub{Code: "PEK", OpenAt: 0, CloseAt: 24 * time.Hour})
	_ = p.AddHub(Hub{Code: "DXB", OpenAt: 0, CloseAt: 24 * time.Hour, HandlingMinutes: 30})
	_ = p.AddHub(Hub{Code: "FRA", OpenAt: 0, CloseAt: 24 * time.Hour})
	now := time.Now().UTC()
	_ = p.AddSegment(Segment{From: "PEK", To: "DXB", Depart: now.Add(time.Hour), Arrive: now.Add(5 * time.Hour), AvailableKg: 100})
	_ = p.AddSegment(Segment{From: "DXB", To: "FRA", Depart: now.Add(6 * time.Hour), Arrive: now.Add(12 * time.Hour), AvailableKg: 100})
	plan, err := p.Plan(context.Background(), "S", "PEK", "FRA", 10, now)
	if err != nil || len(plan.Segments) != 2 {
		t.Fatalf("plan: %v %#v", err, plan)
	}
}
