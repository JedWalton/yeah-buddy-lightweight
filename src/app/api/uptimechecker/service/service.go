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

func (s *UptimeService) StartUptimeServiceDev() {
	isActive := true
	url := "https://lobster-app-dliao.ondigitalocean.app/"
	monitoringInterval := 10
	applicationId, err := s.CreateNewApplication("My Application", "This is my application")
	applicationId2, err := s.CreateNewApplication("My Application2", "This is my application2")
	if err != nil {
		return
	} // Create a new application
	endpointId, err := s.RegisterNewEndpoint(applicationId, "https://lobster-app-dliao.ondigitalocean.app/", 30)
	endpointId2, err := s.RegisterNewEndpoint(applicationId2, "https://lobster-app-dliao.ondigitalocean.app/", 30)
	if err != nil {
		return
	} // Register a new endpoint
	err = s.ActivateEndpoint(endpointId, url, monitoringInterval, isActive)
	err = s.ActivateEndpoint(endpointId2, url, monitoringInterval, isActive)
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
