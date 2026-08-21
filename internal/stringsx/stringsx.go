package stringsx

import (
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strings"
	"unicode"
)

func Clean(v string) string { return strings.Join(strings.Fields(strings.TrimSpace(v)), " ") }
func Slug(v string) string {
	v = Clean(strings.ToLower(v))
	var b strings.Builder
	dash := false
	for _, r := range v {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
func Limit(v string, n int) string {
	if n < 1 {
		return ""
	}
	r := []rune(v)
	if len(r) <= n {
		return v
	}
	return string(r[:n])
}
func EqualFoldNonEmpty(a, b string) error {
	if Clean(a) == "" || Clean(b) == "" || !strings.EqualFold(Clean(a), Clean(b)) {
		return domain.ErrInvalid
	}
	return nil
}
