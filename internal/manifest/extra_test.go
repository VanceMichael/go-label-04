package manifest

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestManifestValidation(t *testing.T) {
	if _, err := Build("M", "S", nil); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	if _, err := Build("M", "S", []Item{{SKU: "", Quantity: 1, WeightKg: 1}}); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	m, err := Build("M", "S", []Item{{SKU: "A", Quantity: 1, WeightKg: 1, Hazardous: true}})
	if err != nil {
		t.Fatal(err)
	}
	if m.HazardousCount() != 1 {
		t.Fatal(m.HazardousCount())
	}
}
