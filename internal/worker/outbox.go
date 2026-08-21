package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/repository"
)

type Publisher interface {
	Publish(context.Context, domain.OutboxEvent) error
}

type PublisherFunc func(context.Context, domain.OutboxEvent) error

func (publish PublisherFunc) Publish(ctx context.Context, event domain.OutboxEvent) error {
	return publish(ctx, event)
}

type OutboxRunner struct {
	repo            repository.OutboxRepository
	publisher       Publisher
	log             *slog.Logger
	done            chan struct{}
	wg              sync.WaitGroup
	interval        time.Duration
	batchSize       int
	finalizeTimeout time.Duration
}

func NewOutboxRunner(repo repository.OutboxRepository, log *slog.Logger) *OutboxRunner {
	return NewOutboxRunnerWithPublisher(repo, PublisherFunc(logPublish), log)
}

func NewOutboxRunnerWithPublisher(repo repository.OutboxRepository, publisher Publisher, log *slog.Logger) *OutboxRunner {
	return &OutboxRunner{
		repo:            repo,
		publisher:       publisher,
		log:             log,
		done:            make(chan struct{}),
		interval:        250 * time.Millisecond,
		batchSize:       20,
		finalizeTimeout: 5 * time.Second,
	}
}

func (runner *OutboxRunner) Run(ctx context.Context) {
	runner.wg.Add(1)
	defer runner.wg.Done()
	ticker := time.NewTicker(runner.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-runner.done:
			return
		case now := <-ticker.C:
			if err := runner.RunOnce(ctx, now); err != nil {
				runner.log.Error("outbox cycle failed", "error", err)
			}
		}
	}
}

func (runner *OutboxRunner) Wait() {
	close(runner.done)
	runner.wg.Wait()
}

func (runner *OutboxRunner) RunOnce(ctx context.Context, now time.Time) error {
	if runner.repo == nil || runner.publisher == nil {
		return fmt.Errorf("outbox repository and publisher are required")
	}
	events, err := runner.repo.Claim(ctx, now, runner.batchSize)
	if err != nil {
		return fmt.Errorf("claim outbox events: %w", err)
	}
	finalizer := deliveryFinalizer{repo: runner.repo, timeout: runner.finalizeTimeout}
	for _, event := range events {
		publishErr := runner.publisher.Publish(ctx, event)
		if publishErr != nil {
			if err := finalizer.failed(ctx, event, now, publishErr); err != nil {
				return err
			}
			continue
		}
		if err := finalizer.published(ctx, event, now); err != nil {
			return err
		}
	}
	return nil
}

func logPublish(ctx context.Context, _ domain.OutboxEvent) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
