package aggregate

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestAggregateVersion(t *testing.T) {
	s := New()
	if err := s.Create(Shipment{ID: "S", Status: domain.ShipmentDraft, Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.Transition("S", domain.ShipmentBooked, 1); err != nil {
		t.Fatal(err)
	}
	if err := s.Transition("S", domain.ShipmentScreening, 1); err != domain.ErrConflict {
		t.Fatal(err)
	}
	v, _ := s.Get("S")
	if v.Version != 2 || len(v.Events) != 1 {
		t.Fatal(v)
	}
}
