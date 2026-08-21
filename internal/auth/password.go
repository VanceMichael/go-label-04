package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(raw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
}
func Compare(hash []byte, raw string) error { return bcrypt.CompareHashAndPassword(hash, []byte(raw)) }
