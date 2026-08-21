package notification

import (
	"context"
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strings"
	"sync"
	"time"
)

type Message struct {
	ID        string
	TenantID  string
	Recipient string
	Channel   string
	Subject   string
	Body      string
	CreatedAt time.Time
}
type Sender interface {
	Send(context.Context, Message) error
}
type MemorySender struct {
	mu       sync.Mutex
	Sent     []Message
	Failures int
}

func (m *MemorySender) Send(ctx context.Context, v Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Failures > 0 {
		m.Failures--
		return fmt.Errorf("temporary notification failure")
	}
	m.Sent = append(m.Sent, v)
	return nil
}

type Router struct{ senders map[string]Sender }

func New() *Router { return &Router{senders: map[string]Sender{}} }
func (r *Router) Register(channel string, s Sender) error {
	if strings.TrimSpace(channel) == "" || s == nil {
		return domain.ErrInvalid
	}
	if _, ok := r.senders[channel]; ok {
		return domain.ErrConflict
	}
	r.senders[channel] = s
	return nil
}
func (r *Router) Deliver(ctx context.Context, v Message) error {
	if v.ID == "" || v.Recipient == "" || v.Body == "" {
		return domain.ErrInvalid
	}
	s, ok := r.senders[v.Channel]
	if !ok {
		return domain.ErrNotFound
	}
	return s.Send(ctx, v)
}
