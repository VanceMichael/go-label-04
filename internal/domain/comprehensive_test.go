package domain

import (
	"testing"
	"time"
)

func TestComprehensiveDomainCases(t *testing.T) {
	cases := []struct {
		name string
		from ShipmentStatus
		to   ShipmentStatus
		ok   bool
	}{
		{"draft book", ShipmentDraft, ShipmentBooked, true},
		{"draft cancel", ShipmentDraft, ShipmentCancelled, true},
		{"book screen", ShipmentBooked, ShipmentScreening, true},
		{"book cancel", ShipmentBooked, ShipmentCancelled, true},
		{"screen clear", ShipmentScreening, ShipmentCleared, true},
		{"screen hold", ShipmentScreening, ShipmentCancelled, true},
		{"clear load", ShipmentCleared, ShipmentLoaded, true},
		{"load depart", ShipmentLoaded, ShipmentDeparted, true},
		{"depart draft", ShipmentDeparted, ShipmentDraft, false},
		{"cancel book", ShipmentCancelled, ShipmentBooked, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.from.CanTransition(tc.to); got != tc.ok {
				t.Fatalf("got %v", got)
			}
		})
	}
}
func TestRulesRejectBadCapacity(t *testing.T) {
	leg := FlightLeg{CapacityKg: 10, ReservedKg: 9}
	if err := ValidateCapacity(leg, 1); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCapacity(leg, 2); err != ErrCapacity {
		t.Fatal(err)
	}
	leg.ReservedKg = -1
	if err := ValidateCapacity(leg, 1); err != ErrCapacity {
		t.Fatal(err)
	}
}
func TestRulesStatusFamilies(t *testing.T) {
	if !CustomsPending.CanTransition(CustomsReview) {
		t.Fatal("customs pending")
	}
	if CustomsReleased.CanTransition(CustomsReview) {
		t.Fatal("customs released")
	}
	if !SecurityPending.CanTransition(SecurityPassed) {
		t.Fatal("security pending")
	}
	if SecurityPassed.CanTransition(SecurityFailed) {
		t.Fatal("security passed")
	}
}
func TestClockAndPagination(t *testing.T) {
	now := time.Now()
	clock := FixedClock{T: now}
	if !clock.Now().Equal(now) {
		t.Fatal("clock changed")
	}
	req := PageRequest{Limit: 1000}.Normalized()
	if req.Limit != 200 {
		t.Fatal(req.Limit)
	}
	req = PageRequest{Limit: 0}.Normalized()
	if req.Limit != 50 {
		t.Fatal(req.Limit)
	}
}
func TestValidationMessages(t *testing.T) {
	bad := []Shipment{{}, {TenantID: "t"}, {TenantID: "t", Reference: "r", Origin: "A", Destination: "B", WeightKg: -1, Pieces: 1}, {TenantID: "t", Reference: "r", Origin: "A", Destination: "B", WeightKg: 1, Pieces: 0}}
	for i, v := range bad {
		if err := ValidateShipment(v); err == nil {
			t.Fatalf("case %d accepted", i)
		}
	}
}
