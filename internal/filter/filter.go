package filter

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strings"
	"time"
)

type Shipment struct {
	ID          string
	TenantID    string
	Reference   string
	Status      domain.ShipmentStatus
	Origin      string
	Destination string
	WeightKg    int64
	CreatedAt   time.Time
}
type Input struct {
	TenantID    string
	Text        string
	Statuses    []domain.ShipmentStatus
	Origin      string
	Destination string
	MinWeight   int64
	MaxWeight   int64
	Since       *time.Time
	Until       *time.Time
}

func Apply(items []Shipment, in Input) []Shipment {
	out := make([]Shipment, 0)
	statuses := map[domain.ShipmentStatus]bool{}
	for _, s := range in.Statuses {
		statuses[s] = true
	}
	text := strings.ToLower(strings.TrimSpace(in.Text))
	for _, item := range items {
		if in.TenantID != "" && item.TenantID != in.TenantID {
			continue
		}
		if text != "" && !strings.Contains(strings.ToLower(item.Reference), text) {
			continue
		}
		if len(statuses) > 0 && !statuses[item.Status] {
			continue
		}
		if in.Origin != "" && item.Origin != in.Origin {
			continue
		}
		if in.Destination != "" && item.Destination != in.Destination {
			continue
		}
		if in.MinWeight > 0 && item.WeightKg < in.MinWeight {
			continue
		}
		if in.MaxWeight > 0 && item.WeightKg > in.MaxWeight {
			continue
		}
		if in.Since != nil && item.CreatedAt.Before(*in.Since) {
			continue
		}
		if in.Until != nil && !item.CreatedAt.Before(*in.Until) {
			continue
		}
		out = append(out, item)
	}
	return out
}
func Normalize(in Input) Input {
	in.Text = strings.TrimSpace(in.Text)
	in.Origin = strings.ToUpper(strings.TrimSpace(in.Origin))
	in.Destination = strings.ToUpper(strings.TrimSpace(in.Destination))
	if in.MinWeight < 0 {
		in.MinWeight = 0
	}
	if in.MaxWeight < 0 {
		in.MaxWeight = 0
	}
	return in
}
