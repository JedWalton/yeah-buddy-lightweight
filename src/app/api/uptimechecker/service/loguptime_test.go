package uptimechecker

import (
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	uptimechecker "i-couldve-got-six-reps/api/uptimechecker/data"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArchiveDay(t *testing.T) {
	// Initialize database and repositories
	database := db.Init()
	defer database.Close()

	userRepo := auth.NewUserRepository(database)

	// Ensure the database is in a clean state
	_ = userRepo.DeleteUser("TestArchiveDay User")

	// Create user and application
	userId, err := userRepo.CreateUser("TestArchiveDay User", "passwordHash")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	repo := uptimechecker.NewUptimeCheckerRepository(database)
	applicationId, err := repo.CreateApplication(
		userId,
		"TestArchiveDay",
		"Uptime Archive Day Test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Create multiple endpoints
	endpointId1, err := repo.AddEndpoint(applicationId, "http://TestArchiveDayEndpoint1.com", 30)
	if err != nil {
		t.Fatalf("Failed to add endpoint: %v", err)
	}
	endpointId2, err := repo.AddEndpoint(applicationId, "http://TestArchiveDayEndpoint2.com", 30)
	if err != nil {
		t.Fatalf("Failed to add second endpoint: %v", err)
	}

	// Generate and log uptime for specific day (two days ago)
	logDate := time.Now().AddDate(0, 0, -2) // Two days ago
	for _, endpointID := range []int{endpointId1, endpointId2} {
		for i := 0; i < 10; i++ {
			log := types.UptimeLog{
				EndpointID:   endpointID,
				StatusCode:   200,
				ResponseTime: 100,
				IsUp:         i%2 == 0,
				Timestamp:    logDate,
			}
			err := repo.LogUptime(log)
			if err != nil {
				t.Errorf("Failed to log uptime: %v", err)
			}
		}
	}

	// Create an instance of UptimeService
	service := UptimeService{repo: repo}

	// Execute ArchiveDay
	err = service.ArchiveDay()
	if err != nil {
		t.Fatalf("ArchiveDay failed: %v", err)
	}

	// Verify the results
	for _, endpointID := range []int{endpointId1, endpointId2} {
		// Check if the uptime logs have been pruned
		var count int
		err = database.QueryRow("SELECT COUNT(*) FROM UptimeLogs WHERE endpoint_id = $1 "+
			"AND DATE(timestamp) = $2", endpointID, logDate.Format("2006-01-02")).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to verify pruned logs: %v", err)
		}
		assert.Equal(t, 0, count, "Expected no logs for the pruned date, but some were found")

		// Check if uptime percentages have been archived
		var uptimePercentage float64
		err = database.QueryRow("SELECT uptime_percentage FROM UptimeLogsDailyArchive WHERE endpoint_id = $1 "+
			"AND date = $2", endpointID, logDate.Format("2006-01-02")).Scan(&uptimePercentage)
		if err != nil {
			t.Fatalf("Failed to retrieve archived uptime percentage: %v", err)
		}
		assert.InEpsilon(t, 50.0, uptimePercentage, 0.1,
			"Expected uptime percentage to be around 50%")
	}

	// Cleanup
	_ = userRepo.DeleteUser("TestArchiveDay User")
	if err != nil {
		t.Logf("Failed to clean up user after tests: %v", err)
	}
}

func TestCalculateUptimePercentageForThisDay(t *testing.T) {
	// table driven tests
	tests := []struct {
		name string
		logs []types.UptimeLog
		want float64
	}{
		{
			name: "All Uptime",
			logs: []types.UptimeLog{
				{IsUp: true},
				{IsUp: true},
				{IsUp: true},
			},
			want: 100.0,
		},
		{
			name: "All Downtime",
			logs: []types.UptimeLog{
				{IsUp: false},
				{IsUp: false},
				{IsUp: false},
			},
			want: 0.0,
		},
		{
			name: "Evenly Split",
			logs: []types.UptimeLog{
				{IsUp: true},
				{IsUp: false},
			},
			want: 50.0,
		},
		{
			name: "No Logs",
			logs: []types.UptimeLog{},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateUptimePercentageForThisDay(tt.logs)
			if got != tt.want {
				t.Errorf("calculateUptimePercentageForThisDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
