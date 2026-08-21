package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoRetries(t *testing.T) {
	attempts := 0
	err := Do(context.Background(), Policy{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond * 2}, func(context.Context, int) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil || attempts != 3 {
		t.Fatalf("%v %d", err, attempts)
	}
}
func TestDoCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Do(ctx, Policy{MaxAttempts: 4}, func(context.Context, int) error { return errors.New("x") }); err != context.Canceled {
		t.Fatal(err)
	}
}
