package uptimechecker

import (
	"github.com/stretchr/testify/assert"
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"testing"
	"time"
)

func TestLogUptime(t *testing.T) {
	// Initialize database and repositories
	database := db.Init()
	defer database.Close()

	userRepo := auth.NewUserRepository(database)

	/* Ensure Database is in cleanstate */
	_ = userRepo.DeleteUser("TestUptimeChecker User")
	/* End of Ensure Database is in cleanstate */

	userId, _ := userRepo.CreateUser("TestUptimeChecker User", "passwordHash")
	repo := NewUptimeCheckerRepository(database)
	applicationId, err := repo.CreateApplication(
		userId,
		"TestUptimeChecker",
		"UptimeChecker test application")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	url := "http://TestAddEndpoint.com"
	monitoringInterval := 30

	// Create multiple endpoints
	endpointId, _ := repo.AddEndpoint(applicationId, url, monitoringInterval)

	// Generate and log uptime for multiple days and multiple endpoints
	log := types.UptimeLog{
		EndpointID:   endpointId,
		StatusCode:   200,
		ResponseTime: 120,
		IsUp:         true,
		Timestamp:    time.Now(),
	}
	err = repo.LogUptime(log)
	if err != nil {
		t.Errorf("Failed to log uptime: %v", err)
	}

	// Test GetAllUptimeLogsForAGivenDayByEndpointIDAndDate
	logs, err := repo.GetAllUptimeLogsForAGivenDayByEndpointIDAndDate(endpointId, time.Now())
	if err != nil {
		t.Errorf("Failed to get uptime logs: %v", err)
	}

	if !logs[0].IsUp {
		t.Errorf("Expected uptime log to be up, got down")
	}

	/* Ensure Database is in cleanstate */
	_ = userRepo.DeleteUser("TestUptimeChecker User")
	/* End of Ensure Database is in cleanstate */
}

// ArchiveUptimePercentageForThisDay archives the uptime percentage for a specific day and endpointID
func TestArchiveUptimePercentageForThisDay(t *testing.T) {
	// Initialize database and repositories
	database := db.Init()
	defer database.Close()

	userRepo := auth.NewUserRepository(database)

	/* Ensure Database is in cleanstate */
	_ = userRepo.DeleteUser("TestArchiveUptimePercentageForThisDay User")
	/* End of Ensure Database is in cleanstate */

	userId, _ := userRepo.CreateUser("TestArchiveUptimePercentageForThisDay User",
		"passwordHash")
	repo := NewUptimeCheckerRepository(database)
	applicationId, err := repo.CreateApplication(
		userId,
		"TestArchiveUptimePercentageForThisDay app",
		"test application")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	url := "http://TestAddEndpoint.com"
	monitoringInterval := 30

	// Create multiple endpoints
	endpointId, _ := repo.AddEndpoint(applicationId, url, monitoringInterval)

	_ = repo.ArchiveUptimePercentageForThisDay(endpointId, 100, time.Now())

	// Archive the uptime percentage.
	date := time.Now()
	err = repo.ArchiveUptimePercentageForThisDay(endpointId, 100, date)
	if err != nil {
		t.Fatalf("Failed to archive uptime percentage: %v", err)
	}

	// Verify the results.
	var (
		epID       int
		uptimePerc float64
		timestamp  time.Time
	)

	err = database.QueryRow("SELECT endpoint_id, uptime_percentage, date FROM UptimeLogsDailyArchive WHERE endpoint_id = $1 AND date = $2", endpointId, date.Format("2006-01-02")).Scan(&epID, &uptimePerc, &timestamp)
	if err != nil {
		t.Errorf("Failed to retrieve archived data: %v", err)
	}

	// Assuming 'date' is the correct column, and it stores just the date
	assert.Equal(t, endpointId, epID, "Endpoint ID should match")
	assert.Equal(t, 100.0, uptimePerc, "Uptime percentage should be 100")
	assert.True(t, date.Format("2006-01-02") == timestamp.Format("2006-01-02"), "Date should match the archived day")

	/* Ensure Database is in cleanstate */
	_ = userRepo.DeleteUser("TestArchiveUptimePercentageForThisDay User")
}

func TestPruneUptimeLogsByEndpointIDAndDate(t *testing.T) {
	// Initialize database and repositories
	database := db.Init()
	defer database.Close()

	userRepo := auth.NewUserRepository(database)

	// Ensure the database is in a clean state
	_ = userRepo.DeleteUser("TestPruneUptimeLogs User")

	// Create user and application
	userId, err := userRepo.CreateUser("TestPruneUptimeLogs User", "passwordHash")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	repo := NewUptimeCheckerRepository(database)
	applicationId, err := repo.CreateApplication(
		userId,
		"TestPruneUptimeLogs",
		"Uptime Logs Pruning Test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Create multiple endpoints
	endpointId1, err := repo.AddEndpoint(applicationId, "http://TestPruneEndpoint1.com", 30)
	if err != nil {
		t.Fatalf("Failed to add endpoint: %v", err)
	}
	endpointId2, err := repo.AddEndpoint(applicationId, "http://TestPruneEndpoint2.com", 30)
	if err != nil {
		t.Fatalf("Failed to add second endpoint: %v", err)
	}

	// Generate and log uptime for different days and endpoints
	logDate := time.Now().AddDate(0, 0, -1) // Logs for yesterday for pruning
	otherDate := time.Now()                 // Logs for today to ensure they aren't pruned

	logsToCreate := 5
	for i := 0; i < logsToCreate; i++ {
		// Logs for endpointId1 on logDate
		log1 := types.UptimeLog{
			EndpointID:   endpointId1,
			StatusCode:   200,
			ResponseTime: 100,
			IsUp:         true,
			Timestamp:    logDate,
		}
		err := repo.LogUptime(log1)
		if err != nil {
			t.Errorf("Failed to log uptime for endpoint 1: %v", err)
		}

		// Logs for endpointId2 on otherDate
		log2 := types.UptimeLog{
			EndpointID:   endpointId2,
			StatusCode:   200,
			ResponseTime: 100,
			IsUp:         true,
			Timestamp:    otherDate,
		}
		err = repo.LogUptime(log2)
		if err != nil {
			t.Errorf("Failed to log uptime for endpoint 2: %v", err)
		}
	}

	// Prune logs for endpointId1 on logDate
	err = repo.PruneUptimeLogsByEndpointIDAndDate(endpointId1, logDate)
	if err != nil {
		t.Errorf("Failed to prune uptime logs: %v", err)
	}

	// Verify the logs for endpointId1 on logDate are pruned
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM UptimeLogs WHERE endpoint_id = $1 AND DATE(timestamp) = $2", endpointId1, logDate.Format("2006-01-02")).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to verify pruned logs for endpoint 1: %v", err)
	}
	assert.Equal(t, 0, count, "Expected no logs for the pruned date for endpoint 1, but some were found")

	// Verify the logs for endpointId2 on otherDate are NOT pruned
	err = database.QueryRow("SELECT COUNT(*) FROM UptimeLogs WHERE endpoint_id = $1 AND DATE(timestamp) = $2", endpointId2, otherDate.Format("2006-01-02")).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to verify non-pruned logs for endpoint 2: %v", err)
	}
	assert.NotEqual(t, 0, count, "Logs for endpoint 2 on a different day should not be pruned")

	// Cleanup
	_ = userRepo.DeleteUser("TestPruneUptimeLogs User")
}
