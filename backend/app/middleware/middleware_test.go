package middleware

import (
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
			req.Header.Add("Authorization", tc.token)
			rr := httptest.NewRecorder()

			r := gin.Default()
			r.Use(AuthMiddleware())
			r.GET("/some-end-point", func(c *gin.Context) {
				username, _ := c.Get("username")
				c.String(http.StatusOK, username.(string))
			})

			r.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("Expected response code %d, but got %d", tc.expectedCode, rr.Code)
			}

			if !strings.Contains(rr.Body.String(), tc.expectedError) && rr.Code == http.StatusUnauthorized {
				t.Errorf("Expected error '%s', but got %s", tc.expectedError, rr.Body.String())
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
