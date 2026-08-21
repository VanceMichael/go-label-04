package approvals

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
	"time"
)

func TestApprovalQueue(t *testing.T) {
	q := New()
	now := time.Now()
	r := Request{ID: "R", TenantID: "T", ObjectID: "S", RequestedBy: "u", CreatedAt: now}
	if err := q.Submit(r); err != nil {
		t.Fatal(err)
	}
	if err := q.Decide("R", "u", true, now.Add(time.Minute)); err != domain.ErrState {
		t.Fatal(err)
	}
	if err := q.Decide("R", "manager", true, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	v, _ := q.Get("R")
	if v.Status != "approved" {
		t.Fatal(v)
	}
}
