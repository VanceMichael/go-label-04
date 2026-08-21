package scan

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestScanInput(t *testing.T) {
	s := New()
	if _, err := s.Check(context.Background(), "", "O", nil); err != domain.ErrInvalid {
		t.Fatal(err)
	}
	if _, err := s.Get("none"); err != domain.ErrNotFound {
		t.Fatal(err)
	}
	r, err := s.Check(context.Background(), "S", "O", map[string]bool{})
	if err != nil || r.Status != domain.SecurityPassed {
		t.Fatal(err)
	}
}
