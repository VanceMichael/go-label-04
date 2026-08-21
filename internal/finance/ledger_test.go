package finance

import (
	"testing"
	"time"
)

func TestLedgerReconcile(t *testing.T) {
	l := New()
	now := time.Now()
	if err := l.Post(Entry{ID: "1", TenantID: "T", ShipmentID: "S", Currency: "CNY", Debit: 100, PostedAt: now}); err != nil {
		t.Fatal(err)
	}
	if err := l.Post(Entry{ID: "2", TenantID: "T", ShipmentID: "S", Currency: "CNY", Credit: 100, PostedAt: now.Add(time.Second)}); err != nil {
		t.Fatal(err)
	}
	if l.Balance("T") != 0 {
		t.Fatal(l.Balance("T"))
	}
	if err := l.Reconcile("T"); err != nil {
		t.Fatal(err)
	}
}
