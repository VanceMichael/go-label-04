package locale

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"time"
)

func Parse(name, value string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(time.RFC3339, value, loc)
	return t.UTC(), err
}
func Window(name string, start, end time.Time) (time.Duration, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return 0, err
	}
	a := start.In(loc)
	b := end.In(loc)
	if !b.After(a) {
		return 0, domain.ErrInvalid
	}
	return b.Sub(a), nil
}
func BusinessDay(name string, t time.Time) bool {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return false
	}
	local := t.In(loc)
	return local.Weekday() != time.Saturday && local.Weekday() != time.Sunday
}
