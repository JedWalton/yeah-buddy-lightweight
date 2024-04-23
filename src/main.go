package main

import (
	"database/sql"
	"i-couldve-got-six-reps/app/auth"
	"i-couldve-got-six-reps/app/db"
	"i-couldve-got-six-reps/htmx"
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

	// init app and frontend
	initApi(r, database)
	initHtmx(r)

	port := os.Getenv("PORT") // Get the PORT environment variable
	if port == "" {
		port = "8080" // Default to 8080 if no PORT environment variable is set
	}

	err := r.Run(":" + port) // listen on the specified port
	if err != nil {
		panic(err) // added panic to handle potential error from r.Run
	}
}

func initApi(r *gin.Engine, database *sql.DB) {
	authService := auth.NewAuthService(database)
	auth.Init(r, authService)

}

func initHtmx(r *gin.Engine) {
	htmx.Init(r)
}