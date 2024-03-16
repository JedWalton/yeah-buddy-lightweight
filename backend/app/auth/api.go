package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(r *gin.Engine) {
	// init auth service
	login(r)
}

//func hello_auth(r *gin.Engine) gin.IRoutes {
//	return r.GET("/auth/ping", func(c *gin.Context) {
//		login(c)
//	})
//}

func login(r *gin.Engine) gin.IRoutes {

	fmt.Printf("login\n")
	return r.POST("/login", func(c *gin.Context) {
		fmt.Printf("login\n")
		LoginHandler(c)
	})
}

func LoginHandler(c *gin.Context) {
	var loginCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Authenticate the user. This example uses hardcoded credentials.
	// Replace this with your actual authentication logic (e.g., database query).
	if loginCredentials.Username != "admin" || loginCredentials.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	tokenString, err := GenerateJWT(loginCredentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Public route

// Protected route
//r.GET("/protected", JWTAuthMiddleware(), func(c *gin.Context) {
//	// If the request reaches this point, it means the user is authenticated
//	username := c.MustGet("username").(string) // Extract username from the token
//	c.JSON(http.StatusOK, gin.H{"username": username, "message": "Welcome to the protected route!"})
//})
