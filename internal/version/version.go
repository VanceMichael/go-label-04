package version

import (
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strings"
)

type Value struct {
	Major int
	Minor int
	Patch int
	Pre   string
}

func Parse(raw string) (Value, error) {
	raw = strings.TrimPrefix(strings.TrimSpace(raw), "v")
	parts := strings.SplitN(raw, "-", 2)
	nums := strings.Split(parts[0], ".")
	if len(nums) != 3 {
		return Value{}, domain.ErrInvalid
	}
	var v Value
	if _, err := fmt.Sscanf(strings.Join(nums, "."), "%d.%d.%d", &v.Major, &v.Minor, &v.Patch); err != nil {
		return Value{}, domain.ErrInvalid
	}
	if len(parts) == 2 {
		v.Pre = parts[1]
	}
	if v.Major < 0 || v.Minor < 0 || v.Patch < 0 {
		return Value{}, domain.ErrInvalid
	}
	return v, nil
}
func (v Value) String() string {
	out := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		out += "-" + v.Pre
	}
	return out
}
func (v Value) AtLeast(other Value) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch >= other.Patch
}
