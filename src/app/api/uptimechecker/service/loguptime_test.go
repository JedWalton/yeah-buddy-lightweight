package uptimechecker

import (
	"database/sql"
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	uptimechecker "i-couldve-got-six-reps/api/uptimechecker/data"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"testing"
	"time"
)

func TestArchiveFunctions(t *testing.T) {
	// Set up the test database
	database := db.Init()
	defer database.Close()

	// Set up repositories and service
	userRepo := auth.NewUserRepository(database)
	_ = userRepo.DeleteUser("TestArchiveFunctions User") // Clean up any existing test data
	userId, err := userRepo.CreateUser("TestArchiveFunctions User", "passwordHash")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := uptimechecker.NewUptimeCheckerRepository(database)
	s := NewUptimeService(database)

	applicationId, err := repo.CreateApplication(userId, "TestArchiveFunctions", "TestArchiveFunctions test application")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	url := "http://TestAddEndpoint.com"
	monitoringInterval := 30

	// Create multiple endpoints
	endpointId1, err := repo.AddEndpoint(applicationId, url, monitoringInterval)
	if err != nil {
		t.Fatalf("Failed to add endpoint: %v", err)
	}
	endpointId2, err := repo.AddEndpoint(applicationId, url, monitoringInterval)
	if err != nil {
		t.Fatalf("Failed to add endpoint: %v", err)
	}

	// Generate and log uptime for multiple days and multiple endpoints
	baseTime := time.Now().AddDate(0, 0, -3) // Three days ago
	days := []int{-2, -1, 0}                 // Logs for three days

	for _, day := range days {
		for _, endpointId := range []int{endpointId1, endpointId2} {
			logTime := baseTime.AddDate(0, 0, day)
			statusCode := 200
			if day == -2 {
				statusCode = 500 // Simulate downtime
			}
			log := types.UptimeLog{
				EndpointID:   endpointId,
				StatusCode:   statusCode,
				ResponseTime: 120,
				IsUp:         statusCode == 200,
				Timestamp:    logTime,
			}
			err := repo.LogUptime(log)
			if err != nil {
				t.Errorf("Failed to log uptime: %v", err)
			}
		}
	}

	// Test archiving for yesterday
	testDate := baseTime.AddDate(0, 0, -1) // Yesterday
	err = s.ArchiveDay(testDate)
	if err != nil {
		t.Errorf("Failed to archive day: %v", err)
	}

	// Verify if archiving was successful and that no data from other days was affected
	verifyArchivingResults(t, database, endpointId1, endpointId2, testDate)
}

func verifyArchivingResults(t *testing.T, db *sql.DB, endpointId1 int, endpointId2 int, date time.Time) {
	// Check that the archived data for the specified date is present and correct
	var count int
	var avgUptime float64
	dateString := date.Format("2006-01-02")

	// Query to verify correct archival of uptime percentage
	query := `SELECT COUNT(*), COALESCE(AVG(uptime_percentage), 0) FROM UptimeLogsDailyArchive 
          WHERE endpoint_id IN ($1, $2) AND date = $3`
	err := db.QueryRow(query, endpointId1, endpointId2, dateString).Scan(&count, &avgUptime)
	if err != nil {
		t.Errorf("Failed to query archived uptime percentages: %v", err)
	}
	if count == 0 {
		t.Errorf("No archival records found for the given date: %s", dateString)
	} else if avgUptime == 0 {
		t.Errorf("Archived uptime percentage is unexpectedly zero.")
	}
	// You may also want to check if the average uptime matches expected values, depending on the input logs.

	// Verify that logs for the specified date are pruned
	query = `SELECT COUNT(*) FROM UptimeLogs WHERE endpoint_id IN ($1, $2) AND DATE(timestamp) = $3`
	err = db.QueryRow(query, endpointId1, endpointId2, dateString).Scan(&count)
	if err != nil {
		t.Errorf("Failed to query UptimeLogs for pruning verification: %v", err)
	}
	if count != 0 {
		t.Errorf("Logs for date %s were not pruned correctly, found %d logs", dateString, count)
	}

	// Optional: Verify that logs for other dates are intact (not pruned)
	query = `SELECT COUNT(*) FROM UptimeLogs WHERE endpoint_id IN ($1, $2) AND DATE(timestamp) != $3`
	err = db.QueryRow(query, endpointId1, endpointId2, dateString).Scan(&count)
	if err != nil {
		t.Errorf("Failed to verify integrity of other dates' logs: %v", err)
	}
	if count == 0 {
		t.Errorf("Logs for other dates appear to have been pruned erroneously")
	}
}
