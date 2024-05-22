package uptimechecker

func (s *UptimeService) SendAlert(channelId, endpointId int, message string) error {
	return s.repo.RecordAlert(channelId, endpointId, message)
}
