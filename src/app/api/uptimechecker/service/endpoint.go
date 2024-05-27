package uptimechecker

import (
	"fmt"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"log"
	"net/http"
	"time"
)

// CRUD
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

func (s *UptimeService) ListAllActiveEndpoints() ([]types.Endpoint, error) {
	return s.repo.ListActiveEndpoints()
}

// END OF CRUD

// Endpoint related operations
func (s *UptimeService) CheckAllEndpoints() {
	endpoints, err := s.repo.ListActiveEndpoints()
	if err != nil {
		log.Printf("Error retrieving endpoints: %v", err)
		return
	}

	for _, endpoint := range endpoints {
		go func(ep types.Endpoint) {
			err := checkEndpointUptime(ep.EndpointID, s)
			if err != nil {
				log.Printf("Uptime check failed for endpoint %d: %v", ep.EndpointID, err)
				s.handleDowntime(ep)
			}
		}(endpoint)
	}
}

func checkEndpointUptime(endpointId int, s *UptimeService) error {
	ticker := time.NewTicker(30 * time.Second) // Const
	defer ticker.Stop()

	log.Printf("Starting to check every 30s for endpoint %d", endpointId)

	allResponses, err := getAllResponsesPerUnitTimeDatabaseEntry(endpointId, ticker, s)
	if err != nil {
		return err
	}

	// Now we can safely log and calculate after the goroutine is done
	log.Printf("All endpoint responses per unit time check: %v", allResponses)
	if len(allResponses) == 0 {
		return fmt.Errorf("no responses collected")
	}

	uptimeLog := calculatePingEndpointResponseDatabaseEntry(allResponses)

	return s.repo.LogUptime(uptimeLog)
}

func getAllResponsesPerUnitTimeDatabaseEntry(endpointId int, ticker *time.Ticker, s *UptimeService) ([]types.UptimeLog, error) {
	quit := make(chan struct{})
	done := make(chan error)
	var allEndpointResponsesPerUnitTimeCheck []types.UptimeLog

	go func() {
		defer close(done)
		for {
			select {
			case <-ticker.C:
				if len(allEndpointResponsesPerUnitTimeCheck) < (s.timeMinutesBetweenDbEntries * 2) { // Const
					pingEndpointResp, err := pingEndpointById(endpointId, s)
					if err != nil {
						log.Printf("Error pinging endpoint %d: %v", endpointId, err)
						done <- err
						return
					}
					allEndpointResponsesPerUnitTimeCheck = append(allEndpointResponsesPerUnitTimeCheck, pingEndpointResp)
				} else {
					done <- nil
					return
				}
			case <-quit:
				return
			}
		}
	}()

	// Wait for the goroutine to finish
	if err := <-done; err != nil {
		return nil, err
	}
	return allEndpointResponsesPerUnitTimeCheck, nil
}

func calculatePingEndpointResponseDatabaseEntry(allEndpointResponsesPer10Min []types.UptimeLog) types.UptimeLog {

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
	uptimeLog := pingEndpoint(endpointId, endpoint.URL)
	log.Printf("Endpoint %d is up: %v, status code: %d, response time: %dms", endpointId,
		uptimeLog.IsUp, uptimeLog.StatusCode, uptimeLog.ResponseTime)
	return uptimeLog, nil
}

func pingEndpoint(endpointId int, url string) types.UptimeLog {
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

	return types.UptimeLog{
		EndpointID:   endpointId,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
		IsUp:         isUp,
		Timestamp:    time.Now()}
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
