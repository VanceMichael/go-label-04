package warehouse

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestSlotConflicts(t *testing.T) {
	c := New()
	now := time.Now()
	_ = c.Add(Slot{ID: "S", Terminal: "T", MaxKg: 100, OpenAt: now, CloseAt: now.Add(time.Hour)})
	b := Booking{ID: "B", ShipmentID: "X", SlotID: "S", WeightKg: 60, Start: now.Add(time.Minute), End: now.Add(20 * time.Minute)}
	if err := c.Book(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	b.ID = "C"
	b.WeightKg = 10
	if err := c.Book(context.Background(), b); err != domain.ErrConflict {
		t.Fatal(err)
	}
}
