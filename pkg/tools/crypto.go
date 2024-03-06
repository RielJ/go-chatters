package tools

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword function
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Println(hash)
	return string(hash), nil
}

// CheckPasswordHash function
func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
