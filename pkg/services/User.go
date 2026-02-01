package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gofor-little/env"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
}

func (this *User) HashPassword(password string) (string, error) {

	pepper := env.Get("HASH_KEY", "default")
	h := hmac.New(sha256.New, []byte(pepper))
	h.Write([]byte(password))

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (this *User) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
