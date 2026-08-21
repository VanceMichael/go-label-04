package sortx

import (
	"sort"
	"time"
)

type Record struct {
	ID    string
	Time  time.Time
	Score int64
}

func ByTime(records []Record, descending bool) []Record {
	out := append([]Record(nil), records...)
	sort.SliceStable(out, func(i, j int) bool {
		if descending {
			return out[i].Time.After(out[j].Time)
		}
		return out[i].Time.Before(out[j].Time)
	})
	return out
}
func ByScore(records []Record) []Record {
	out := append([]Record(nil), records...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].ID < out[j].ID
		}
		return out[i].Score > out[j].Score
	})
	return out
}
