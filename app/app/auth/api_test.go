package auth

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"i-couldve-got-six-reps/app/db"
	"i-couldve-got-six-reps/app/db/middleware"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
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
	_, err := repo.DB.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table before TestLoginHandler: %v", err)
	}

	/* Test login user not exist */
	// Create a request to pass to our handler.
	var jsonStr = []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password))
	req, err := http.NewRequest("POST", "/auth/public/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	repo.CreateUser(username, passwordHash)

	/* Test login wrong password */
	// Create a request to pass to our handler.
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", "wrongPassword")
	req, err = http.NewRequest("POST", "/auth/public/login", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Record the response
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	/* Test login success */
	// Create a request to pass to our handler.
	// For "application/x-www-form-urlencoded"
	formData = url.Values{}
	formData.Set("username", username)
	formData.Set("password", password)
	req, err = http.NewRequest("POST", "/auth/public/login", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Record the response
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the "Set-Cookie" header contains "auth_token"
	setCookie := rr.Header().Get("Set-Cookie")
	if !strings.Contains(setCookie, "auth_token") {
		t.Error("auth_token cookie not set in response")
	}

	// Optionally, you might want to check if the cookie's value (the token) is valid,
	// which could involve decoding the JWT and checking its payload, but that might
	// be beyond the scope of this particular test and could require additional setup.

	// Clear the users table after the test
	_, err = repo.DB.Exec("DELETE FROM users WHERE username = $1", username)
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

	// Create a request, this time without setting an Authorization header
	req, _ := http.NewRequest("GET", "/auth/protected/account-info", nil)

	// Instead, set the token as a cookie
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: getValidToken(),
	})

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

	// Create a request, this time also without an Authorization header but expecting it to fail
	req, _ := http.NewRequest("GET", "/auth/protected/account-info", nil)

	// Set an invalid token as a cookie
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: "SomeInvalidToken",
	})

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
