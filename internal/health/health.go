package health

import (
	"context"
	"sync"
	"time"
)

type Check func(context.Context) error
type Status struct {
	Name     string
	OK       bool
	Error    string
	Duration time.Duration
}
type Registry struct {
	mu     sync.RWMutex
	checks map[string]Check
}

func New() *Registry { return &Registry{checks: map[string]Check{}} }
func (r *Registry) Register(name string, fn Check) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks[name] = fn
}
func (r *Registry) Run(ctx context.Context) []Status {
	r.mu.RLock()
	checks := map[string]Check{}
	for k, v := range r.checks {
		checks[k] = v
	}
	r.mu.RUnlock()
	out := make([]Status, 0, len(checks))
	for name, fn := range checks {
		start := time.Now()
		err := fn(ctx)
		s := Status{Name: name, OK: err == nil, Duration: time.Since(start)}
		if err != nil {
			s.Error = err.Error()
		}
		out = append(out, s)
	}
	return out
}
