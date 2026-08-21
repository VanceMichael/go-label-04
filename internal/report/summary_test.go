package report

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestBuildSummary(t *testing.T) {
	now := time.Now()
	s := Build([]ShipmentRow{{Reference: "A", Status: domain.ShipmentBooked, WeightKg: 10, CreatedAt: now}, {Reference: "B", Status: domain.ShipmentDraft, WeightKg: 5, CreatedAt: now.Add(time.Minute)}}, 1)
	if s.Total != 2 || s.WeightKg != 15 || len(s.Latest) != 1 {
		t.Fatalf("%#v", s)
	}
	if s.ByStatus[domain.ShipmentDraft] != 1 {
		t.Fatal(s.ByStatus)
	}
}
