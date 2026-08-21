package scan

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestSignals(t *testing.T) {
	s := New()
	r, err := s.Check(context.Background(), "S", "O", map[string]bool{"seal": false})
	if err != nil || r.Status != domain.SecurityPassed {
		t.Fatal(err)
	}
	r, err = s.Check(context.Background(), "S", "O", map[string]bool{"xray": true})
	if err != nil || r.Status != domain.SecurityFailed || !s.Blocked("S") {
		t.Fatalf("%v %#v", err, r)
	}
}
