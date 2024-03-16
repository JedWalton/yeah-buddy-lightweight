package auth

import (
	"bytes"
	"encoding/json"
	"i-couldve-got-six-reps/app/db"
	"i-couldve-got-six-reps/app/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginHandler(t *testing.T) {
	// Setup Gin router for testing
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Apply any necessary middleware
	r.Use(middleware.DB(db.Init()))

	// Initialize routes
	Init(r) // This initializes your routes, adjust as necessary

	// Create a request to pass to our handler.
	var jsonStr = []byte(`{"username":"admin", "password":"password"}`)
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
}
