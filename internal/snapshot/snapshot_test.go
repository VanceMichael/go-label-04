package snapshot

import (
	"testing"
	"time"
)

func TestSnapshotIntegrity(t *testing.T) {
	s, err := Make(1, "T", nil, nil, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	copy := s
	if !s.EqualPayload(copy) {
		t.Fatal("same snapshot differs")
	}
	copy.Version = 2
	if s.EqualPayload(copy) {
		t.Fatal("version ignored")
	}
}
