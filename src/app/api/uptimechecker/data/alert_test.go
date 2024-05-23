package uptimechecker

import (
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

// TestRecordAlert tests the insertion of an alert into the database.
func TestRecordAlert(t *testing.T) {
	db := db.Init()
	defer db.Close()

	repo := NewUptimeCheckerRepository(db)
	userrepo := auth.NewUserRepository(db)

	// Create a user to associate with the application
	userId, err := userrepo.CreateUser("testuser", "passwordHash")
	assert.NoError(t, err, "Failed to create user")

	// Create necessary entries for the test to satisfy foreign key constraints
	applicationId, err := repo.CreateApplication(userId, "TestRecordAlert",
		"TestRecordAlert test application")
	assert.NoError(t, err, "Failed to create application")

	channelId, err := repo.AddNotificationChannel(applicationId, "Email", "{\"email\":\"test@recordalert.com\"}")
	assert.NoError(t, err, "Failed to add notification channel")

	endpointId, err := repo.AddEndpoint(applicationId, "http://example.com", 30)
	assert.NoError(t, err, "Failed to add endpoint")

	// Now test recording an alert
	message := "Endpoint down"
	err = repo.RecordAlert(channelId, endpointId, message)
	assert.NoError(t, err, "Failed to record alert")

	// Verify the alert was inserted correctly
	var count int
	query := `SELECT COUNT(*) FROM Alerts WHERE channel_id = $1 AND endpoint_id = $2 AND message = $3`
	err = db.QueryRow(query, channelId, endpointId, message).Scan(&count)
	assert.NoError(t, err, "Failed to query alert")
	assert.Equal(t, 1, count, "Alert record not found")

	// Clean up
	userrepo.DeleteUserById(userId) // Assuming DeleteUser also cascades to delete applications and other dependent records
}

// Additional test implementations can be added here to cover all other repository functions.
