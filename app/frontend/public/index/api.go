package index

import (
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	initPublic(r)
}

func initPublic(r *gin.Engine) {
	// Serve the index.html file for the root path
	indexhtml(r)
}

func indexhtml(r *gin.Engine) gin.IRoutes {
	return r.GET("/", func(c *gin.Context) {
		c.File("./frontend/public/index/index.html")
	})
}
