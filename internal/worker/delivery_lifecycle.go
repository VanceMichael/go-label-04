package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/repository"
)

type deliveryFinalizer struct {
	repo    repository.OutboxRepository
	timeout time.Duration
}

func (finalizer deliveryFinalizer) published(ctx context.Context, event domain.OutboxEvent, at time.Time) error {
	if event.ID == "" {
		return fmt.Errorf("published outbox event has no identity")
	}
	if finalizer.timeout <= 0 {
		return fmt.Errorf("outbox finalization timeout must be positive")
	}
	if err := finalizer.repo.MarkPublished(ctx, event.ID, at); err != nil {
		return fmt.Errorf("finalize published outbox event %s: %w", event.ID, err)
	}
	return nil
}

func (finalizer deliveryFinalizer) failed(ctx context.Context, event domain.OutboxEvent, at time.Time, publishErr error) error {
	if event.ID == "" {
		return fmt.Errorf("failed outbox event has no identity")
	}
	if publishErr == nil {
		return fmt.Errorf("failed outbox event %s has no publish error", event.ID)
	}
	if finalizer.timeout <= 0 {
		return fmt.Errorf("outbox finalization timeout must be positive")
	}
	if err := finalizer.repo.MarkFailed(ctx, event.ID, at, publishErr.Error()); err != nil {
		return fmt.Errorf("release failed outbox event %s: %w", event.ID, err)
	}
	return nil
}
