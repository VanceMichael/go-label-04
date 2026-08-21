package eta

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestEstimate(t *testing.T) {
	now := time.Now()
	e, err := Calculate("S", []Stop{{Code: "A", Arrive: now.Add(time.Hour), Depart: now.Add(2 * time.Hour)}}, now)
	if err != nil || e.Arrival.IsZero() {
		t.Fatal(err)
	}
	e, err = AddDelay(e, 121)
	if err != nil || e.Confidence != "low" {
		t.Fatalf("%v %#v", err, e)
	}
	if _, err := Calculate("", nil, now); err != domain.ErrInvalid {
		t.Fatal(err)
	}
}
