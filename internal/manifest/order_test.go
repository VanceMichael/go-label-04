package manifest

import (
	"testing"
)

func TestBuildOrderStable(t *testing.T) {
	m, err := Build("M", "S", []Item{{SKU: "B", Quantity: 1, WeightKg: 1}, {SKU: "A", Quantity: 1, WeightKg: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if m.Items[0].SKU != "A" || m.Items[1].SKU != "B" {
		t.Fatal(m.Items)
	}
}
