package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//	func loginHandler(c *gin.Context) {
//		db, ok := c.MustGet("db").(*sql.DB)
//		if !ok {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
//			return
//		}
//		userRepo := NewUserRepository(db)
//
//		var loginCredentials struct {
//			Username string `form:"username" json:"username"`
//			Password string `form:"password" json:"password"`
//		}
//
//		if err := c.ShouldBind(&loginCredentials); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//			return
//		}
//
//		user, err := userRepo.GetUserByUsername(loginCredentials.Username)
//		if err != nil {
//			if err == sql.ErrNoRows {
//				c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
//			} else {
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
//			}
//			return
//		}
//
//		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginCredentials.Password)); err != nil {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
//			return
//		}
//
//		tokenString, err := generateJWT(user.Username)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
//			return
//		}
//
//		// Set the JWT as a cookie in the response
//		http.SetCookie(c.Writer, &http.Cookie{
//			Name:     "auth_token",            // Cookie name, e.g., "auth_token"
//			Value:    tokenString,             // The JWT
//			Path:     "/",                     // Cookie path. Using "/" means it's sent for all paths.
//			HttpOnly: true,                    // HttpOnly prevents JavaScript access to the cookie, enhancing security
//			Secure:   true,                    // Secure flag ensures the cookie is sent over HTTPS only, enhancing security
//			SameSite: http.SameSiteStrictMode, // SameSite=Strict prevents the cookie from being sent with cross-site requests
//			// Set the MaxAge or Expires field if you want the cookie to expire
//		})
//
//		// Optionally, if you want HTMX to replace part of your page with a response,
//		// you can return HTML instead of JSON
//		// For example, to update a "registration message" div:
//		//c.HTML(http.StatusOK, "registration_success.html", gin.H{
//		//	"Username": user.Username,
//		//})
//
//		// Optionally, redirect the user or send a success response
//		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
//	}
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

//func createUserHandler(c *gin.Context) {
//	db, ok := c.MustGet("db").(*sql.DB)
//	if !ok {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
//		return
//	}
//	userRepo := NewUserRepository(db)
//
//	var user struct {
//		Username string `form:"username" json:"username"`
//		Password string `form:"password" json:"password"`
//	}
//
//	if err := c.ShouldBind(&user); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//		return
//	}
//
//	passwordHash, err := hashPassword(user.Password)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
//		return
//	}
//
//	if err := userRepo.CreateUser(user.Username, passwordHash); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
//		return
//	}
//
//	// Optionally, if you want HTMX to replace part of your page with a response,
//	// you can return HTML instead of JSON
//	// For example, to update a "registration message" div:
//	//c.HTML(http.StatusOK, "registration_success.html", gin.H{
//	//	"Username": user.Username,
//	//})
//
//	// Or simply return a success message as JSON if that's your preference
//	c.JSON(http.StatusOK, gin.H{"message": "User created"})
//}

func createUserHandler(c *gin.Context, authService *AuthService) {
	var user struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
	}

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := authService.CreateUser(user.Username, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Response
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
