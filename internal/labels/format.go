package labels

import (
	"sort"
	"strings"
)

func Format(in Set) string {
	keys := Keys(in)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+in[k])
	}
	return strings.Join(parts, ",")
}
func Parse(raw string) Set {
	out := Set{}
	for _, part := range strings.Split(raw, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return out
}
func SortedValues(in Set) []string {
	keys := Keys(in)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, in[k])
	}
	sort.Strings(out)
	return out
}
