package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a hash out of a password passed into it
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hash), nil
}

// CompareHash checks if a password and a hash are the same
func CompareHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
