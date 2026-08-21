package bulk

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
)

type Result[T any] struct {
	Index int
	Value T
	Err   error
}
type Processor[T any] struct {
	StopOnError bool
	Workers     int
	Handle      func(context.Context, T) error
}

func (p Processor[T]) Run(ctx context.Context, values []T) []Result[T] {
	if p.Workers < 1 {
		p.Workers = 1
	}
	if p.Handle == nil {
		return []Result[T]{{Err: domain.ErrInvalid}}
	}
	jobs := make(chan int)
	results := make(chan Result[T], len(values))
	var wg sync.WaitGroup
	var stop sync.Once
	cancelled := false
	var mu sync.Mutex
	for worker := 0; worker < p.Workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				mu.Lock()
				halt := cancelled
				mu.Unlock()
				if halt {
					results <- Result[T]{Index: idx, Err: context.Canceled}
					continue
				}
				err := p.Handle(ctx, values[idx])
				results <- Result[T]{Index: idx, Value: values[idx], Err: err}
				if err != nil && p.StopOnError {
					stop.Do(func() { mu.Lock(); cancelled = true; mu.Unlock() })
				}
			}
		}()
	}
	for i := range values {
		select {
		case <-ctx.Done():
			mu.Lock()
			cancelled = true
			mu.Unlock()
		case jobs <- i:
		}
	}
	close(jobs)
	wg.Wait()
	close(results)
	out := make([]Result[T], 0, len(values))
	for r := range results {
		out = append(out, r)
	}
	return out
}
