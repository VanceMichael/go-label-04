package manifest

import (
	"testing"
	"time"
)

func TestManifestSealAndWeight(t *testing.T) {
	m, err := Build("M", "S", []Item{{SKU: "B", Quantity: 2, WeightKg: 3}, {SKU: "A", Quantity: 1, WeightKg: 5}})
	if err != nil {
		t.Fatal(err)
	}
	if m.TotalWeight() != 11 {
		t.Fatal(m.TotalWeight())
	}
	sealed, err := m.Seal(time.Now().UTC())
	if err != nil || sealed.SealedAt == nil {
		t.Fatal(err)
	}
	if _, err = sealed.Seal(time.Now()); err == nil {
		t.Fatal("sealed twice")
	}
}
