package pricing

import (
	"testing"
	"time"
)

func TestQuote(t *testing.T) {
	table := NewTable()
	now := time.Now().UTC()
	_ = table.Add(Rate{Origin: "PEK", Destination: "FRA", MinKg: 1, MaxKg: 100, CentsPerKg: 10, FuelCents: 50, EffectiveFrom: now.Add(-time.Hour), EffectiveTo: now.Add(time.Hour)})
	q, err := table.Quote("PEK", "FRA", 20, now)
	if err != nil || q.TotalCents != 250 {
		t.Fatalf("%v %#v", err, q)
	}
}
