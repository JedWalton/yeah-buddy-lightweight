package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Init() *sql.DB {
	// Use environment variables
	postgresURL := os.Getenv("POSTGRESQL_URL")
	log.Printf("Connecting to PostgreSQL")

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to PostgreSQL")

	return db
}
