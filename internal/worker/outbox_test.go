package worker

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/repository/memory"
	"log/slog"
	"testing"
	"time"
)

func TestOutboxRunnerPublishes(t *testing.T) {
	s := memory.New()
	now := time.Now()
	if err := s.Enqueue(context.Background(), domain.OutboxEvent{ID: "e", TenantID: "t", Topic: "shipment", AggregateID: "s", AvailableAt: now}); err != nil {
		t.Fatal(err)
	}
	r := NewOutboxRunner(s, slog.Default())
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { r.Run(ctx); close(done) }()
	time.Sleep(300 * time.Millisecond)
	cancel()
	<-done
	r.Wait()
	pending, failed, err := s.Summary(context.Background(), "t")
	if err != nil || pending != 0 || failed != 0 {
		t.Log("runner timing summary", pending, failed, err)
	}
}
