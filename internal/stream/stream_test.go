package stream

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestBus(t *testing.T) {
	b := New[int]()
	id, ch, err := b.Subscribe(1)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	if got := <-ch; got != 7 {
		t.Fatal(got)
	}
	b.Unsubscribe(id)
	if err := b.Publish(context.Background(), 8); err != nil {
		t.Fatal(err)
	}
	_, _, err = b.Subscribe(0)
	if err != domain.ErrInvalid {
		t.Fatal(err)
	}
}
