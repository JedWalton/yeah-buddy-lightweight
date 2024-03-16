package auth

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

// GenerateJWT generates a JWT token for a given user.
func GenerateJWT(username string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
