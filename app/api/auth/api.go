package auth

import (
	"github.com/gin-gonic/gin"
	"i-couldve-got-six-reps/api/auth/middleware"
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
