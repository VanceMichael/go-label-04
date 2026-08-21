package capacity

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestCapacityLifecycle(t *testing.T) {
	b := New()
	if err := b.Define("", 1); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	if err := b.Define("x", 0); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	if _, err := b.Get("none"); err != domain.ErrNotFound {
		t.Fatal(err)
	}
	if err := b.Release("none", 1); err != domain.ErrNotFound {
		t.Fatal(err)
	}
}
