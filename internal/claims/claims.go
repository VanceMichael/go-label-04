package claims

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strings"
	"time"
)

type Claim struct {
	ID         string
	TenantID   string
	ShipmentID string
	FiledBy    string
	Reason     string
	Status     string
	FiledAt    time.Time
	ResolvedAt *time.Time
}

func Open(c Claim) error {
	if c.ID == "" || c.TenantID == "" || c.ShipmentID == "" || c.FiledBy == "" || strings.TrimSpace(c.Reason) == "" || c.FiledAt.IsZero() {
		return domain.ErrInvalid
	}
	return nil
}
func Resolve(c Claim, at time.Time) (Claim, error) {
	if c.Status == "resolved" || c.Status == "rejected" {
		return Claim{}, domain.ErrState
	}
	if at.Before(c.FiledAt) {
		return Claim{}, domain.ErrInvalid
	}
	c.Status = "resolved"
	c.ResolvedAt = &at
	return c, nil
}
func Reject(c Claim, at time.Time) (Claim, error) {
	if c.Status == "resolved" || c.Status == "rejected" {
		return Claim{}, domain.ErrState
	}
	if at.Before(c.FiledAt) {
		return Claim{}, domain.ErrInvalid
	}
	c.Status = "rejected"
	c.ResolvedAt = &at
	return c, nil
}
func IsOpen(c Claim) bool { return c.Status == "open" || c.Status == "investigating" }
