package uptimechecker

import (
	"database/sql"
	"i-couldve-got-six-reps/api/auth"
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

func (s *UptimeService) Stop() {
	s.cron.Stop() // Stop the cron scheduler
}

func (s *UptimeService) StartUptimeService() {
	s.scheduleEndpointChecksAndDbEntry() // Start scheduled monitoring
	s.cron.Start()                       // Start the cron scheduler
}

func (s *UptimeService) StartUptimeServiceDev() {
	if scheduleEndpointChecksDev(s) {
		return
	}
	s.cron.Start() // Start the cron scheduler
}

func (s *UptimeService) scheduleEndpointChecksAndDbEntry() {
	// Setup the periodic checks
	cronSchedule := "@every 10m" // Using a cron expression for every 30 seconds
	_, err := s.cron.AddFunc(cronSchedule, s.CheckAllEndpoints)
	if err != nil {
		log.Printf("Failed to schedule endpoint checks: %v", err)
	}
}

func scheduleEndpointChecksDev(s *UptimeService) bool {
	db := s.repo.DB
	userRepo := auth.NewUserRepository(db)
	err := userRepo.DeleteUser("1) Dev User One")
	if err != nil {
		return true
	}
	err = userRepo.DeleteUser("2) Dev User Two")
	if err != nil {
		return true
	}
	userIdOne, _ := userRepo.CreateUser("1) Dev User One", "passwordHash")
	userIdTwo, _ := userRepo.CreateUser("2) Dev User Two", "passwordHash")
	isActive := true
	url := "https://lobster-app-dliao.ondigitalocean.app/"
	monitoringInterval := 10
	applicationId, err := s.CreateNewApplication(userIdOne, "1) Dev Application", "Dev application")
	applicationId2, err := s.CreateNewApplication(userIdTwo, "2) Dev Application 2", "Dev application 2")
	if err != nil {
		return true
	} // Create a new application
	endpointId, err := s.RegisterNewEndpoint(applicationId, "https://lobster-app-dliao.ondigitalocean.app/", monitoringInterval)
	endpointId2, err := s.RegisterNewEndpoint(applicationId2, "https://lobster-app-dliao.ondigitalocean.app/", monitoringInterval)
	if err != nil {
		return true
	} // Register a new endpoint
	err = s.ActivateEndpoint(endpointId, url, monitoringInterval, isActive)
	err = s.ActivateEndpoint(endpointId2, url, monitoringInterval, isActive)
	s.scheduleEndpointChecksAndDbEntry() // Start scheduled monitoring
	return false
}
