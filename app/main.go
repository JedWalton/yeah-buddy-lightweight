package main

import (
	"database/sql"
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	"i-couldve-got-six-reps/api/db/middleware"
	"i-couldve-got-six-reps/frontend/public/index"
	"os"

	"github.com/gin-gonic/gin"
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

	// init api and frontend
	initApi(r)
	initFrontend(r)

	port := os.Getenv("PORT") // Get the PORT environment variable
	if port == "" {
		port = "8080" // Default to 8080 if no PORT environment variable is set
	}

	err := r.Run(":" + port) // listen on the specified port
	if err != nil {
		panic(err) // added panic to handle potential error from r.Run
	}
}

func initGlobalMiddleware(r *gin.Engine, database *sql.DB) {
	r.Use(middleware.DB(database))
}

func initApi(r *gin.Engine) {
	auth.Init(r)
}

func initFrontend(r *gin.Engine) {
	index.Init(r)
}
