package worker

import (
	"context"
	"errors"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/repository"
	"github.com/VanceMichael/go-base-airbridge/internal/repository/memory"
	"log/slog"
	"testing"
	"time"
)

type lifecycleContextKey struct{}

type lifecycleRepository struct {
	event             domain.OutboxEvent
	published         bool
	failed            bool
	finalContextValue string
	finalDeadline     bool
}

var _ repository.OutboxRepository = (*lifecycleRepository)(nil)

func (repository *lifecycleRepository) Enqueue(context.Context, domain.OutboxEvent) error {
	return nil
}
func (repository *lifecycleRepository) Claim(ctx context.Context, _ time.Time, _ int) ([]domain.OutboxEvent, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return []domain.OutboxEvent{repository.event}, nil
}
func (repository *lifecycleRepository) MarkPublished(ctx context.Context, _ string, _ time.Time) error {
	if err := repository.captureFinalContext(ctx); err != nil {
		return err
	}
	repository.published = true
	return nil
}
func (repository *lifecycleRepository) MarkFailed(ctx context.Context, _ string, _ time.Time, _ string) error {
	if err := repository.captureFinalContext(ctx); err != nil {
		return err
	}
	repository.failed = true
	return nil
}
func (repository *lifecycleRepository) Summary(context.Context, string) (int, int, error) {
	return 0, 0, nil
}
func (repository *lifecycleRepository) captureFinalContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	repository.finalContextValue, _ = ctx.Value(lifecycleContextKey{}).(string)
	_, repository.finalDeadline = ctx.Deadline()
	if repository.finalContextValue == "" {
		return errors.New("finalization context lost operation metadata")
	}
	if !repository.finalDeadline {
		return errors.New("finalization context has no deadline")
	}
	return nil
}

func TestShutdownFinalizesClaimedEventsWithBoundedContext(t *testing.T) {
	brokerUnavailable := errors.New("broker unavailable during shutdown")
	for _, testCase := range []struct {
		name       string
		publishErr error
		published  bool
		failed     bool
	}{
		{name: "publisher completed", published: true},
		{name: "publisher interrupted", publishErr: brokerUnavailable, failed: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			event := domain.OutboxEvent{ID: "event-shutdown", TenantID: "tenant-a", Topic: "shipment.booked", AggregateID: "shipment-a", AvailableAt: time.Now().UTC()}
			repository := &lifecycleRepository{event: event}
			ctx, cancel := context.WithCancel(context.WithValue(context.Background(), lifecycleContextKey{}, "cycle-shutdown"))
			publisher := PublisherFunc(func(context.Context, domain.OutboxEvent) error {
				cancel()
				return testCase.publishErr
			})
			runner := NewOutboxRunnerWithPublisher(repository, publisher, slog.Default())

			if err := runner.RunOnce(ctx, time.Now().UTC()); err != nil {
				t.Fatalf("RunOnce() error = %v", err)
			}
			if repository.published != testCase.published || repository.failed != testCase.failed {
				t.Fatalf("published=%v failed=%v", repository.published, repository.failed)
			}
			if repository.finalContextValue != "cycle-shutdown" || !repository.finalDeadline {
				t.Fatalf("final context value=%q deadline=%v", repository.finalContextValue, repository.finalDeadline)
			}
		})
	}
}

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
