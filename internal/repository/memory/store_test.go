package memory

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestMemoryUsersAndSessions(t *testing.T) {
	s := New()
	now := time.Now()
	u := domain.User{ID: "u", TenantID: "t", Email: "u@example.com", Active: true, Role: domain.RoleCoordinator, CreatedAt: now}
	if err := s.CreateUser(context.Background(), u); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetUserByEmail(context.Background(), u.Email); err != nil {
		t.Fatal(err)
	}
	session := domain.Session{ID: "s", UserID: "u", TokenHash: "hash", ExpiresAt: now.Add(time.Hour), CreatedAt: now}
	if err := s.CreateSession(context.Background(), session); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetSession(context.Background(), "wrong"); err != domain.ErrNotFound {
		t.Fatal(err)
	}
	if err := s.DeactivateUser(context.Background(), "u", now); err != nil {
		t.Fatal(err)
	}
	got, _ := s.GetUser(context.Background(), "u")
	if got.Active {
		t.Fatal("user still active")
	}
}
func TestMemoryShipmentIdempotency(t *testing.T) {
	s := New()
	now := time.Now()
	v := domain.Shipment{ID: "s", TenantID: "t", Reference: "R", Origin: "PEK", Destination: "FRA", WeightKg: 1, Pieces: 1, Status: domain.ShipmentDraft, IdempotencyKey: "k", Version: 1, CreatedAt: now, UpdatedAt: now}
	if err := s.CreateShipment(context.Background(), v); err != nil {
		t.Fatal(err)
	}
	if _, err := s.FindByIdempotency(context.Background(), "t", "k"); err != nil {
		t.Fatal(err)
	}
	v.ID = "s2"
	if err := s.CreateShipment(context.Background(), v); err != domain.ErrConflict {
		t.Fatal(err)
	}
	page, err := s.ListShipments(context.Background(), "t", domain.PageRequest{Limit: 10})
	if err != nil || page.Total != 1 {
		t.Fatal(err, page)
	}
}
func TestMemoryLegConcurrencyVersion(t *testing.T) {
	s := New()
	now := time.Now()
	l := domain.FlightLeg{ID: "l", TenantID: "t", FlightNumber: "F", Origin: "PEK", Destination: "FRA", DepartureAt: now.Add(time.Hour), ArrivalAt: now.Add(2 * time.Hour), CapacityKg: 10, Status: domain.LegPlanned, Version: 1, CreatedAt: now}
	if err := s.CreateLeg(context.Background(), l); err != nil {
		t.Fatal(err)
	}
	if err := s.ReserveCapacity(context.Background(), "t", "l", 5, 1); err != nil {
		t.Fatal(err)
	}
	if err := s.ReserveCapacity(context.Background(), "t", "l", 5, 1); err != domain.ErrConflict {
		t.Fatal(err)
	}
	if err := s.ReserveCapacity(context.Background(), "t", "l", 6, 2); err != domain.ErrCapacity {
		t.Fatal(err)
	}
}

func TestMemoryBookingIsAtomicAndRouteScoped(t *testing.T) {
	s := New()
	now := time.Now().UTC()
	shipment := domain.Shipment{ID: "s", TenantID: "t", Reference: "R", Origin: "PEK", Destination: "FRA", WeightKg: 6, Pieces: 1, Status: domain.ShipmentDraft, Version: 1, CreatedAt: now, UpdatedAt: now}
	if err := s.CreateShipment(context.Background(), shipment); err != nil {
		t.Fatal(err)
	}
	leg := domain.FlightLeg{ID: "l", TenantID: "t", FlightNumber: "AB1", Origin: "PEK", Destination: "FRA", DepartureAt: now.Add(time.Hour), ArrivalAt: now.Add(2 * time.Hour), CapacityKg: 10, Status: domain.LegOpen, Version: 1, CreatedAt: now}
	if err := s.CreateLeg(context.Background(), leg); err != nil {
		t.Fatal(err)
	}
	booked, err := s.BookShipment(context.Background(), "t", "s", "l", 1, 1, now)
	if err != nil || booked.Status != domain.ShipmentBooked {
		t.Fatalf("%v %#v", err, booked)
	}
	storedLeg, _ := s.GetLeg(context.Background(), "t", "l")
	if storedLeg.ReservedKg != 6 {
		t.Fatal(storedLeg.ReservedKg)
	}
}

func TestMemoryBookingRejectsRouteMismatchWithoutReservation(t *testing.T) {
	s := New()
	now := time.Now().UTC()
	_ = s.CreateShipment(context.Background(), domain.Shipment{ID: "s", TenantID: "t", Reference: "R", Origin: "PEK", Destination: "FRA", WeightKg: 6, Pieces: 1, Status: domain.ShipmentDraft, Version: 1, CreatedAt: now, UpdatedAt: now})
	_ = s.CreateLeg(context.Background(), domain.FlightLeg{ID: "l", TenantID: "t", FlightNumber: "AB1", Origin: "PVG", Destination: "FRA", DepartureAt: now.Add(time.Hour), ArrivalAt: now.Add(2 * time.Hour), CapacityKg: 10, Status: domain.LegOpen, Version: 1, CreatedAt: now})
	if _, err := s.BookShipment(context.Background(), "t", "s", "l", 1, 1, now); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	leg, _ := s.GetLeg(context.Background(), "t", "l")
	if leg.ReservedKg != 0 {
		t.Fatal("reservation leaked", leg.ReservedKg)
	}
}
