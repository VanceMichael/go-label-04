package snapshot

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestSnapshotMutationDetected(t *testing.T) {
	s, err := Make(1, "T", []domain.Shipment{{ID: "S"}}, nil, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	s.Shipments[0].Reference = "changed"
	if err := s.Validate(); err != domain.ErrConflict {
		t.Fatal(err)
	}
}
