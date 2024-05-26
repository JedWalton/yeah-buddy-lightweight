package uptimechecker

import (
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"log"
)

func ArchiveToday(s *UptimeService) {
	endpoints, err := s.repo.ListActiveEndpoints()
	if err != nil {
		return
	}

	for _, endpoint := range endpoints {
		logs, err := s.repo.GetAllUptimeLogsForADayPerEndpointID(endpoint.EndpointID)
		if err != nil {
			continue
		}
		uptimePercentageForThisDay := CalculateUptimePercentageForThisDay(logs)
		// Archive the uptime percentage for this day
		err = ArchiveUptimePercentageForThisDay(s, endpoint.EndpointID, uptimePercentageForThisDay)
		if err != nil {
			continue
		}
	}
}

func ArchiveUptimePercentageForThisDay(s *UptimeService, endpointID int, uptimePercentage float64) error {
	// Archive the uptime percentage for this day
	err := s.repo.ArchiveUptimePercentageForThisDay(endpointID, uptimePercentage)
	if err != nil {
		log.Printf("Error archiving uptime percentage for endpoint %d: %v", endpointID, err)
		return err
	}
	err = s.repo.RemoveUptimeLogsForToday(endpointID)
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
