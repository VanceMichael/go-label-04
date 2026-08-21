package pricing

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"math"
	"sort"
	"time"
)

type Rate struct {
	Origin        string
	Destination   string
	MinKg         int64
	MaxKg         int64
	CentsPerKg    int64
	FuelCents     int64
	EffectiveFrom time.Time
	EffectiveTo   time.Time
}
type Quote struct {
	Currency   string
	BaseCents  int64
	FuelCents  int64
	Surcharges int64
	TotalCents int64
	Rate       Rate
}
type Table struct{ rates []Rate }

func NewTable() *Table { return &Table{rates: []Rate{}} }
func (t *Table) Add(r Rate) error {
	if r.Origin == "" || r.Destination == "" || r.MinKg <= 0 || r.MaxKg < r.MinKg || r.CentsPerKg < 0 || !r.EffectiveTo.After(r.EffectiveFrom) {
		return domain.ErrInvalid
	}
	t.rates = append(t.rates, r)
	return nil
}
func (t *Table) Quote(origin, destination string, weight int64, at time.Time) (Quote, error) {
	if weight <= 0 {
		return Quote{}, domain.ErrInvalid
	}
	candidates := make([]Rate, 0)
	for _, r := range t.rates {
		if r.Origin == origin && r.Destination == destination && weight >= r.MinKg && weight <= r.MaxKg && !at.Before(r.EffectiveFrom) && at.Before(r.EffectiveTo) {
			candidates = append(candidates, r)
		}
	}
	if len(candidates) == 0 {
		return Quote{}, domain.ErrNotFound
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].EffectiveFrom.After(candidates[j].EffectiveFrom) })
	r := candidates[0]
	base := r.CentsPerKg * weight
	fuel := r.FuelCents
	total := base + fuel
	if total < 0 || total > math.MaxInt64 {
		return Quote{}, domain.ErrInvalid
	}
	return Quote{Currency: "CNY", BaseCents: base, FuelCents: fuel, TotalCents: total, Rate: r}, nil
}
