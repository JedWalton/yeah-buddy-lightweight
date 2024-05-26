package uptimechecker

import "i-couldve-got-six-reps/api/uptimechecker/types"

type PingEndpointResponseModel struct {
	responseTime int
	isUp         bool
}

// Uptime Log Management
func (r *Repository) LogUptime(log types.UptimeLog) error {
	query := `INSERT INTO UptimeLogs(endpoint_id, status_code, response_time, is_up, timestamp)
				VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.DB.Exec(query, log.EndpointID, log.StatusCode, log.ResponseTime, log.IsUp)
	return err
}
