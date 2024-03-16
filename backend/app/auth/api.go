package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//func InitProtected(r *gin.Engine) {
//	protected := r.Group("/auth/protected")
//	protected.Use(middleware.AuthMiddleware())
//	login(protected)
//}

func initPublic(r *gin.Engine) {
	public := r.Group("/auth/public")
	login(public)
}

func Init(r *gin.Engine) {
	//InitProtected(r)
	initPublic(r)
}

func login(r *gin.RouterGroup) gin.IRoutes {
	return r.POST("/login", func(c *gin.Context) {
		loginHandler(c)
	})
}

func loginHandler(c *gin.Context) {
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
