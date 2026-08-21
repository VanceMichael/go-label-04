package metrics

import (
	"sync"
	"time"
)

type Counter struct {
	mu     sync.Mutex
	values map[string]int64
}

func NewCounter() *Counter         { return &Counter{values: map[string]int64{}} }
func (c *Counter) Inc(name string) { c.Add(name, 1) }
func (c *Counter) Add(name string, n int64) {
	if name == "" {
		return
	}
	c.mu.Lock()
	c.values[name] += n
	c.mu.Unlock()
}
func (c *Counter) Get(name string) int64 { c.mu.Lock(); defer c.mu.Unlock(); return c.values[name] }
func (c *Counter) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := map[string]int64{}
	for k, v := range c.values {
		out[k] = v
	}
	return out
}

type Timer struct{ started time.Time }

func StartTimer() Timer                { return Timer{started: time.Now()} }
func (t Timer) Elapsed() time.Duration { return time.Since(t.started) }
