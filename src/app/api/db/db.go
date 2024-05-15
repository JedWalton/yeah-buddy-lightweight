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

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatal("Failed to open connection to postgres", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to open connection to postgres", err)
	}

	log.Printf("Connected to PostgreSQL")

	return db
}
