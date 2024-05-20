package htmxapp

import (
	"i-couldve-got-six-reps/htmxapp/pages/doc"
	"i-couldve-got-six-reps/htmxapp/pages/home"

	"github.com/gin-gonic/gin"
)

func InitHtmxApp(r *gin.Engine) {
	r.LoadHTMLGlob("htmxapp/pages/**/*.html")
	home.Init(r)
	doc.Init(r)
}
