package idempotency

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestIdempotencyExpiry(t *testing.T) {
	s := New()
	at := time.Now()
	fp := Fingerprint("POST", "/x", nil)
	_, _, _ = s.Begin("T", "k", fp, at)
	_ = s.Complete("T", "k", 201, []byte("ok"))
	r, replay, err := s.Begin("T", "k", fp, at.Add(time.Hour))
	if err != nil || !replay || r.Status != 201 {
		t.Fatalf("%v %v %#v", err, replay, r)
	}
	if err := s.Complete("T", "missing", 200, nil); err != domain.ErrNotFound {
		t.Fatal(err)
	}
	if n := s.Cleanup(at.Add(48 * time.Hour)); n != 1 {
		t.Fatal(n)
	}
}
