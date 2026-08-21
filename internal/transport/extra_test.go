package transport

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestAssignmentBounds(t *testing.T) {
	p := New()
	now := time.Now()
	_ = p.AddVehicle(Vehicle{ID: "V", Registration: "R", MaxKg: 1, AvailableFrom: now, AvailableTo: now.Add(time.Hour)})
	if err := p.Assign(context.Background(), Assignment{ID: "A", ShipmentID: "S", VehicleID: "V", WeightKg: 2}, now.Add(time.Minute)); err != domain.ErrCapacity {
		t.Fatal(err)
	}
	if err := p.Assign(context.Background(), Assignment{ID: "A", ShipmentID: "S", VehicleID: "X", WeightKg: 1}, now); err != domain.ErrNotFound {
		t.Fatal(err)
	}
}
