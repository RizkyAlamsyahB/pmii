package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims struktur untuk JWT payload
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte
var jwtExpirationHours int

// InitJWT inisialisasi JWT secret key dan expiration
// Harus dipanggil saat aplikasi start dengan config
func InitJWT(secret string, expirationHours int) {
	jwtSecret = []byte(secret)
	jwtExpirationHours = expirationHours
}

// GenerateJWT membuat JWT token dengan user ID dan role
// Token valid selama 24 jam
func GenerateJWT(userID int, role string) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret not initialized")
	}

	// Set expiration time sesuai config (default 24 jam)
	expirationTime := time.Now().Add(time.Duration(jwtExpirationHours) * time.Hour)

	// Buat claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT memvalidasi JWT token dan return claims jika valid
func ValidateJWT(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
