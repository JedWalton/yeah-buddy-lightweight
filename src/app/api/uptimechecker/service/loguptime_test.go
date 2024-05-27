package uptimechecker

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	uptimechecker "i-couldve-got-six-reps/api/uptimechecker/data"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"testing"
	"time"
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
	logDate := time.Now().AddDate(0, 0, -2)   // Two days ago
	otherDate := time.Now().AddDate(0, 0, -1) // One day ago (should not be pruned or archived)
	for _, endpointID := range []int{endpointId1, endpointId2} {
		for i := 0; i < 10; i++ {
			log1 := types.UptimeLog{
				EndpointID:   endpointID,
				StatusCode:   200,
				ResponseTime: 100,
				IsUp:         i%2 == 0,
				Timestamp:    logDate,
			}
			log2 := types.UptimeLog{
				EndpointID:   endpointID,
				StatusCode:   200,
				ResponseTime: 100,
				IsUp:         i%2 == 0,
				Timestamp:    otherDate,
			}
			// Logs for the logDate
			err := repo.LogUptime(log1)
			if err != nil {
				t.Errorf("Failed to log uptime for logDate: %v", err)
			}
			// Logs for the otherDate
			err = repo.LogUptime(log2)
			if err != nil {
				t.Errorf("Failed to log uptime for otherDate: %v", err)
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
		// Check if the logs for the logDate are pruned and archived correctly
		verifyArchiveAndPrune(t, database, endpointID, logDate, 50.0)

		// Ensure that logs from the otherDate are not pruned or archived
		verifyNotAffected(t, database, endpointID, otherDate)
	}

	// Cleanup
	_ = userRepo.DeleteUser("TestArchiveDay User")
}

func verifyArchiveAndPrune(t *testing.T, db *sql.DB, endpointID int, date time.Time, expectedPerc float64) {
	var count int
	var uptimePercentage float64
	// Verify pruned logs
	err := db.QueryRow("SELECT COUNT(*) FROM UptimeLogs WHERE endpoint_id = $1 AND DATE(timestamp) = $2", endpointID, date.Format("2006-01-02")).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to verify pruned logs: %v", err)
	}
	assert.Equal(t, 0, count, "Expected no logs for the pruned date, but some were found")
	// Verify archived uptime
	err = db.QueryRow("SELECT uptime_percentage FROM UptimeLogsDailyArchive WHERE endpoint_id = $1 AND date = $2", endpointID, date.Format("2006-01-02")).Scan(&uptimePercentage)
	if err != nil {
		t.Fatalf("Failed to retrieve archived uptime percentage: %v", err)
	}
	assert.InEpsilon(t, expectedPerc, uptimePercentage, 0.1, "Expected uptime percentage to be around the expected value")
}

func verifyNotAffected(t *testing.T, db *sql.DB, endpointID int, date time.Time) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM UptimeLogs WHERE endpoint_id = $1 AND DATE(timestamp) = $2", endpointID, date.Format("2006-01-02")).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to verify non-pruned logs: %v", err)
	}
	assert.NotEqual(t, 0, count, "Logs for a day not targeted should not be pruned")

	var uptimePercentage float64
	err = db.QueryRow("SELECT uptime_percentage FROM UptimeLogsDailyArchive WHERE endpoint_id = $1 AND date = $2", endpointID, date.Format("2006-01-02")).Scan(&uptimePercentage)
	assert.Error(t, err, "No archive should exist for non-targeted days")
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
