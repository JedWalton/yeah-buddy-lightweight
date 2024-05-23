package auth

import (
	"i-couldve-got-six-reps/api/auth/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine, authService *AuthService) {
	initPublic(r, authService)
	initProtected(r)
}

func initPublic(r *gin.Engine, authService *AuthService) {
	public := r.Group("/api/auth")
	login(public, authService)
	createUser(public, authService)
}

func initProtected(r *gin.Engine) {
	protected := r.Group("/auth/protected")
	protected.Use(middleware.AuthMiddleware())
}

func login(r *gin.RouterGroup, authService *AuthService) gin.IRoutes {
	return r.POST("/login", func(c *gin.Context) {
		loginHandler(c, authService)
	})
}

func createUser(r *gin.RouterGroup, authService *AuthService) gin.IRoutes {
	return r.POST("/create", func(c *gin.Context) {
		createUserHandler(c, authService)
	})
}

func loginHandler(c *gin.Context, authService *AuthService) {
	var loginCredentials struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
	}

	if err := c.ShouldBind(&loginCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tokenString, err := authService.AuthenticateUser(loginCredentials.Username, loginCredentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	// Set the JWT as a cookie in the response
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func createUserHandler(c *gin.Context, authService *AuthService) {
	var user struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
	}

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if _, err := authService.CreateUser(user.Username, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}
