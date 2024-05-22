package uptimechecker

import (
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"log"
)

// Notification Channel Management
func (r *Repository) AddNotificationChannel(applicationId int, typeStr string, details string) (int, error) {
	var channelId int
	query := `INSERT INTO types.NotificationChannels (application_id, type, details, is_active) VALUES ($1, $2, $3, true) RETURNING channel_id`
	err := r.DB.QueryRow(query, applicationId, typeStr, details).Scan(&channelId)
	if err != nil {
		return 0, err
	}
	return channelId, nil
}

func (r *Repository) UpdateNotificationChannel(channelId int, details string, isActive bool) error {
	query := `UPDATE NotificationChannels SET details = $1, is_active = $2, updated_at = NOW() WHERE channel_id = $3`
	_, err := r.DB.Exec(query, details, isActive, channelId)
	return err
}

func (r *Repository) DeleteNotificationChannel(channelId int) error {
	query := `DELETE FROM NotificationChannels WHERE channel_id = $1`
	_, err := r.DB.Exec(query, channelId)
	return err
}

// ListNotificationChannels retrieves all active notification channels for a given application.
func (r *Repository) ListNotificationChannels(applicationID int) ([]types.NotificationChannel, error) {
	var channels []types.NotificationChannel
	query := `SELECT channel_id, application_id, type, details, is_active FROM NotificationChannels WHERE application_id = $1 AND is_active = TRUE`
	rows, err := r.DB.Query(query, applicationID)
	if err != nil {
		log.Printf("Error querying notification channels: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ch types.NotificationChannel
		err := rows.Scan(&ch.ChannelID, &ch.ApplicationID, &ch.Type, &ch.Details, &ch.IsActive)
		if err != nil {
			log.Printf("Error scanning notification channel: %v", err)
			continue
		}
		channels = append(channels, ch)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return nil, err
	}

	return channels, nil
}
