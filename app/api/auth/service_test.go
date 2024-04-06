package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"testing"
)

// Testing GenerateJWT function
func TestGenerateJWT(t *testing.T) {
	// Mocking environment variable
	os.Setenv("JWT_SECRET", "test_secret_key")

	// Table driven tests
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "All fields valid",
			username: "John Doe",
			wantErr:  false,
		},
		{
			name:     "Empty username",
			username: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.username)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil && strings.Contains(err.Error(), "token contains an invalid number of segments") {
				t.Errorf("Token couldn't be parsed properly, error: %v", err)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	var tests = []struct {
		name          string
		password      string
		expectedError error
	}{
		{
			name:     "empty_password",
			password: "",
		},
		{
			name:     "short_password",
			password: "123",
		},
		{
			name:     "average_password",
			password: "qwerty",
		},
		{
			name:     "long_password",
			password: "ThisIsAVeryLongPasswordWithSpecialCharacters!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := hashPassword(tt.password)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(tt.password))
			if err != nil {
				t.Fatalf("Fail to compare password and hashed password: %v", err)
			}
		})
	}
}
