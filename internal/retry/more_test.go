package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
	defer cancel()
	err := Do(ctx, Policy{MaxAttempts: 100, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond}, func(context.Context, int) error { return errors.New("x") })
	if err == nil {
		t.Fatal("expected deadline")
	}
}
