package uptimechecker

//func TestUptimeService_ArchiveDay(t *testing.T) {
//	// Initialize database and repositories
//	database := db.Init()
//	defer database.Close()
//
//	userRepo := auth.NewUserRepository(database)
//
//	/* Ensure Database is in cleanstate */
//	_ = userRepo.DeleteUser("TestUptimeService_ArchiveDay User")
//	/* End of Ensure Database is in cleanstate */
//
//	userId, _ := userRepo.CreateUser("TestUptimeService_ArchiveDay User", "passwordHash")
//	repo := uptimechecker.NewUptimeCheckerRepository(database)
//	applicationId, err := repo.CreateApplication(userId, "TestUptimeService_ArchiveDay",
//		"TestUptimeService_ArchiveDay test application")
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//	url := "http://TestAddEndpoint.com"
//	monitoringInterval := 30
//
//	// Create multiple endpoints
//	endpointId1, _ := repo.AddEndpoint(applicationId, url, monitoringInterval)
//	endpointId2, _ := repo.AddEndpoint(applicationId, url, monitoringInterval)
//
//	// Generate and log uptime for multiple days and multiple endpoints
//	baseTime := time.Now().AddDate(0, 0, -3) // Go back three days
//	days := []int{-2, -1, 0}                 // logs for three days
//
//	for _, day := range days {
//		for _, endpointId := range []int{endpointId1, endpointId2} {
//			logTime := baseTime.AddDate(0, 0, day)
//			log := types.UptimeLog{
//				EndpointID:   endpointId,
//				StatusCode:   200,
//				ResponseTime: 120,
//				IsUp:         day != -2, // Make one day have all down logs
//				Timestamp:    logTime,
//			}
//			err := repo.LogUptime(log)
//			if err != nil {
//				t.Errorf("Failed to log uptime: %v", err)
//			}
//			log2 := types.UptimeLog{
//				EndpointID:   endpointId,
//				StatusCode:   200,
//				ResponseTime: 129,
//				IsUp:         day != -2, // Make one day have all down logs
//				Timestamp:    logTime,
//			}
//			err2 := repo.LogUptime(log2)
//			if err2 != nil {
//				t.Errorf("Failed to log uptime: %v", err2)
//			}
//			log3 := types.UptimeLog{
//				EndpointID:   endpointId,
//				StatusCode:   200,
//				ResponseTime: 125,
//				IsUp:         day != -2, // Make one day have all down logs
//				Timestamp:    logTime,
//			}
//			err3 := repo.LogUptime(log3)
//			if err3 != nil {
//				t.Errorf("Failed to log uptime: %v", err3)
//			}
//		}
//	}
//
//	s := NewUptimeService(database)
//	err = s.ArchiveDay()
//	if err != nil {
//		t.Errorf("Failed to archive uptime: %v", err)
//	}
//
//}
//
//func TestCalculateUptimePercentageForThisDay(t *testing.T) {
//	// table driven tests
//	tests := []struct {
//		name string
//		logs []types.UptimeLog
//		want float64
//	}{
//		{
//			name: "All Uptime",
//			logs: []types.UptimeLog{
//				{IsUp: true},
//				{IsUp: true},
//				{IsUp: true},
//			},
//			want: 100.0,
//		},
//		{
//			name: "All Downtime",
//			logs: []types.UptimeLog{
//				{IsUp: false},
//				{IsUp: false},
//				{IsUp: false},
//			},
//			want: 0.0,
//		},
//		{
//			name: "Evenly Split",
//			logs: []types.UptimeLog{
//				{IsUp: true},
//				{IsUp: false},
//			},
//			want: 50.0,
//		},
//		{
//			name: "No Logs",
//			logs: []types.UptimeLog{},
//			want: 0.0,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := calculateUptimePercentageForThisDay(tt.logs)
//			if got != tt.want {
//				t.Errorf("calculateUptimePercentageForThisDay() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//// ArchiveUptimePercentageForThisDay archives the uptime percentage for a specific day and endpointID
//func TestArchiveUptimePercentageForThisDay(t *testing.T) {
//	// Initialize database and repositories
//	database := db.Init()
//	defer database.Close()
//
//	userRepo := auth.NewUserRepository(database)
//
//	/* Ensure Database is in cleanstate */
//	_ = userRepo.DeleteUser("TestArchiveUptimePercentageForThisDay User")
//	/* End of Ensure Database is in cleanstate */
//
//	userId, _ := userRepo.CreateUser("TestArchiveUptimePercentageForThisDay User",
//		"passwordHash")
//	repo := NewUptimeCheckerRepository(database)
//	applicationId, err := repo.CreateApplication(
//		userId,
//		"TestArchiveUptimePercentageForThisDay app",
//		"test application")
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//	url := "http://TestAddEndpoint.com"
//	monitoringInterval := 30
//
//	// Create multiple endpoints
//	endpointId, _ := repo.AddEndpoint(applicationId, url, monitoringInterval)
//
//	// Generate and log uptime for multiple days and multiple endpoints
//	_ = repo.LogUptime(types.UptimeLog{
//		EndpointID:   endpointId,
//		StatusCode:   200,
//		ResponseTime: 121,
//		IsUp:         true,
//		Timestamp:    time.Now(),
//	})
//	_ = repo.LogUptime(types.UptimeLog{
//		EndpointID:   endpointId,
//		StatusCode:   200,
//		ResponseTime: 124,
//		IsUp:         true,
//		Timestamp:    time.Now(),
//	})
//	_ = repo.LogUptime(types.UptimeLog{
//		EndpointID:   endpointId,
//		StatusCode:   200,
//		ResponseTime: 130,
//		IsUp:         false,
//		Timestamp:    time.Now(),
//	})
//
//	// Test GetAllUptimeLogsForAGivenDayByEndpointIDAndDate
//	logs, err := repo.GetAllUptimeLogsForAGivenDayByEndpointIDAndDate(endpointId, time.Now())
//	if err != nil {
//		t.Errorf("Failed to get uptime logs: %v", err)
//	}
//
//	uptimePercentageForThisDay := calculateUptimePercentageForThisDay(logs)
//
//
//	/* Ensure Database is in cleanstate */
//	_ = userRepo.DeleteUser("TestArchiveUptimePercentageForThisDay User")
//	/* End of Ensure Database is in cleanstate */
