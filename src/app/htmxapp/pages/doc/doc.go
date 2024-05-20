package doc

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// Static file routes
	// Serve CSS, JS, and images for the 'doc' page
	r.Static("/doc/css", "./htmxapp/pages/doc/css")
	r.Static("/doc/images", "./htmxapp/pages/doc/images")
	r.Static("/doc/js", "./htmxapp/pages/doc/js")

	// Serve the main page
	r.GET("/docs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "doc.html", nil)
	})

	r.GET("/doc/htmxspecificlogic", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", nil)
	})
}
