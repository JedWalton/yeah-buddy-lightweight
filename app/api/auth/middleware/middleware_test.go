package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "somerandomstring")
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name          string
		token         string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Valid Token",
			token:         getValidToken(),
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Invalid Token",
			token:         "invalid_token",
			expectedCode:  http.StatusUnauthorized,
			expectedError: "Unauthorized",
		},
		{
			name:          "Missing Token",
			token:         "",
			expectedCode:  http.StatusUnauthorized,
			expectedError: "Unauthorized",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/some-end-point", nil)

			// Instead of setting the Authorization header, set the token as a cookie
			if tc.token != "" {
				req.AddCookie(&http.Cookie{
					Name:  "auth_token",
					Value: tc.token,
				})
			}

			rr := httptest.NewRecorder()
			r := gin.Default()
			r.Use(AuthMiddleware())
			r.GET("/some-end-point", func(c *gin.Context) {
				username, exists := c.Get("username")
				if exists {
					c.String(http.StatusOK, fmt.Sprintf("%v", username))
				} else {
					c.String(http.StatusOK, "No username found")
				}
			})

			r.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("Expected response code %d, but got %d for case '%s'", tc.expectedCode, rr.Code, tc.name)
			}

			if tc.expectedError != "" && !strings.Contains(rr.Body.String(), tc.expectedError) {
				t.Errorf("Expected error '%s', but got '%s' for case '%s'", tc.expectedError, rr.Body.String(), tc.name)
			}
		})
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
