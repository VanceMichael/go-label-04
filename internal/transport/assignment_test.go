package transport

import (
	"context"
	"testing"
	"time"
)

func TestAssignment(t *testing.T) {
	p := New()
	now := time.Now()
	if err := p.AddVehicle(Vehicle{ID: "V", Registration: "R", MaxKg: 100, AvailableFrom: now.Add(-time.Hour), AvailableTo: now.Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if err := p.Assign(context.Background(), Assignment{ID: "A", ShipmentID: "S", VehicleID: "V", WeightKg: 10}, now); err != nil {
		t.Fatal(err)
	}
	if err := p.Complete("A"); err != nil {
		t.Fatal(err)
	}
	if len(p.List("V")) != 1 {
		t.Fatal("missing assignment")
	}
}
