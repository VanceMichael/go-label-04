package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

type Policy struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
}

func (p Policy) Normalize() Policy {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.InitialDelay <= 0 {
		p.InitialDelay = 100 * time.Millisecond
	}
	if p.MaxDelay <= 0 {
		p.MaxDelay = 10 * time.Second
	}
	return p
}
func (p Policy) Delay(attempt int) time.Duration {
	p = p.Normalize()
	if attempt < 1 {
		attempt = 1
	}
	factor := math.Pow(2, float64(attempt-1))
	d := time.Duration(float64(p.InitialDelay) * factor)
	if d > p.MaxDelay {
		return p.MaxDelay
	}
	return d
}
func Do(ctx context.Context, p Policy, fn func(context.Context, int) error) error {
	p = p.Normalize()
	var last error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := fn(ctx, attempt)
		if err == nil {
			return nil
		}
		last = err
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		if attempt == p.MaxAttempts {
			break
		}
		timer := time.NewTimer(p.Delay(attempt))
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return last
}
