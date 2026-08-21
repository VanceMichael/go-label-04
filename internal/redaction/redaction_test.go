package redaction

import (
	"testing"
)

func TestJSONRedactsSecrets(t *testing.T) {
	raw, err := JSON([]byte(`{"token":"abc","nested":{"password":"x"},"name":"ok"}`))
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	if s == "" || contains(s, "abc") || contains(s, "x") || !contains(s, "REDACTED") {
		t.Fatal(s)
	}
}
func contains(s, w string) bool {
	for i := 0; i+len(w) <= len(s); i++ {
		if s[i:i+len(w)] == w {
			return true
		}
	}
	return false
}
