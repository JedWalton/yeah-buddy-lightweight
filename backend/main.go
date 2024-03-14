package main

import (
	"github.com/gin-gonic/gin"
	"i-couldve-got-six-reps/app/auth"
	"i-couldve-got-six-reps/app/db"
	"i-couldve-got-six-reps/app/middleware"
)

func main() {
	r := gin.Default()

	// init database
	db := db.Init()
	defer db.Close()

	r.Use(middleware.DB(db))

	// init services
	init_service(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func init_service(r *gin.Engine) {
	// init auth service
	auth.Init(r)
	// init other services
	// ...
}
