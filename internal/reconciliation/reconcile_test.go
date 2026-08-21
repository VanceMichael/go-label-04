package reconciliation

import (
	"testing"
	"time"
)

func TestCompare(t *testing.T) {
	now := time.Now()
	left := []Entry{{ExternalID: "1", Kind: "charge", Amount: 10, OccurredAt: now}, {ExternalID: "2", Kind: "charge", Amount: 5, OccurredAt: now}}
	right := []Entry{{ExternalID: "1", Kind: "charge", Amount: 11, OccurredAt: now}, {ExternalID: "3", Kind: "charge", Amount: 5, OccurredAt: now}}
	d := Compare(left, right)
	if len(d) != 3 {
		t.Fatalf("%#v", d)
	}
}
