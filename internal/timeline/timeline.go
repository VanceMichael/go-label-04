package timeline

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"sort"
	"time"
)

type Event struct {
	Kind    string
	At      time.Time
	Actor   string
	Details string
}
type Timeline struct{ Events []Event }

func New() *Timeline { return &Timeline{Events: []Event{}} }
func (t *Timeline) Append(e Event) error {
	if e.Kind == "" || e.Actor == "" || e.At.IsZero() {
		return domain.ErrInvalid
	}
	for _, old := range t.Events {
		if old.Kind == e.Kind && old.At.Equal(e.At) {
			return domain.ErrConflict
		}
	}
	t.Events = append(t.Events, e)
	sort.SliceStable(t.Events, func(i, j int) bool { return t.Events[i].At.Before(t.Events[j].At) })
	return nil
}
func (t Timeline) Between(start, end time.Time) []Event {
	out := make([]Event, 0)
	for _, e := range t.Events {
		if !e.At.Before(start) && e.At.Before(end) {
			out = append(out, e)
		}
	}
	return out
}
func (t Timeline) Latest(kind string) (Event, error) {
	for i := len(t.Events) - 1; i >= 0; i-- {
		if t.Events[i].Kind == kind {
			return t.Events[i], nil
		}
	}
	return Event{}, domain.ErrNotFound
}
func (t Timeline) Duration(startKind, endKind string) (time.Duration, error) {
	a, err := t.Latest(startKind)
	if err != nil {
		return 0, err
	}
	b, err := t.Latest(endKind)
	if err != nil {
		return 0, err
	}
	if b.At.Before(a.At) {
		return 0, domain.ErrState
	}
	return b.At.Sub(a.At), nil
}
