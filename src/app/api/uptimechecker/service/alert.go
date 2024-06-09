package uptimechecker

type alertManager interface {
	SendAlert(channelId, endpointId int, message string) error
}

func (s *UptimeService) SendAlert(channelId, endpointId int, message string) error {
	return s.repo.RecordAlert(channelId, endpointId, message)
}

// Alert user of downtime.
