package workerpool

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
)

type Job func(context.Context) error
type Pool struct {
	workers int
	jobs    chan Job
	done    chan struct{}
	wg      sync.WaitGroup
	once    sync.Once
}

func New(workers, queue int) *Pool {
	if workers < 1 {
		workers = 1
	}
	if queue < workers {
		queue = workers
	}
	return &Pool{workers: workers, jobs: make(chan Job, queue), done: make(chan struct{})}
}
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-p.done:
					return
				case job := <-p.jobs:
					if job != nil {
						_ = job(ctx)
					}
				}
			}
		}()
	}
}
func (p *Pool) Submit(ctx context.Context, job Job) error {
	if job == nil {
		return domain.ErrInvalid
	}
	select {
	case p.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.done:
		return domain.ErrState
	}
}
func (p *Pool) Stop() { p.once.Do(func() { close(p.done); p.wg.Wait() }) }
