package alerts

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestAlertRegistry(t *testing.T) {
	r := New()
	now := time.Now()
	a := Alert{ID: "A", TenantID: "T", Code: "late", Message: "late", OpenedAt: now}
	if err := r.Open(a); err != nil {
		t.Fatal(err)
	}
	if len(r.Active("T")) != 1 {
		t.Fatal("missing alert")
	}
	if err := r.Open(a); err != domain.ErrConflict {
		t.Fatal(err)
	}
	if err := r.Close("A", now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if len(r.Active("T")) != 0 {
		t.Fatal("alert still active")
	}
}
