package ratios

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"math"
)

func Percent(part, total int64) (float64, error) {
	if total <= 0 || part < 0 || part > total {
		return 0, domain.ErrInvalid
	}
	return float64(part) * 100 / float64(total), nil
}
func Round(v float64, places int) float64 {
	if places < 0 {
		places = 0
	}
	factor := math.Pow10(places)
	return math.Round(v*factor) / factor
}
func Weighted(values []int64, weights []int64) (float64, error) {
	if len(values) == 0 || len(values) != len(weights) {
		return 0, domain.ErrInvalid
	}
	var sum, weight int64
	for i, v := range values {
		if v < 0 || weights[i] < 0 {
			return 0, domain.ErrInvalid
		}
		sum += v * weights[i]
		weight += weights[i]
	}
	if weight == 0 {
		return 0, domain.ErrInvalid
	}
	return float64(sum) / float64(weight), nil
}
