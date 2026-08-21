package healthcheck

import (
	"sort"
	"strings"
)

func Healthy(results []Result) bool {
	if len(results) == 0 {
		return false
	}
	for _, r := range results {
		if !r.OK {
			return false
		}
	}
	return true
}
func Failed(results []Result) []Result {
	out := make([]Result, 0)
	for _, r := range results {
		if !r.OK {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func Render(results []Result) string {
	parts := make([]string, 0, len(results))
	for _, r := range results {
		state := "ok"
		if !r.OK {
			state = "failed"
		}
		parts = append(parts, r.Name+":"+state)
	}
	return strings.Join(parts, ",")
}
