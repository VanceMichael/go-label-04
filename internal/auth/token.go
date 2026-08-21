package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
