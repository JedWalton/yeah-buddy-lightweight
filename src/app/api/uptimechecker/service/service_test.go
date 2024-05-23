package uptimechecker

import (
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

func TestUptimeServiceIntegration(t *testing.T) {
	// Initialize the test database connection
	db := db.Init()
	defer db.Close()

	// Create UptimeService instance
	service := NewUptimeService(db)

	// Start the uptime service in a development environment
	go service.StartUptimeServiceDev() // Run this as a goroutine to avoid blocking
	defer service.Stop()               // Ensure the cron is stopped after the test

	// Allow some time for the service to start and schedule jobs
	time.Sleep(2 * time.Second)

	// Perform validations
	// Check if applications and endpoints are correctly registered and activated
	// This should query the database and check for the existence of expected rows
	// This is a simplified example:
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Applications").Scan(&count)
	assert.NoError(t, err)
	assert.True(t, count >= 2, "Expected at least two applications to be registered")

	// Additional checks can include checking for registered endpoints, active statuses, etc.
	err = db.QueryRow("SELECT COUNT(*) FROM Endpoints WHERE is_active = true").Scan(&count)
	assert.NoError(t, err)
	assert.True(t, count >= 2, "Expected at least two endpoints to be active")

	// Cleanup
	authService := auth.NewAuthService(db)
	authService.DeleteUser("TestUptimeServiceIntegrationUser")
}
