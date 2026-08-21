package masking

import (
	"strings"
	"unicode/utf8"
)

func Email(v string) string {
	at := strings.Index(v, "@")
	if at <= 1 {
		return "***"
	}
	return v[:1] + "***" + v[at:]
}
func Token(v string) string {
	if utf8.RuneCountInString(v) <= 8 {
		return "***"
	}
	return v[:4] + "..." + v[len(v)-4:]
}
func Phone(v string) string {
	r := []rune(v)
	if len(r) <= 4 {
		return "***"
	}
	return string(r[:2]) + "***" + string(r[len(r)-2:])
}
func Secret(v string) string {
	if strings.TrimSpace(v) == "" {
		return ""
	}
	return "[REDACTED]"
}
