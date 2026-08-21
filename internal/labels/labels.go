package labels

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"strings"
)

type Set map[string]string

func Clean(in Set) Set {
	out := Set{}
	for k, v := range in {
		k = strings.TrimSpace(strings.ToLower(k))
		v = strings.TrimSpace(v)
		if k != "" && v != "" {
			out[k] = v
		}
	}
	return out
}
func Merge(base, override Set) Set {
	out := Set{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		out[k] = v
	}
	return out
}
func Keys(in Set) []string {
	out := make([]string, 0, len(in))
	for k := range in {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func Validate(in Set) error {
	for k, v := range in {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			return domain.ErrInvalid
		}
	}
	return nil
}
