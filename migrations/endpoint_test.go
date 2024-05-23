package migrations

import (
	"i-couldve-got-six-reps/api/db"
	uptimechecker "i-couldve-got-six-reps/api/uptimechecker/data"
	"testing"
)

func TestAddEndpoint(t *testing.T) {
	// Set up the test database
	db, teardown := db.Init()
	defer teardown()

	// Initialize the repository
	repo := uptimechecker.NewUptimeCheckerRepository(db)

	// Test data
	applicationId := 1
	url := "http://example.com"
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

	err = db.QueryRow("SELECT application_id, url, monitoring_interval FROM Endpoints WHERE endpoint_id = $1", endpointId).Scan(&actualApplicationId, &actualUrl, &actualMonitoringInterval)
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
}
