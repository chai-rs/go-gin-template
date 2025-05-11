package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// MustHashPassword hashes a password and panics if it fails.
func MustHashPassword(password string) string {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		panic(err)
	}

	return hashedPassword
}

// ComparePassword compares a password with a hashed password.
func ComparePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
