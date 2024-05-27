package uptimechecker

import (
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"testing"
	"time"
)

func TestUptimeChecker(t *testing.T) {
	// Initialize database and repositories
	database := db.Init()
	defer database.Close()

	userRepo := auth.NewUserRepository(database)

	/* Ensure Database is in cleanstate */
	_ = userRepo.DeleteUser("TestUptimeChecker User")
	/* End of Ensure Database is in cleanstate */

	userId, _ := userRepo.CreateUser("TestUptimeChecker User", "passwordHash")
	repo := NewUptimeCheckerRepository(database)
	applicationId, err := repo.CreateApplication(userId, "TestUptimeChecker", "UptimeChecker test application")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	url := "http://TestAddEndpoint.com"
	monitoringInterval := 30

	// Create multiple endpoints
	endpointId1, _ := repo.AddEndpoint(applicationId, url, monitoringInterval)
	endpointId2, _ := repo.AddEndpoint(applicationId, url, monitoringInterval)

	// Generate and log uptime for multiple days and multiple endpoints
	baseTime := time.Now().AddDate(0, 0, -3) // Go back three days
	days := []int{-2, -1, 0}                 // logs for three days

	for _, day := range days {
		for _, endpointId := range []int{endpointId1, endpointId2} {
			logTime := baseTime.AddDate(0, 0, day)
			log := types.UptimeLog{
				EndpointID:   endpointId,
				StatusCode:   200,
				ResponseTime: 120,
				IsUp:         day != -2, // Make one day have all down logs
				Timestamp:    logTime,
			}
			err := repo.LogUptime(log)
			if err != nil {
				t.Errorf("Failed to log uptime: %v", err)
			}
		}
	}

	// Test specific day pruning for endpointId1
	pruneDate := baseTime.AddDate(0, 0, -2).Format("2006-01-02")
	err = repo.PruneUptimeLogsByEndpointIDAndDate(endpointId1, pruneDate)
	if err != nil {
		t.Errorf("Failed to prune uptime logs: %v", err)
	}

	// Check logs for other endpoint on same day
	logs, err := repo.GetAllUptimeLogsForAGivenDayByEndpointIDAndDate(endpointId2, baseTime.AddDate(0, 0, -2))
	if err != nil || len(logs) == 0 {
		t.Errorf("Logs for other endpoints on same day were incorrectly pruned, error: %v", err)
	}

	// Check logs for other days for endpointId1
	logs, err = repo.GetAllUptimeLogsForAGivenDayByEndpointIDAndDate(endpointId1, baseTime.AddDate(0, 0, -1))
	if err != nil || len(logs) == 0 {
		t.Errorf("Logs for other days were incorrectly pruned, error: %v", err)
	}

	/* Ensure Database is in cleanstate */
	_ = userRepo.DeleteUser("TestUptimeChecker User")
	/* End of Ensure Database is in cleanstate */
}
