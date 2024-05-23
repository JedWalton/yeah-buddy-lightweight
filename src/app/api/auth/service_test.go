package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"i-couldve-got-six-reps/api/db"
	"os"
	"strings"
	"testing"
)

func TestCreateAndAuthenticateUser(t *testing.T) {
	// Setup Gin router for testing
	database := db.Init()
	authService := NewAuthService(database)

	username := "testCreateAndAuthenticateUser"
	password := "testCreateAndAuthenticateUser"

	// Clear the users table before each test
	_, err := database.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table before TestCreateAndAuthenticateUser: %v", err)
	}

	// Create a user
	_, err = authService.CreateUser(username, password)
	if err != nil {
		t.Fatalf("Could not create user: %v", err)
	}

	// Test authenticated user
	tokenString, err := authService.AuthenticateUser(username, password)
	if err != nil {
		t.Fatalf("Could not authenticate user: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		t.Fatalf("Error parsing token: %v", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Assuming the username is stored in the token's claims
		if claims["sub"] != username {
			t.Fatalf("Username does not match")
		}
	}

	// Clear the users table before each test
	_, err = database.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table after testing TestCreateAndAuthenticateUser: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	database := db.Init()
	authService := NewAuthService(database)

	username := "testDeleteUser"
	password := "testDeleteUserPassword"

	_, err := authService.CreateUser(username, password)
	if err != nil {
		t.Fatalf("Could not create user: %v", err)
	}

	// Verify the user was deleted form the db
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		t.Errorf("Could not query users table: %v", err)
	}
	if count != 1 {
		t.Errorf("User was not created in db. Got %d, want 1", count)
	}

	err = authService.DeleteUser(username)
	if err != nil {
		t.Errorf("Could not delete user: %v", err)
	}

	// Verify the user was deleted form the db
	count = 0 // Reset count
	err = database.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		t.Fatalf("Could not query users table: %v", err)
	}
	if count != 0 {
		t.Errorf("User was not deleted from the db. Got %d, want 0", count)
	}

	_, err = database.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table after TestDeleteUser test case: %v", err)
	}
}

// Testing generateJWT function
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
			token, err := generateJWT(tt.username)

			if (err != nil) != tt.wantErr {
				t.Errorf("generateJWT() error = %v, wantErr %v", err, tt.wantErr)
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
