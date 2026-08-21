package booking

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestCapacityLedger(t *testing.T) {
	l := NewLedger()
	if err := l.DefineLeg("L1", 100); err != nil {
		t.Fatal(err)
	}
	a, err := l.Reserve(context.Background(), Request{TenantID: "T", ShipmentID: "S1", LegID: "L1", WeightKg: 60})
	if err != nil || !a.Accepted {
		t.Fatalf("reserve: %v", err)
	}
	if _, err = l.Reserve(context.Background(), Request{TenantID: "T", ShipmentID: "S2", LegID: "L1", WeightKg: 50}); err != domain.ErrCapacity {
		t.Fatalf("want capacity, got %v", err)
	}
	available, _ := l.Available("L1")
	if available != 40 {
		t.Fatalf("available %d", available)
	}
	if err := l.Release(a.ID); err != nil {
		t.Fatal(err)
	}
	available, _ = l.Available("L1")
	if available != 100 {
		t.Fatalf("release did not restore: %d", available)
	}
}
func TestLedgerTenantIsolation(t *testing.T) {
	l := NewLedger()
	_ = l.DefineLeg("L", 10)
	_, _ = l.Reserve(context.Background(), Request{TenantID: "A", ShipmentID: "S", LegID: "L", WeightKg: 1})
	if got := len(l.List("B")); got != 0 {
		t.Fatalf("leaked %d allocations", got)
	}
	if _, err := l.Reserve(context.Background(), Request{LegID: "L", WeightKg: 0}); err != domain.ErrInvalid {
		t.Fatal(err)
	}
}
func TestLedgerCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	l := NewLedger()
	_ = l.DefineLeg("L", 10)
	if _, err := l.Reserve(ctx, Request{TenantID: "T", ShipmentID: "S", LegID: "L", WeightKg: 1}); err != context.Canceled {
		t.Fatal(err)
	}
	_ = time.Now()
}
