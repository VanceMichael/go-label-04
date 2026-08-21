package domain

import "time"

func (s ShipmentStatus) CanTransition(next ShipmentStatus) bool {
	allowed := map[ShipmentStatus]map[ShipmentStatus]bool{
		ShipmentDraft:     {ShipmentBooked: true, ShipmentCancelled: true},
		ShipmentBooked:    {ShipmentScreening: true, ShipmentCancelled: true},
		ShipmentScreening: {ShipmentCleared: true, ShipmentCancelled: true},
		ShipmentCleared:   {ShipmentLoaded: true, ShipmentCancelled: true},
		ShipmentLoaded:    {ShipmentDeparted: true}, ShipmentDeparted: {}, ShipmentCancelled: {},
	}
	return allowed[s][next]
}

func (s Shipment) CanTransition(next ShipmentStatus) bool { return s.Status.CanTransition(next) }

func (s LegStatus) CanTransition(next LegStatus) bool {
	return map[LegStatus]map[LegStatus]bool{LegPlanned: {LegOpen: true}, LegOpen: {LegBoarding: true, LegClosed: true}, LegBoarding: {LegDeparted: true, LegClosed: true}, LegDeparted: {LegClosed: true}, LegClosed: {}}[s][next]
}
func (s CustomsStatus) CanTransition(next CustomsStatus) bool {
	return map[CustomsStatus]map[CustomsStatus]bool{CustomsPending: {CustomsReview: true, CustomsHeld: true}, CustomsReview: {CustomsReleased: true, CustomsHeld: true}, CustomsHeld: {CustomsReview: true}, CustomsReleased: {}}[s][next]
}
func (s SecurityStatus) CanTransition(next SecurityStatus) bool {
	return map[SecurityStatus]map[SecurityStatus]bool{SecurityPending: {SecurityPassed: true, SecurityFailed: true}, SecurityFailed: {SecurityPending: true}, SecurityPassed: {}}[s][next]
}

func ValidateShipment(s Shipment) error {
	if s.TenantID == "" || s.Reference == "" || s.Origin == "" || s.Destination == "" || s.WeightKg <= 0 || s.Pieces <= 0 {
		return ErrInvalid
	}
	if s.Origin == s.Destination {
		return ErrInvalid
	}
	return nil
}
func ValidateLeg(l FlightLeg, now time.Time) error {
	if l.TenantID == "" || l.FlightNumber == "" || l.Origin == "" || l.Destination == "" || l.CapacityKg <= 0 || !l.DepartureAt.After(now) || !l.ArrivalAt.After(l.DepartureAt) {
		return ErrInvalid
	}
	return nil
}
func ValidateCapacity(l FlightLeg, weight int64) error {
	if weight <= 0 || l.ReservedKg < 0 || l.ReservedKg+weight > l.CapacityKg {
		return ErrCapacity
	}
	return nil
}
func CustomsAllowsLoading(c CustomsCase) bool    { return c.Status == CustomsReleased }
func SecurityAllowsLoading(s SecurityCheck) bool { return s.Status == SecurityPassed }
func IsSessionActive(s Session, now time.Time) error {
	if s.RevokedAt != nil {
		return ErrRevoked
	}
	if !now.Before(s.ExpiresAt) {
		return ErrExpired
	}
	return nil
}
