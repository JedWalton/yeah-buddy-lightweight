package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

// GenerateJWT generates a JWT token for a given user.
func GenerateJWT(username string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if jwtSecret == nil || len(jwtSecret) == 0 {
		log.Println("JWT_SECRET is not set or empty")
		return "", errors.New("JWT secret is not set or empty")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Failed to sign the token: %v\n", err)
		return "", err
	}

	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword), err
}
