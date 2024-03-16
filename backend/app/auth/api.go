package auth

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
