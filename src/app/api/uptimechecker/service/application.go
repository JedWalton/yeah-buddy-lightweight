package uptimechecker

// Application Services
func (s *UptimeService) CreateNewApplication(userId int, name, description string) (int, error) {
	return s.repo.CreateApplication(userId, name, description)
}

func (s *UptimeService) UpdateExistingApplication(applicationId int, name, description string) error {
	return s.repo.UpdateApplication(applicationId, name, description)
}

func (s *UptimeService) RemoveApplication(applicationId int) error {
	return s.repo.DeleteApplication(applicationId)
}
