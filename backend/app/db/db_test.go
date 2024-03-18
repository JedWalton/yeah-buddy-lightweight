package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

// TestInit tries to connect to the database using the provided POSTGRESQL_URL environment variable.
func TestInit(t *testing.T) {
	// Load the environment variables from the .env file.
	Init()

	tests := []struct {
		name    string
		env     string
		wantErr bool
	}{
		{"valid postgres url", os.Getenv("POSTGRESQL_URL"), false},
		// Note: For an invalid test, ensure the URL is indeed invalid and points to a non-existent or inaccessible database.
		{"invalid postgres url", "postgres://invalid:invalid@localhost:5432/invalid?sslmode=disable", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily set the environment variable for the PostgreSQL URL.
			os.Setenv("POSTGRESQL_URL", tt.env)

			db, err := sql.Open("postgres", os.Getenv("POSTGRESQL_URL"))
			if err != nil {
				t.Fatalf("Could not open database connection: %v", err)
			}

			// Try to ping the database to check if the connection is valid.
			err = db.Ping()

			if tt.wantErr {
				assert.Error(t, err, "Expected an error for the database connection, but got none")
			} else {
				assert.NoError(t, err, "Did not expect an error for the database connection")
			}

			// Clean up: Close the database connection after the test.
			if db != nil {
				db.Close()
			}
		})
	}
}
