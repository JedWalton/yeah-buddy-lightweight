package uptimechecker

import (
	"fmt"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"log"
	"net/http"
	"time"
)

// Endpoint Services
func (s *UptimeService) RegisterNewEndpoint(applicationId int, url string, monitoringInterval int) (int, error) {
	return s.repo.AddEndpoint(applicationId, url, monitoringInterval)
}

func (s *UptimeService) ModifyEndpoint(endpointId int, url string, monitoringInterval int, isActive bool) error {
	return s.repo.UpdateEndpoint(endpointId, url, monitoringInterval, isActive)
}

func (s *UptimeService) DeactivateEndpoint(endpointId int) error {
	return s.repo.UpdateEndpoint(endpointId, "", 0, false)
}

func (s *UptimeService) ActivateEndpoint(endpointId int) error {
	return s.repo.UpdateEndpoint(endpointId, "", 0, true)
}

func (s *UptimeService) CheckEndpointUptime(endpointId int) error {
	endpoint, err := s.repo.GetEndpoint(endpointId)
	if err != nil {
		return err
	}
	statusCode, responseTime, isUp := s.pingEndpoint(endpoint.URL)
	return s.repo.LogUptime(endpointId, statusCode, responseTime, isUp)
}

func (s *UptimeService) pingEndpoint(url string) (statusCode int, responseTime int, isUp bool) {
	// Start the timer to measure the response time
	start := time.Now()

	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		// If there's an error, we consider the endpoint down
		return 0, 0, false
	}
	defer resp.Body.Close()

	// Calculate the response time
	duration := time.Since(start)
	responseTime = int(duration.Milliseconds())

	// Check if the HTTP status code is in the range of 200-299
	isUp = resp.StatusCode >= 200 && resp.StatusCode <= 299
	statusCode = resp.StatusCode

	return statusCode, responseTime, isUp
}

func (s *UptimeService) CheckAllEndpoints() {
	endpoints, err := s.repo.ListActiveEndpoints()
	if err != nil {
		log.Printf("Error retrieving endpoints: %v", err)
		return
	}
	for _, endpoint := range endpoints {
		go func(ep types.Endpoint) {
			err := s.CheckEndpointUptime(ep.EndpointID)
			if err != nil {
				log.Printf("Uptime check failed for endpoint %d: %v", ep.EndpointID, err)
				s.handleDowntime(ep)
			}
		}(endpoint)
	}
}

func (s *UptimeService) handleDowntime(endpoint types.Endpoint) {
	// Logic to handle when an endpoint is detected as down
	channels, err := s.repo.ListNotificationChannels(endpoint.ApplicationID)
	if err != nil {
		log.Printf("Error retrieving notification channels for application %d: %v", endpoint.ApplicationID, err)
		return
	}
	for _, channel := range channels {
		message := fmt.Sprintf("Downtime alert for %s", endpoint.URL)
		err := s.SendAlert(channel.ChannelID, endpoint.EndpointID, message)
		if err != nil {
			log.Printf("Failed to send alert for endpoint %d: %v", endpoint.EndpointID, err)
		}
	}
}

func (s *UptimeService) ListAllActiveEndpoints() ([]types.Endpoint, error) {
	// This would be a repository method to get all active endpoints
	return s.repo.ListActiveEndpoints()
}
