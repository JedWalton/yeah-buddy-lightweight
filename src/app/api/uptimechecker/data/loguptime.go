package uptimechecker

// Uptime Log Management
func (r *Repository) LogUptime(endpointId int, statusCode int, responseTime int, isUp bool) error {
	query := `INSERT INTO UptimeLogs(endpoint_id, status_code, response_time, is_up, timestamp)
				VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.DB.Exec(query, endpointId, statusCode, responseTime, isUp)
	return err
}
