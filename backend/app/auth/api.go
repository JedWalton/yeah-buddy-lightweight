package auth

import (
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// init auth service
	hello_auth(r)
}

func hello_auth(r *gin.Engine) gin.IRoutes {
	return r.GET("/ping", func(c *gin.Context) {
		login(c)
	})
}
