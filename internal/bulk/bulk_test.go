package bulk

import (
	"context"
	"errors"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestStopOnError(t *testing.T) {
	seen := 0
	p := Processor[int]{Workers: 1, StopOnError: true, Handle: func(context.Context, int) error {
		seen++
		if seen == 1 {
			return errors.New("stop")
		}
		return nil
	}}
	out := p.Run(context.Background(), []int{1, 2, 3})
	if len(out) != 3 {
		t.Fatal(len(out))
	}
	if out[0].Err == nil {
		t.Fatal("first must fail")
	}
	if out[1].Err != context.Canceled && out[2].Err != context.Canceled {
		t.Fatalf("no cancellation: %#v", out)
	}
	_ = domain.ErrInvalid
}
