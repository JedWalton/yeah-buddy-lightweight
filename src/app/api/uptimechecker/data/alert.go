package uptimechecker

// Alert Log Management
func (r *Repository) RecordAlert(channelId, endpointId int, message string) error {
	query := `INSERT INTO Alerts (channel_id, endpoint_id, message, sent_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.DB.Exec(query, channelId, endpointId, message)
	return err
}
