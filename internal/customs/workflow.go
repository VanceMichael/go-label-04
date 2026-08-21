package customs

import (
	"context"
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strings"
	"sync"
	"time"
)

type Document struct {
	Number    string
	Kind      string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Hash      string
}
type Case struct {
	ID         string
	ShipmentID string
	TenantID   string
	Documents  []Document
	Status     domain.CustomsStatus
	Notes      []string
	UpdatedAt  time.Time
}
type Workflow struct {
	mu    sync.RWMutex
	cases map[string]Case
}

func NewWorkflow() *Workflow { return &Workflow{cases: map[string]Case{}} }
func (w *Workflow) Open(_ context.Context, c Case) error {
	if c.ID == "" || c.ShipmentID == "" || c.TenantID == "" {
		return domain.ErrInvalid
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.cases[c.ShipmentID]; ok {
		return domain.ErrConflict
	}
	c.Status = domain.CustomsPending
	c.UpdatedAt = time.Now().UTC()
	w.cases[c.ShipmentID] = c
	return nil
}
func (w *Workflow) Attach(_ context.Context, shipment string, d Document) error {
	if strings.TrimSpace(d.Number) == "" || d.ExpiresAt.Before(d.IssuedAt) {
		return domain.ErrInvalid
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	c, ok := w.cases[shipment]
	if !ok {
		return domain.ErrNotFound
	}
	for _, old := range c.Documents {
		if old.Number == d.Number {
			return domain.ErrConflict
		}
	}
	c.Documents = append(c.Documents, d)
	c.UpdatedAt = time.Now().UTC()
	w.cases[shipment] = c
	return nil
}
func (w *Workflow) Review(_ context.Context, shipment, actor string, now time.Time) (Case, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	c, ok := w.cases[shipment]
	if !ok {
		return Case{}, domain.ErrNotFound
	}
	if actor == "" {
		return Case{}, domain.ErrForbidden
	}
	if len(c.Documents) == 0 {
		return Case{}, fmt.Errorf("%w: documents missing", domain.ErrState)
	}
	for _, d := range c.Documents {
		if !now.Before(d.ExpiresAt) {
			c.Status = domain.CustomsHeld
			c.Notes = append(c.Notes, "expired document "+d.Number)
			w.cases[shipment] = c
			return c, domain.ErrExpired
		}
	}
	c.Status = domain.CustomsReleased
	c.Notes = append(c.Notes, "reviewed by "+actor)
	c.UpdatedAt = now
	w.cases[shipment] = c
	return c, nil
}
func (w *Workflow) Get(shipment string) (Case, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	c, ok := w.cases[shipment]
	if !ok {
		return Case{}, domain.ErrNotFound
	}
	c.Documents = append([]Document(nil), c.Documents...)
	c.Notes = append([]string(nil), c.Notes...)
	return c, nil
}
