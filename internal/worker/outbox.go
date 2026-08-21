package worker

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/repository"
	"log/slog"
	"sync"
	"time"
)

type OutboxRunner struct {
	repo repository.OutboxRepository
	log  *slog.Logger
	done chan struct{}
	wg   sync.WaitGroup
}

func NewOutboxRunner(repo repository.OutboxRepository, log *slog.Logger) *OutboxRunner {
	return &OutboxRunner{repo: repo, log: log, done: make(chan struct{})}
}
func (r *OutboxRunner) Run(ctx context.Context) {
	r.wg.Add(1)
	defer r.wg.Done()
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-r.done:
			return
		case now := <-ticker.C:
			r.tick(ctx, now)
		}
	}
}
func (r *OutboxRunner) Wait() { close(r.done); r.wg.Wait() }
func (r *OutboxRunner) tick(ctx context.Context, now time.Time) {
	events, err := r.repo.Claim(ctx, now, 20)
	if err != nil {
		r.log.Error("outbox claim failed", "error", err)
		return
	}
	for _, event := range events {
		if err := r.publish(ctx, event); err != nil {
			_ = r.repo.MarkFailed(ctx, event.ID, now, err.Error())
			continue
		}
		_ = r.repo.MarkPublished(ctx, event.ID, now)
	}
}
func (r *OutboxRunner) publish(ctx context.Context, event domain.OutboxEvent) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
