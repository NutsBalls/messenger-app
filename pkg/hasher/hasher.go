package hasher

import (
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func HashToken(token string) (string, error) {
	h := sha256.New()
	h.Write([]byte(token))
	return Hash(string(h.Sum(nil)))
}
