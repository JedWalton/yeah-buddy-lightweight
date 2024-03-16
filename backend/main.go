package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"i-couldve-got-six-reps/app/auth"
	"i-couldve-got-six-reps/app/db"
	"i-couldve-got-six-reps/app/middleware"
)

func main() {
	r := gin.Default()

	// init database
	database := db.Init()
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			panic(err)
		}
	}(database)

	// init middleware
	initGlobalMiddleware(r, database)

	// init services
	initService(r)

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}

func initGlobalMiddleware(r *gin.Engine, database *sql.DB) {
	r.Use(middleware.DB(database))
}

func initService(r *gin.Engine) {
	auth.Init(r)

	// init other services
	// ...
}
