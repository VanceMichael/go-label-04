package route

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestHubAndSegmentValidation(t *testing.T) {
	p := NewPlanner()
	h := Hub{Code: "PEK", OpenAt: 0, CloseAt: 24 * time.Hour}
	if err := p.AddHub(h); err != nil {
		t.Fatal(err)
	}
	if err := p.AddHub(h); err != domain.ErrConflict {
		t.Fatal(err)
	}
	if err := p.AddSegment(Segment{From: "PEK", To: "PEK", Depart: time.Now(), Arrive: time.Now().Add(time.Hour), AvailableKg: 1}); err != domain.ErrInvalid {
		t.Fatal(err)
	}
}
