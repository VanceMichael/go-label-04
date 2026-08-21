package search

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"strings"
)

type Query struct {
	Term        string
	Statuses    []domain.ShipmentStatus
	Origin      string
	Destination string
	Limit       int
}
type Item struct {
	ID          string
	Reference   string
	Origin      string
	Destination string
	Status      domain.ShipmentStatus
	Score       int
}

func Filter(items []Item, q Query) []Item {
	term := strings.ToLower(strings.TrimSpace(q.Term))
	status := map[domain.ShipmentStatus]bool{}
	for _, v := range q.Statuses {
		status[v] = true
	}
	out := make([]Item, 0)
	for _, v := range items {
		if term != "" && !strings.Contains(strings.ToLower(v.Reference), term) {
			continue
		}
		if q.Origin != "" && v.Origin != q.Origin {
			continue
		}
		if q.Destination != "" && v.Destination != q.Destination {
			continue
		}
		if len(status) > 0 && !status[v.Status] {
			continue
		}
		v.Score = score(v, term)
		out = append(out, v)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].ID < out[j].ID
		}
		return out[i].Score > out[j].Score
	})
	if q.Limit > 0 && len(out) > q.Limit {
		out = out[:q.Limit]
	}
	return out
}
func score(v Item, term string) int {
	if term == "" {
		return 0
	}
	score := 0
	if strings.EqualFold(v.Reference, term) {
		score += 100
	}
	if strings.Contains(strings.ToLower(v.Reference), term) {
		score += 10
	}
	return score
}
