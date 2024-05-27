package uptimechecker

import (
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"log"
	"time"
)

func ArchiveDay(s *UptimeService, date time.Time) error {
	endpoints, err := s.repo.ListActiveEndpoints()
	if err != nil {
		return err // Properly return the error if it occurs
	}

	for _, endpoint := range endpoints {
		logs, err := s.repo.GetAllUptimeLogsForAGivenDayByEndpointIDAndDate(endpoint.EndpointID, date)
		if err != nil {
			log.Printf("Error retrieving logs for endpoint %d: %v", endpoint.EndpointID, err)
			continue // Log the error and continue processing other endpoints
		}
		uptimePercentageForThisDay := CalculateUptimePercentageForThisDay(logs)
		// Archive the uptime percentage for this day
		err = ArchiveUptimePercentageForThisDay(s, endpoint.EndpointID, uptimePercentageForThisDay)
		if err != nil {
			log.Printf("Error archiving uptime percentage for endpoint %d: %v", endpoint.EndpointID, err)
			continue // Log the error and continue processing other endpoints
		}
	}
	return nil // Return nil to indicate successful execution
}

func ArchiveUptimePercentageForThisDay(s *UptimeService, endpointID int, uptimePercentage float64) error {
	date := time.Now().AddDate(0, 0, -2).Format("1998-08-13")
	// Archive the uptime percentage for this day
	err := s.repo.ArchiveUptimePercentageForThisDay(endpointID, uptimePercentage, date)
	if err != nil {
		log.Printf("Error archiving uptime percentage for endpoint %d: %v", endpointID, err)
		return err
	}
	err = s.repo.PruneUptimeLogsByEndpointIDAndDate(endpointID, date)
	if err != nil {
		log.Printf("Error removing uptime logs for endpoint %d: %v", endpointID, err)
		return err
	}
	return nil
}

func CalculateUptimePercentageForThisDay(logs []types.UptimeLog) float64 {
	uptime := 0
	downtime := 0
	for _, log := range logs {
		if log.IsUp {
			uptime++
		} else {
			downtime++
		}
	}
	total := uptime + downtime
	if total == 0 {
		return 0.0 // Avoid division by zero if there are no logs
	}
	return float64(uptime) / float64(total) * 100.0
	// Archive the uptime percentage for this day
}
