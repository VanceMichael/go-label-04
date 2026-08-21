package warehouse

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestCancelRestoresCapacity(t *testing.T) {
	c := New()
	now := time.Now()
	_ = c.Add(Slot{ID: "S", Terminal: "T", MaxKg: 20, OpenAt: now, CloseAt: now.Add(time.Hour)})
	b := Booking{ID: "B", ShipmentID: "X", SlotID: "S", WeightKg: 10, Start: now.Add(time.Minute), End: now.Add(20 * time.Minute)}
	if err := c.Book(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	u, _ := c.Utilization("S")
	if u != .5 {
		t.Fatal(u)
	}
	if err := c.Cancel("B"); err != nil {
		t.Fatal(err)
	}
	u, _ = c.Utilization("S")
	if u != 0 {
		t.Fatal(u)
	}
	if err := c.Cancel("B"); err != domain.ErrNotFound {
		t.Fatal(err)
	}
}
