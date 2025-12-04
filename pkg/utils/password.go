package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword melakukan hashing password menggunakan bcrypt
// Cost menggunakan bcrypt.DefaultCost (10) untuk keamanan optimal
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash memverifikasi apakah password cocok dengan hash
// Returns true jika password valid
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
