package sequence

import (
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sync"
	"time"
)

type Generator struct {
	mu     sync.Mutex
	prefix string
	day    string
	next   int64
}

func New(prefix string) *Generator { return &Generator{prefix: prefix} }
func (g *Generator) Next(now time.Time) (string, error) {
	if g.prefix == "" || now.IsZero() {
		return "", domain.ErrInvalid
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	day := now.UTC().Format("20060102")
	if day != g.day {
		g.day = day
		g.next = 0
	}
	g.next++
	return fmt.Sprintf("%s-%s-%06d", g.prefix, day, g.next), nil
}
func (g *Generator) Peek() int64 { g.mu.Lock(); defer g.mu.Unlock(); return g.next + 1 }
