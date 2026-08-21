package validation

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

var codePattern = regexp.MustCompile(`^[A-Z]{3}$`)

func Email(v string) error {
	v = strings.TrimSpace(v)
	if _, err := mail.ParseAddress(v); err != nil {
		return domain.ErrInvalid
	}
	return nil
}
func AirportCode(v string) error {
	if !codePattern.MatchString(strings.ToUpper(strings.TrimSpace(v))) {
		return domain.ErrInvalid
	}
	return nil
}
func Reference(v string) error {
	v = strings.TrimSpace(v)
	if utf8.RuneCountInString(v) < 3 || utf8.RuneCountInString(v) > 64 {
		return domain.ErrInvalid
	}
	return nil
}
func Notes(v string) error {
	if utf8.RuneCountInString(v) > 1000 {
		return domain.ErrInvalid
	}
	return nil
}
func Required(values ...string) error {
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			return domain.ErrInvalid
		}
	}
	return nil
}
