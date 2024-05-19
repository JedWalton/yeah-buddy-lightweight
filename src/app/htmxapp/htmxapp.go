package htmxapp

import (
	"github.com/gin-gonic/gin"
	"i-couldve-got-six-reps/htmxapp/pages/home"
)

func InitHtmxApp(r *gin.Engine) {
	r.LoadHTMLGlob("htmxapp/pages/**/*")
	home.Init(r)
}
