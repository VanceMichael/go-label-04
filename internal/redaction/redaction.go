package redaction

import (
	"encoding/json"
	"strings"
)

var sensitive = map[string]bool{"password": true, "token": true, "authorization": true, "password_hash": true, "secret": true}

func Map(input map[string]any) map[string]any {
	out := make(map[string]any, len(input))
	for k, v := range input {
		if sensitive[strings.ToLower(k)] {
			out[k] = "[REDACTED]"
			continue
		}
		switch x := v.(type) {
		case map[string]any:
			out[k] = Map(x)
		case []any:
			out[k] = slice(x)
		default:
			out[k] = v
		}
	}
	return out
}
func slice(input []any) []any {
	out := make([]any, len(input))
	for i, v := range input {
		if m, ok := v.(map[string]any); ok {
			out[i] = Map(m)
		} else {
			out[i] = v
		}
	}
	return out
}
func JSON(input []byte) ([]byte, error) {
	var v map[string]any
	if err := json.Unmarshal(input, &v); err != nil {
		return nil, err
	}
	return json.Marshal(Map(v))
}
