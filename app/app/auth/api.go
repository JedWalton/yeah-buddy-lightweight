package auth

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"i-couldve-got-six-reps/app/auth/middleware"
	"net/http"
)

func Init(r *gin.Engine) {
	initPublic(r)
	initProtected(r)
}

func initPublic(r *gin.Engine) {
	public := r.Group("/auth/public")
	login(public)
	createUser(public)
}

func initProtected(r *gin.Engine) {
	protected := r.Group("/auth/protected")
	protected.Use(middleware.AuthMiddleware())
	getAccountInfo(protected)
}

func login(r *gin.RouterGroup) gin.IRoutes {
	return r.POST("/login", func(c *gin.Context) {
		loginHandler(c)
	})
}

func createUser(r *gin.RouterGroup) gin.IRoutes {
	return r.POST("/create", func(c *gin.Context) {
		createUserHandler(c)
	})
}

func getAccountInfo(r *gin.RouterGroup) gin.IRoutes {
	return r.GET("/account-info", func(c *gin.Context) {
		getAccountInfoHandler(c)
	})
}

func loginHandler(c *gin.Context) {
	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}
	userRepo := NewUserRepository(db)

	var loginCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := userRepo.GetUserByUsername(loginCredentials.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginCredentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	tokenString, err := GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func createUserHandler(c *gin.Context) {
	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}
	userRepo := NewUserRepository(db)

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	passwordHash, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := userRepo.CreateUser(user.Username, passwordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func getAccountInfoHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username not found in context"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username})
}
