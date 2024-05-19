package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(r *gin.Engine) {
	// Serve the main page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hello, World! Path: /",
		})
	})

	// HTMX endpoint to handle the dynamic content loading from an external HTML file
	r.GET("/hello", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", nil)
	})
}
