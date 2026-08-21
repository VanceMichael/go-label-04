package ratelimit

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
	"time"
)

type Window struct {
	Started time.Time
	Count   int
}
type Limiter struct {
	mu       sync.Mutex
	limit    int
	duration time.Duration
	windows  map[string]Window
}

func New(limit int, d time.Duration) *Limiter {
	if limit < 1 {
		limit = 1
	}
	if d <= 0 {
		d = time.Minute
	}
	return &Limiter{limit: limit, duration: d, windows: map[string]Window{}}
}
func (l *Limiter) Allow(key string, now time.Time) bool {
	if key == "" {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	w, ok := l.windows[key]
	if !ok || !now.Before(w.Started.Add(l.duration)) {
		l.windows[key] = Window{Started: now, Count: 1}
		return true
	}
	if w.Count >= l.limit {
		return false
	}
	w.Count++
	l.windows[key] = w
	return true
}
func (l *Limiter) Reset(key string) { l.mu.Lock(); defer l.mu.Unlock(); delete(l.windows, key) }
func (l *Limiter) Remaining(key string, now time.Time) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	w, ok := l.windows[key]
	if !ok || !now.Before(w.Started.Add(l.duration)) {
		return l.limit
	}
	if w.Count >= l.limit {
		return 0
	}
	return l.limit - w.Count
}

var _ = domain.ErrInvalid
