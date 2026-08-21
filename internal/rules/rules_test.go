package rules

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestRuleSet(t *testing.T) {
	s := New()
	if err := s.Add(Before(time.Now().Add(time.Hour))); err != nil {
		t.Fatal(err)
	}
	if err := s.Evaluate(map[string]any{"at": time.Now()}); err != nil {
		t.Fatal(err)
	}
	if err := s.Evaluate(map[string]any{"at": time.Now().Add(2 * time.Hour)}); err != domain.ErrExpired {
		t.Fatal(err)
	}
}
