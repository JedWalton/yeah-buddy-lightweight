package uptimechecker

import (
	"database/sql"
	uptimechecker "i-couldve-got-six-reps/api/uptimechecker/data"
	"log"

	"github.com/robfig/cron/v3"
)

type UptimeService struct {
	repo uptimechecker.Repository
	cron *cron.Cron // Add a cron field to manage the lifecycle of the cron scheduler
}

func NewUptimeService(database *sql.DB) *UptimeService {
	repo := uptimechecker.NewUptimeCheckerRepository(database)
	return &UptimeService{
		repo: repo,
		cron: cron.New(cron.WithSeconds()), // Initialize the cron with second precision
	}
}

func (s *UptimeService) StartUptimeService() {
	s.ScheduleEndpointChecks() // Start scheduled monitoring
	s.cron.Start()             // Start the cron scheduler
}

func (s *UptimeService) ScheduleEndpointChecks() {
	// Setup the periodic checks
	cronSchedule := "@every 30s" // Using a cron expression for every 30 seconds
	_, err := s.schedule(cronSchedule, s.CheckAllEndpoints)
	if err != nil {
		log.Printf("Failed to schedule endpoint checks: %v", err)
	}
}

func (s *UptimeService) schedule(spec string, job func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, job) // Schedule the job with the cron expression
}

func (s *UptimeService) Stop() {
	s.cron.Stop() // Stop the cron scheduler
}
