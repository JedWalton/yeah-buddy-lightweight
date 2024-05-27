package uptimechecker

import (
	"fmt"
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"time"
)

// Uptime Log Management
func (r *Repository) LogUptime(log types.UptimeLog) error {
	query := `INSERT INTO UptimeLogs(endpoint_id, status_code, response_time, is_up, timestamp)
				VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.DB.Exec(query, log.EndpointID, log.StatusCode, log.ResponseTime, log.IsUp)
	return err
}

func (r *Repository) GetAllUptimeLogsForAGivenDayByEndpointIDAndDate(endpointID int, date time.Time) ([]types.UptimeLog, error) {
	var logs []types.UptimeLog
	dateFormatted := date.Format("2006-01-02") // Correct date format for YYYY-MM-DD
	query := `SELECT endpoint_id, status_code, response_time, is_up, timestamp
				FROM UptimeLogs
				WHERE endpoint_id = $1 AND DATE(timestamp) = $2`

	rows, err := r.DB.Query(query, endpointID, dateFormatted)
	if err != nil {
		return nil, fmt.Errorf("error querying logs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var log types.UptimeLog
		if err := rows.Scan(&log.EndpointID, &log.StatusCode, &log.ResponseTime, &log.IsUp, &log.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// ArchiveUptimePercentageForThisDay archives the uptime percentage for a specific day and endpointID
func (r *Repository) ArchiveUptimePercentageForThisDay(endpointID int, uptimePercentage float64, date string) error {
	query := `INSERT INTO UptimeLogsDailyArchive(endpoint_id, uptime_percentage, timestamp)
				VALUES ($1, $2, $3)`

	_, err := r.DB.Exec(query, endpointID, uptimePercentage, date)
	if err != nil {
		return fmt.Errorf("error archiving uptime percentage: %w", err)
	}
	return nil
}

func (r *Repository) PruneUptimeLogsByEndpointIDAndDate(endpointID int, date string) error {
	query := `DELETE FROM UptimeLogs WHERE endpoint_id = $1 AND DATE(timestamp) = $2`

	// Execute the query with endpointID and date as parameters
	if _, err := r.DB.Exec(query, endpointID, date); err != nil {
		return fmt.Errorf("error deleting logs: %w", err)
	}
	return nil
}
