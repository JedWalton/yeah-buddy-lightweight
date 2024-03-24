package htmx

import (
	"i-couldve-got-six-reps/app/auth/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// This package should contain the API routes for the HTMX service.
// This will include page navigation and other HTMX specific routes.

func Init(r *gin.Engine) {
	initPublic(r)
	initProtected(r)
}

func initPublic(r *gin.Engine) {
	public := r.Group("/htmx/public")
	helloWorld(public)
}

func initProtected(r *gin.Engine) {
	protected := r.Group("/htmx/protected")
	protected.Use(middleware.AuthMiddleware())
	helloWorldButItsProtected(protected)
}

func helloWorld(r *gin.RouterGroup) gin.IRoutes {
	return r.GET("/htmx", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte("Hello, World from HTMX"))
	})
}

func helloWorldButItsProtected(r *gin.RouterGroup) gin.IRoutes {
	return r.GET("/htmx", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte("Hello, World from HTMX"))
	})
}
