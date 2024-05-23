package uptimechecker

import (
	"github.com/stretchr/testify/assert"
	"i-couldve-got-six-reps/api/auth"
	"i-couldve-got-six-reps/api/db"
	"testing"
)

func TestCreateApplication(t *testing.T) {
	db := db.Init()
	defer db.Close()

	userRepo := auth.NewUserRepository(db)
	userId, _ := userRepo.CreateUser("1) Test Create Application User One", "passwordHash")

	repo := NewUptimeCheckerRepository(db)
	appName := "Test App"
	appDesc := "A test application for monitoring."

	appID, err := repo.CreateApplication(userId, appName, appDesc)
	assert.NoError(t, err)
	assert.NotZero(t, appID)

	userRepo.DeleteUserById(userId)
}

func TestUpdateAndDeleteApplication(t *testing.T) {
	db := db.Init()
	defer db.Close()

	userRepo := auth.NewUserRepository(db)
	userId, _ := userRepo.CreateUser("1) Test Update And Delete Application User One", "passwordHash")

	repo := NewUptimeCheckerRepository(db)
	appName := "Test App"
	appDesc := "A test application for monitoring."

	// Create application to update
	appID, err := repo.CreateApplication(userId, appName, appDesc)
	assert.NoError(t, err)
	assert.NotZero(t, appID)

	// Update application
	newName := "Updated Test App"
	newDesc := "Updated description for the test application."
	err = repo.UpdateApplication(appID, newName, newDesc)
	assert.NoError(t, err)

	// Verify update
	app, err := repo.FindApplication(appID)
	assert.NoError(t, err)
	assert.Equal(t, newName, app.Name)
	assert.Equal(t, newDesc, app.Description)

	// Delete application
	err = repo.DeleteApplication(appID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.FindApplication(appID)
	assert.Error(t, err) // Should error because the application no longer exists

	userRepo.DeleteUserById(userId)
}
