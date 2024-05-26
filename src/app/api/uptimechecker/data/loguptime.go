package uptimechecker

import "i-couldve-got-six-reps/api/uptimechecker/types"

// Uptime Log Management
func (r *Repository) LogUptime(log types.UptimeLog) error {
	query := `INSERT INTO UptimeLogsToday(endpoint_id, status_code, response_time, is_up, timestamp)
				VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.DB.Exec(query, log.EndpointID, log.StatusCode, log.ResponseTime, log.IsUp)
	return err
}

func (r *Repository) GetAllUptimeLogsForADayPerEndpointID(endpointID int) ([]types.UptimeLog, error) {
	var logs []types.UptimeLog
	query := `SELECT endpoint_id, status_code, response_time, is_up, timestamp
				FROM UptimeLogsLast24Hours WHERE endpoint_id = $1`
	rows, err := r.DB.Query(query, endpointID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log types.UptimeLog
		err := rows.Scan(&log.EndpointID, &log.StatusCode, &log.ResponseTime, &log.IsUp, &log.Timestamp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *Repository) ArchiveUptimePercentageForThisDay(endpointID int, uptimePercentage float64) error {
	query := `INSERT INTO UptimeLogsDailyArchive(endpoint_id, uptime_percentage, timestamp)
				VALUES ($1, $2, NOW())`
	_, err := r.DB.Exec(query, endpointID, uptimePercentage)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) RemoveUptimeLogsForToday(endpointID int) error {
	query := `DELETE FROM UptimeLogsLast24Hours WHERE endpoint_id = $1`
	_, err := r.DB.Exec(query, endpointID)
	if err != nil {
		return err
	}
	return nil
}
