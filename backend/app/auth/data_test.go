package auth

import (
	"i-couldve-got-six-reps/app/db"
	"testing"
)

func TestGetUserByUsername(t *testing.T) {
	database := db.Init()

	repo := NewUserRepository(database)

	user, err := repo.GetUserByUsername("admin")
	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if user.Username != "admin" {
		t.Fatalf("Username does not match: got %v want %v", user.Username, "admin")
	}
	// Add more test assertions as needed
}
