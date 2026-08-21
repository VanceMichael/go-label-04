package rules

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strings"
	"time"
)

type Rule struct {
	Name    string
	Enabled bool
	Check   func(context map[string]any) error
}
type Set struct{ rules []Rule }

func New() *Set { return &Set{rules: []Rule{}} }
func (s *Set) Add(r Rule) error {
	if strings.TrimSpace(r.Name) == "" || r.Check == nil {
		return domain.ErrInvalid
	}
	for _, old := range s.rules {
		if old.Name == r.Name {
			return domain.ErrConflict
		}
	}
	s.rules = append(s.rules, r)
	return nil
}
func (s *Set) Evaluate(values map[string]any) error {
	for _, r := range s.rules {
		if r.Enabled {
			if err := r.Check(values); err != nil {
				return err
			}
		}
	}
	return nil
}
func Before(deadline time.Time) Rule {
	return Rule{Name: "before_deadline", Enabled: true, Check: func(values map[string]any) error {
		raw, ok := values["at"].(time.Time)
		if !ok || !raw.Before(deadline) {
			return domain.ErrExpired
		}
		return nil
	}}
}
