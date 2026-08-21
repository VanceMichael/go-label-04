package auth

import (
	"testing"
)

func TestPasswordHash(t *testing.T) {
	hash, err := HashPassword("correct horse")
	if err != nil {
		t.Fatal(err)
	}
	if Compare(hash, "correct horse") != nil {
		t.Fatal("password did not match")
	}
	if Compare(hash, "wrong") == nil {
		t.Fatal("wrong password matched")
	}
}
func TestTokenHashDeterministic(t *testing.T) {
	a := HashToken("abc")
	b := HashToken("abc")
	if a != b {
		t.Fatal("hash changed")
	}
	if a == HashToken("def") {
		t.Fatal("hash collision")
	}
}
