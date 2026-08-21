package retention

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"time"
)

type Policy struct {
	Name      string
	Keep      time.Duration
	LegalHold bool
}
type Record struct {
	ID        string
	CreatedAt time.Time
	DeletedAt *time.Time
	Hold      bool
}

func Expired(p Policy, r Record, now time.Time) bool {
	if p.Name == "" || p.Keep <= 0 || r.CreatedAt.IsZero() || r.Hold || p.LegalHold {
		return false
	}
	return !now.Before(r.CreatedAt.Add(p.Keep))
}
func MarkDeleted(p Policy, r Record, now time.Time) (Record, error) {
	if !Expired(p, r, now) {
		return Record{}, domain.ErrState
	}
	r.DeletedAt = &now
	return r, nil
}
func Purgeable(r Record, now time.Time) bool {
	return r.DeletedAt != nil && now.After(r.DeletedAt.Add(24*time.Hour))
}
