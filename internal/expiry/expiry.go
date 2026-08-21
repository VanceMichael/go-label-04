package expiry

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"time"
)

type Value struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
	Grace     time.Duration
}

func New(issued time.Time, ttl, grace time.Duration) (Value, error) {
	if issued.IsZero() || ttl <= 0 || grace < 0 {
		return Value{}, domain.ErrInvalid
	}
	return Value{IssuedAt: issued, ExpiresAt: issued.Add(ttl), Grace: grace}, nil
}
func (v Value) Active(at time.Time) bool { return at.Before(v.ExpiresAt) }
func (v Value) WithinGrace(at time.Time) bool {
	return !at.Before(v.ExpiresAt) && at.Before(v.ExpiresAt.Add(v.Grace))
}
func (v Value) Remaining(at time.Time) time.Duration {
	if !at.Before(v.ExpiresAt) {
		return 0
	}
	return v.ExpiresAt.Sub(at)
}
