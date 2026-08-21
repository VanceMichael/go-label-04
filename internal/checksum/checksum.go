package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"io"
)

func Bytes(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
func Reader(r io.Reader) (string, error) {
	if r == nil {
		return "", domain.ErrInvalid
	}
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
func Equal(a, b string) bool {
	return len(a) == len(b) && a != "" && Bytes([]byte(a)) == Bytes([]byte(b))
}
