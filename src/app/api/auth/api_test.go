package auth

import (
	"bytes"
	"fmt"
	"i-couldve-got-six-reps/api/db"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginHandler(t *testing.T) {
	// Setup Gin router for testing
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	database := db.Init()
	authService := NewAuthService(database)

	// Initialize routes
	Init(r, authService) // This initializes your routes, adjust as necessary

	username := "testUserLoginHandler"
	password := "testPasswordLoginHandler"
	passwordHash, _ := hashPassword(password)

	// Clear the users table before each test
	_, err := database.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table before TestLoginHandler: %v", err)
	}

	/* Test login user not exist */
	// Create a request to pass to our handler.
	var jsonStr = []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password))
	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonStr))
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

	_, err = authService.userRepo.CreateUser(username, passwordHash)
	if err != nil {
		return
	}

	/* Test login wrong password */
	// Create a request to pass to our handler.
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", "wrongPassword")
	req, err = http.NewRequest("POST", "/api/auth/login", strings.NewReader(formData.Encode()))
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
	req, err = http.NewRequest("POST", "/api/auth/login", strings.NewReader(formData.Encode()))
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
	_, err = database.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table after TestLoginHandler: %v", err)
	}
}

func TestCreateUserHandler(t *testing.T) {
	// Setup Gin router for testing
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	database := db.Init()
	authService := NewAuthService(database)

	// Initialize routes
	Init(r, authService) // This initializes your routes, adjust as necessary

	username := "testUserCreateUserHandler"
	password := "testPasswordCreateUserHandler"
	passwordHash, _ := hashPassword(password)

	// Clear the users table before each test
	_, err := database.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table before TestCreateUserHandler: %v", err)
	}

	// Create a request to pass to our handler.
	var jsonStr = []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, passwordHash))
	req, err := http.NewRequest("POST", "/api/auth/create", bytes.NewBuffer(jsonStr))
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
	_, err = database.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table after TestCreateUserHandler: %v", err)
	}
}

//
//// Get a valid token for testing.
//func getValidToken() string {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"username": "testuser",
//	})
//
//	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
//
//	return tokenString
//}
