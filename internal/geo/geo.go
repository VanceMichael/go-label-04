package geo

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"math"
)

type Point struct {
	Lat float64
	Lon float64
}

func Valid(p Point) bool {
	return p.Lat >= -90 && p.Lat <= 90 && p.Lon >= -180 && p.Lon <= 180 && !math.IsNaN(p.Lat) && !math.IsNaN(p.Lon)
}
func Distance(a, b Point) (float64, error) {
	if !Valid(a) || !Valid(b) {
		return 0, domain.ErrInvalid
	}
	const rad = math.Pi / 180
	lat1 := a.Lat * rad
	lat2 := b.Lat * rad
	dlat := (b.Lat - a.Lat) * rad
	dlon := (b.Lon - a.Lon) * rad
	h := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	return 6371 * 2 * math.Asin(math.Sqrt(h)), nil
}
