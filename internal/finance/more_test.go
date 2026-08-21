package finance

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestTenantBalance(t *testing.T) {
	l := New()
	now := time.Now()
	_ = l.Post(Entry{ID: "a", TenantID: "A", ShipmentID: "S", Currency: "CNY", Credit: 10, PostedAt: now})
	_ = l.Post(Entry{ID: "b", TenantID: "B", ShipmentID: "S", Currency: "CNY", Credit: 20, PostedAt: now})
	if l.Balance("A") != 10 || l.Balance("B") != 20 {
		t.Fatal(l.Balance("A"), l.Balance("B"))
	}
	if err := l.Reconcile("A"); err != nil {
		t.Fatal(err)
	}
	if err := l.Post(Entry{ID: "c", TenantID: "A", ShipmentID: "S", Currency: "CNY", Debit: 30, PostedAt: now}); err != nil {
		t.Fatal(err)
	}
	if l.Balance("A") != -20 {
		t.Fatal(l.Balance("A"))
	}
	_ = domain.ErrInvalid
}
