package healthcheck

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"time"
)

type Probe struct {
	Name    string
	Timeout time.Duration
	Run     func(context.Context) error
}
type Result struct {
	Name    string
	OK      bool
	Elapsed time.Duration
	Error   string
}

func Execute(ctx context.Context, p Probe) Result {
	start := time.Now()
	r := Result{Name: p.Name}
	if p.Name == "" || p.Run == nil {
		r.Error = domain.ErrInvalid.Error()
		return r
	}
	child, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()
	if err := p.Run(child); err != nil {
		r.Error = err.Error()
	} else {
		r.OK = true
	}
	r.Elapsed = time.Since(start)
	return r
}
func All(ctx context.Context, probes []Probe) []Result {
	out := make([]Result, 0, len(probes))
	for _, p := range probes {
		out = append(out, Execute(ctx, p))
	}
	return out
}
