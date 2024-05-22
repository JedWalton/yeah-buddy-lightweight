package uptimechecker

// Notification Services
func (s *UptimeService) AddNotificationChannel(applicationId int, typeStr string, details string) (int, error) {
	return s.repo.AddNotificationChannel(applicationId, typeStr, details)
}

func (s *UptimeService) UpdateNotificationChannel(channelId int, details string, isActive bool) error {
	return s.repo.UpdateNotificationChannel(channelId, details, isActive)
}

func (s *UptimeService) RemoveNotificationChannel(channelId int) error {
	return s.repo.DeleteNotificationChannel(channelId)
}
