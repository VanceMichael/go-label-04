package pagination

import (
	"encoding/base64"
	"fmt"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"strconv"
	"strings"
)

type Cursor struct {
	CreatedAt int64
	ID        string
}

func Encode(c Cursor) string {
	raw := fmt.Sprintf("%d|%s", c.CreatedAt, c.ID)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}
func Decode(raw string) (Cursor, error) {
	if raw == "" {
		return Cursor{}, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return Cursor{}, domain.ErrInvalid
	}
	parts := strings.SplitN(string(data), "|", 2)
	if len(parts) != 2 {
		return Cursor{}, domain.ErrInvalid
	}
	ts, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || parts[1] == "" {
		return Cursor{}, domain.ErrInvalid
	}
	return Cursor{CreatedAt: ts, ID: parts[1]}, nil
}
func Next[T any](items []T, limit int, hasMore bool, encode func(T) Cursor) string {
	if !hasMore || len(items) == 0 {
		return ""
	}
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}
	return Encode(encode(items[len(items)-1]))
}
