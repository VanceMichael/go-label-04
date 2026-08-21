package capacity

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestVersionedHolds(t *testing.T) {
	b := New()
	if err := b.Define("k", 100); err != nil {
		t.Fatal(err)
	}
	v, _ := b.Get("k")
	if err := b.Hold("k", 40, v.Version); err != nil {
		t.Fatal(err)
	}
	if err := b.Hold("k", 40, v.Version); err != domain.ErrConflict {
		t.Fatal(err)
	}
	v, _ = b.Get("k")
	if err := b.Hold("k", 70, v.Version); err != domain.ErrCapacity {
		t.Fatal(err)
	}
	if err := b.Release("k", 40); err != nil {
		t.Fatal(err)
	}
}
