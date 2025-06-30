package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the provided password using the bcrypt library. It can't handle passwords longer than 72
// characters, so make sure it is below that limit
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash checks if a plain password matches a hashed password
func CheckPasswordHash(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
