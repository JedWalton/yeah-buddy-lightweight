package main

import (
	"database/sql"
	"i-couldve-got-six-reps/app/auth"
	"i-couldve-got-six-reps/app/db"
	"i-couldve-got-six-reps/app/db/middleware"
	"i-couldve-got-six-reps/app/htmx"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	r := gin.Default()

	// serve frontend
	r.Static("/app", "../frontend")

	// init database
	database := db.Init()
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			panic(err)
		}
	}(database)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(os.Stdout)

	// init middleware
	initGlobalMiddleware(r, database)

	// init services
	initService(r)

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

func initService(r *gin.Engine) {
	auth.Init(r)
	htmx.Init(r)
	//payment.Init(r)
	// init other services
	// ...
}
