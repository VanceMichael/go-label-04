package report

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"time"
)

type ShipmentRow struct {
	Reference string
	Status    domain.ShipmentStatus
	WeightKg  int64
	CreatedAt time.Time
}
type Summary struct {
	Total    int
	WeightKg int64
	ByStatus map[domain.ShipmentStatus]int
	Latest   []ShipmentRow
}

func Build(rows []ShipmentRow, limit int) Summary {
	if limit < 1 {
		limit = 10
	}
	out := Summary{ByStatus: map[domain.ShipmentStatus]int{}}
	copyRows := append([]ShipmentRow(nil), rows...)
	sort.Slice(copyRows, func(i, j int) bool { return copyRows[i].CreatedAt.After(copyRows[j].CreatedAt) })
	for _, r := range copyRows {
		out.Total++
		out.WeightKg += r.WeightKg
		out.ByStatus[r.Status]++
	}
	if len(copyRows) > limit {
		copyRows = copyRows[:limit]
	}
	out.Latest = copyRows
	return out
}
func Window(rows []ShipmentRow, start, end time.Time) []ShipmentRow {
	out := make([]ShipmentRow, 0)
	for _, r := range rows {
		if !r.CreatedAt.Before(start) && r.CreatedAt.Before(end) {
			out = append(out, r)
		}
	}
	return out
}
