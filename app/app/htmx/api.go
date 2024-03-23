package htmx

import (
	"github.com/gin-gonic/gin"
	"i-couldve-got-six-reps/app/auth/middleware"
	"net/http"
)

func Init(r *gin.Engine) {
	initPublic(r)
	initProtected(r)
}

func initPublic(r *gin.Engine) {
	public := r.Group("/htmx/public")
	submitForm(public)
}

func initProtected(r *gin.Engine) {
	r.Group("/htmx/protected").Use(middleware.AuthMiddleware())
}

func submitForm(r *gin.RouterGroup) gin.IRoutes {
	return r.GET("/htmx", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte("Hello, World from HTMX"))
	})
}
