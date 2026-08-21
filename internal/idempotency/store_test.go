package idempotency

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestTenantScoped(t *testing.T) {
	s := New()
	now := time.Now()
	fp := Fingerprint("POST", "/v1/shipments", []byte("x"))
	if _, replayed, err := s.Begin("A", "k", fp, now); err != nil || replayed {
		t.Fatal(err, replayed)
	}
	if _, replayed, err := s.Begin("A", "k", fp, now); err != nil || !replayed {
		t.Fatal(err, replayed)
	}
	if _, _, err := s.Begin("B", "k", fp, now); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Begin("A", "k", "other", now); err != domain.ErrConflict {
		t.Fatal(err)
	}
}
