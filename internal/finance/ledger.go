package finance

import (
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"sync"
	"time"
)

type Entry struct {
	ID         string
	TenantID   string
	ShipmentID string
	Currency   string
	Debit      int64
	Credit     int64
	Memo       string
	PostedAt   time.Time
}
type Ledger struct {
	mu       sync.RWMutex
	entries  []Entry
	balances map[string]int64
}

func New() *Ledger { return &Ledger{entries: make([]Entry, 0), balances: make(map[string]int64)} }
func (l *Ledger) Post(e Entry) error {
	if e.ID == "" || e.TenantID == "" || e.ShipmentID == "" || e.Currency == "" || e.PostedAt.IsZero() {
		return domain.ErrInvalid
	}
	if e.Debit < 0 || e.Credit < 0 || e.Debit == 0 && e.Credit == 0 {
		return domain.ErrInvalid
	}
	if e.Debit > 0 && e.Credit > 0 {
		return domain.ErrInvalid
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, old := range l.entries {
		if old.ID == e.ID {
			return domain.ErrConflict
		}
	}
	amount := e.Credit - e.Debit
	l.balances[e.TenantID] += amount
	l.entries = append(l.entries, e)
	return nil
}
func (l *Ledger) Balance(tenant string) int64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.balances[tenant]
}
func (l *Ledger) Entries(tenant string) []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, 0)
	for _, e := range l.entries {
		if e.TenantID == tenant {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PostedAt.Before(out[j].PostedAt) })
	return out
}
func (l *Ledger) Reconcile(tenant string) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var total int64
	for _, e := range l.entries {
		if e.TenantID == tenant {
			total += e.Credit - e.Debit
		}
	}
	if total != l.balances[tenant] {
		return fmt.Errorf("%w: ledger balance", domain.ErrConflict)
	}
	return nil
}
