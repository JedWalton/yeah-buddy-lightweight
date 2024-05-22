package uptimechecker

import (
	"fmt"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"log"
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
	// This function should perform an actual HTTP GET request to the URL
	// and measure the response time and status code. Stubbed for example.
	start := time.Now()
	// Simulate a request
	time.Sleep(time.Millisecond * 100) // Simulated delay
	duration := time.Since(start)

	return 200, int(duration.Milliseconds()), true
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
