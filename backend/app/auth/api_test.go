package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"i-couldve-got-six-reps/app/db"
	"i-couldve-got-six-reps/app/db/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginHandler(t *testing.T) {
	// Setup Gin router for testing
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	database := db.Init()
	repo := NewUserRepository(database)

	// Apply any necessary middleware
	r.Use(middleware.DB(database))

	// Initialize routes
	Init(r) // This initializes your routes, adjust as necessary

	username := "testUserLoginHandler"
	password := "testPasswordLoginHandler"
	passwordHash, _ := hashPassword(password)

	// Clear the users table before each test
	_, err := repo.DB.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table before TestLoginHandler: %v", err)
	}

	repo.CreateUser(username, passwordHash)

	// Create a request to pass to our handler.
	var jsonStr = []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password))
	req, err := http.NewRequest("POST", "/auth/public/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response body to check for the token
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("could not decode response body", err)
	}

	// Check if the token is present and not empty
	token, ok := response["token"]
	if !ok || token == "" {
		t.Errorf("handler returned unexpected body, token not present or empty")
	}

	// Clear the users table before each test
	_, err = repo.DB.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table after TestLoginHandler: %v", err)
	}
}

func TestCreateUserHandler(t *testing.T) {
	// Setup Gin router for testing
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	database := db.Init()
	repo := NewUserRepository(database)

	// Apply any necessary middleware
	r.Use(middleware.DB(database))

	// Initialize routes
	Init(r) // This initializes your routes, adjust as necessary

	username := "testUserCreateUserHandler"
	password := "testPasswordCreateUserHandler"
	passwordHash, _ := hashPassword(password)

	// Clear the users table before each test
	_, err := repo.DB.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table before TestCreateUserHandler: %v", err)
	}

	// Create a request to pass to our handler.
	var jsonStr = []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, passwordHash))
	req, err := http.NewRequest("POST", "/auth/public/create", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.String() != "{\"message\":\"User created\"}" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "User created")
	}

	// Clear the users table after each test
	_, err = repo.DB.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table after TestCreateUserHandler: %v", err)
	}
}

func TestGetAccountInfoHandlerAuthorized(t *testing.T) {
	// Set up Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	Init(r)

	// Set up protected route with AuthMiddleware

	// Create a request with Authorization header
	req, _ := http.NewRequest("GET", "/auth/protected/account-info", nil)
	req.Header.Set("Authorization", getValidToken())

	// Record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert on the response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

}

func TestGetAccountInfoHandlerUnauthorized(t *testing.T) {
	// Set up Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	Init(r)

	// Set up protected route with AuthMiddleware

	// Create a request with Authorization header
	req, _ := http.NewRequest("GET", "/auth/protected/account-info", nil)
	req.Header.Set("Authorization", "SomeBullshitInvalidToken")

	// Record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert on the response
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}

}

// Get a valid token for testing.
func getValidToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString
}
