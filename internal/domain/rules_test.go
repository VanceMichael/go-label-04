package domain

import (
	"testing"
	"time"
)

func TestShipmentTransitions(t *testing.T) {
	cases := []struct {
		from ShipmentStatus
		to   ShipmentStatus
		ok   bool
	}{{ShipmentDraft, ShipmentBooked, true}, {ShipmentBooked, ShipmentScreening, true}, {ShipmentScreening, ShipmentCleared, true}, {ShipmentCleared, ShipmentLoaded, true}, {ShipmentLoaded, ShipmentDeparted, true}, {ShipmentDraft, ShipmentDeparted, false}, {ShipmentDeparted, ShipmentDraft, false}}
	for _, tc := range cases {
		if got := tc.from.CanTransition(tc.to); got != tc.ok {
			t.Fatalf("%s to %s got %v", tc.from, tc.to, got)
		}
	}
}
func TestLegTransitions(t *testing.T) {
	if !LegPlanned.CanTransition(LegOpen) {
		t.Fatal("planned should open")
	}
	if LegOpen.CanTransition(LegDeparted) {
		t.Fatal("open cannot depart")
	}
	if LegClosed.CanTransition(LegOpen) {
		t.Fatal("closed cannot reopen")
	}
}
func TestValidation(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := ValidateShipment(Shipment{TenantID: "t", Reference: "R", Origin: "PEK", Destination: "FRA", WeightKg: 10, Pieces: 1}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateShipment(Shipment{TenantID: "t", Reference: "R", Origin: "PEK", Destination: "PEK", WeightKg: 10, Pieces: 1}); err == nil {
		t.Fatal("same route should fail")
	}
	if err := ValidateLeg(FlightLeg{TenantID: "t", FlightNumber: "AB1", Origin: "PEK", Destination: "FRA", DepartureAt: now.Add(time.Hour), ArrivalAt: now.Add(2 * time.Hour), CapacityKg: 1}, now); err != nil {
		t.Fatal(err)
	}
}
func TestSessionBoundaries(t *testing.T) {
	now := time.Unix(100, 0)
	s := Session{ExpiresAt: now}
	if IsSessionActive(s, now) != ErrExpired {
		t.Fatal("expiry must be inclusive")
	}
	s.ExpiresAt = now.Add(time.Second)
	if IsSessionActive(s, now) != nil {
		t.Fatal("session should be active")
	}
}
