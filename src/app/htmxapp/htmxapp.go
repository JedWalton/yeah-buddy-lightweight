package htmxapp

import (
	"github.com/gin-gonic/gin"
	"i-couldve-got-six-reps/htmxapp/home"
)

func InitHtmxApp(r *gin.Engine) {
	r.LoadHTMLGlob("htmxapp/**/*")
	home.Init(r)
}
