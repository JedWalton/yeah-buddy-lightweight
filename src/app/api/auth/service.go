package auth

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo Repository
}

func NewAuthService(database *sql.DB) *AuthService {
	userRepo := NewUserRepository(database)
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) AuthenticateUser(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		log.Printf("Failed to get user by username: %v\n", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("Failed to compare password hashes: %v\n", err)
		return "", err
	}

	tokenString, err := generateJWT(user.Username)
	if err != nil {
		log.Printf("Failed to generate JWT: %v\n", err)
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) CreateUser(username, password string) (int, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return 0, err
	}

	userId, err := s.userRepo.CreateUser(username, passwordHash)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *AuthService) DeleteUser(username string) error {
	return s.userRepo.DeleteUser(username)
}

// generateJWT generates a JWT token for a given user.
func generateJWT(username string) (string, error) {
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
