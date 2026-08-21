package reconciliation

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestEntryValidation(t *testing.T) {
	if err := Validate(Entry{}); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	if err := Validate(Entry{ExternalID: "x", Kind: "charge", OccurredAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
}
