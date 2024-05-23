package uptimechecker

import (
	"github.com/stretchr/testify/assert"
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"

	"testing"
)

func TestAddEndpoint(t *testing.T) {
	// Set up the test database
	database := db.Init()
	defer database.Close()

	userRepo := auth.NewUserRepository(database)
	userId, _ := userRepo.CreateUser("1) Test Create Application User One", "passwordHash")
	// Initialize the repository
	repo := NewUptimeCheckerRepository(database)

	// First, create necessary entries for the test to satisfy foreign key constraints
	applicationId, err := repo.CreateApplication(userId, "TestRecordAlert", "TestRecordAlert test application")
	assert.NoError(t, err, "Failed to create application")
	// Test data
	url := "http://TestAddEndpoint.com"
	monitoringInterval := 30

	// Run the method
	endpointId, err := repo.AddEndpoint(applicationId, url, monitoringInterval)

	// Check for errors
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the returned endpointId
	if endpointId == 0 {
		t.Fatalf("Expected a valid endpoint ID, got 0")
	}

	// Verify the endpoint was actually inserted into the database
	var actualApplicationId int
	var actualUrl string
	var actualMonitoringInterval int

	err = database.QueryRow("SELECT application_id, url, monitoring_interval FROM Endpoints WHERE endpoint_id = $1", endpointId).Scan(&actualApplicationId, &actualUrl, &actualMonitoringInterval)
	if err != nil {
		t.Fatalf("Expected no error when querying the inserted endpoint, got %v", err)
	}

	if actualApplicationId != applicationId {
		t.Errorf("Expected application ID to match, got %d", actualApplicationId)
	}
	if actualUrl != url {
		t.Errorf("Expected URL to match, got %s", actualUrl)
	}
	if actualMonitoringInterval != monitoringInterval {
		t.Errorf("Expected monitoring interval to match, got %d", actualMonitoringInterval)
	}

	userRepo.DeleteUserById(userId)
}
