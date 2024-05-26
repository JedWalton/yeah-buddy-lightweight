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

func (s *UptimeService) ActivateEndpoint(endpointId int, url string, monitoringInterval int, isActive bool) error {
	return s.repo.UpdateEndpoint(endpointId, url, monitoringInterval, isActive)
}

func (s *UptimeService) CheckEndpointUptimeCalledEvery10Mins(endpointId int) error {
	// Add a cron here every 30s.
	// This will check the endpoint every 30s.
	// Store results into data structure.
	// After 10mins, do necessary calculations and store into db.
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	iterations := 0
	var allEndpointResponsesPer10Min []types.UptimeLog
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				if iterations < 20 {
					pingEndpointResp, err := pingEndpointById(endpointId, s)
					if err != nil {
						log.Printf("Error pinging endpoint %d: %v", endpointId, err)
						close(quit)
						return
					}
					allEndpointResponsesPer10Min = append(allEndpointResponsesPer10Min, pingEndpointResp)
				} else {
					close(quit)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Calculate the average response time and uptime percentage
	uptimeLog := calculatePingEndpointResponseDatabaseEntry(allEndpointResponsesPer10Min)

	return s.repo.LogUptime(uptimeLog)
}

func calculatePingEndpointResponseDatabaseEntry(
	allEndpointResponsesPer10Min []types.UptimeLog) types.UptimeLog {

	var count int
	for _, response := range allEndpointResponsesPer10Min {
		// Calculate the average response time and uptime percentage
		if !response.IsUp {
			return types.UptimeLog{allEndpointResponsesPer10Min[0].EndpointID,
				0, 0, false, time.Now()}
		}
		count++
	}

	var toAvgResponseTimeSum int
	for _, response := range allEndpointResponsesPer10Min {
		toAvgResponseTimeSum += response.ResponseTime
	}
	avgResponseTime := toAvgResponseTimeSum / count

	return types.UptimeLog{allEndpointResponsesPer10Min[0].EndpointID, 200,
		avgResponseTime, true, time.Now()}
}

func pingEndpointById(endpointId int, s *UptimeService) (types.UptimeLog, error) {
	endpoint, err := s.repo.GetEndpoint(endpointId)
	if err != nil {
		return types.UptimeLog{}, err
	}
	uptimeLog := s.pingEndpoint(endpointId, endpoint.URL)
	log.Printf("Endpoint %d is up: %v, status code: %d, response time: %dms", endpointId,
		uptimeLog.IsUp, uptimeLog.StatusCode, uptimeLog.ResponseTime)
	return uptimeLog, nil
}

func (s *UptimeService) pingEndpoint(endpointId int, url string) types.UptimeLog {
	// Start the timer to measure the response time
	start := time.Now()

	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		// If there's an error, we consider the endpoint down
		return types.UptimeLog{
			EndpointID:   endpointId,
			ResponseTime: 0,
			IsUp:         false,
			Timestamp:    time.Time{},
		}
	}
	defer resp.Body.Close()

	// Calculate the response time
	duration := time.Since(start)
	responseTime := int(duration.Milliseconds())

	// Check if the HTTP status code is in the range of 200-299
	isUp := resp.StatusCode >= 200 && resp.StatusCode <= 299
	statusCode := resp.StatusCode

	return types.UptimeLog{endpointId, statusCode, responseTime, isUp, time.Now()}
}

// CheckAllEndpoints Triggered every 10mins. This will check all active endpoints
// And trigger 30s checks for each endpoint.
// Only write to db after 10mins.
func (s *UptimeService) CheckAllEndpoints() {
	endpoints, err := s.repo.ListActiveEndpoints()
	if err != nil {
		log.Printf("Error retrieving endpoints: %v", err)
		return
	}

	for _, endpoint := range endpoints {
		go func(ep types.Endpoint) {
			err := s.CheckEndpointUptimeCalledEvery10Mins(ep.EndpointID)
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
