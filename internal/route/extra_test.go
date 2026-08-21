package route

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestPlannerRejectsInvalid(t *testing.T) {
	p := NewPlanner()
	if err := p.AddHub(Hub{Code: "PEK", OpenAt: time.Hour, CloseAt: time.Hour}); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	if err := p.AddSegment(Segment{From: "PEK", To: "FRA", Depart: time.Now(), Arrive: time.Now().Add(time.Hour), AvailableKg: 1}); err != domain.ErrNotFound {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := p.Plan(ctx, "S", "A", "B", 1, time.Now())
	if err != domain.ErrNotFound && err != domain.ErrInvalid {
		t.Fatal(err)
	}
}
