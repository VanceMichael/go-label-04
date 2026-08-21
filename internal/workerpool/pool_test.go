package workerpool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestPoolProcessesJobs(t *testing.T) {
	p := New(2, 4)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p.Start(ctx)
	var n atomic.Int32
	for i := 0; i < 4; i++ {
		if err := p.Submit(ctx, func(context.Context) error { n.Add(1); return nil }); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(10 * time.Millisecond)
	p.Stop()
	if n.Load() != 4 {
		t.Fatal(n.Load())
	}
}
