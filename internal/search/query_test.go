package search

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestFilter(t *testing.T) {
	items := []Item{{ID: "2", Reference: "ABC-2", Status: domain.ShipmentBooked}, {ID: "1", Reference: "ABC-1", Status: domain.ShipmentDraft}}
	out := Filter(items, Query{Term: "ABC", Statuses: []domain.ShipmentStatus{domain.ShipmentBooked}})
	if len(out) != 1 || out[0].ID != "2" {
		t.Fatalf("%#v", out)
	}
}
