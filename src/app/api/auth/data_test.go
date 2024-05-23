package auth

import (
	"i-couldve-got-six-reps/api/db"
	"testing"
)

func TestGetUserByUsername(t *testing.T) {
	database := db.Init()

	repo := NewUserRepository(database)

	username := "admin"
	password := "password"
	passwordHash, err := hashPassword(password)

	// Clear the users table before each test
	_, err = repo.DB.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table: %v", err)
	}

	repo.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, passwordHash)

	user, err := repo.GetUserByUsername("admin")
	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if user.Username != "admin" {
		t.Fatalf("Username does not match: got %v want %v", user.Username, "admin")
	}
	// Add more test assertions as needed

	repo.DB.Exec("DELETE FROM users where username = $1", username)
}

func TestUserRepository_CreateUser(t *testing.T) {
	database := db.Init()
	repo := NewUserRepository(database)

	// Define a test user
	username := "testUserRepositoryCreateUser"
	passwordHash, _ := hashPassword("testUserRepositoryPasswordHashCreateUser")

	// Clear the users table before each test
	_, err := repo.DB.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table: %v", err)
	}

	// Test CreateUser function
	_, err = repo.CreateUser(username, passwordHash)
	if err != nil {
		t.Errorf("CreateUser failed: %v", err)
	}

	// Verify the user was added to the db
	var count int
	err = repo.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1 AND password_hash = $2", username, passwordHash).Scan(&count)
	if err != nil {
		t.Fatalf("Could not query users table: %v", err)
	}
	if count != 1 {
		t.Errorf("User was not inserted into the db. Got %d, want 1", count)
	}

	_, err = repo.DB.Exec("DELETE FROM users where username = $1 AND password_hash = $2", username, passwordHash)
	if err != nil {
		t.Fatalf("Could not clear users table after test case: %v", err)
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	database := db.Init()
	repo := NewUserRepository(database)

	// Define a test user
	username := "testUserRepositoryDeleteUser"
	passwordHash, _ := hashPassword("testUserRepositoryDeleteUserPassword")

	// Clear the users table before each test
	_, err := repo.DB.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		t.Fatalf("Could not clear users table: %v", err)
	}

	// Test CreateUser function
	_, err = repo.CreateUser(username, passwordHash)
	if err != nil {
		t.Errorf("CreateUser failed: %v", err)
	}

	err = repo.DeleteUser(username)

	// Verify the user was deleted form the db
	var count int
	err = repo.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1 AND password_hash = $2", username, passwordHash).Scan(&count)
	if err != nil {
		t.Fatalf("Could not query users table: %v", err)
	}
	if count != 0 {
		t.Errorf("User was not deleted from the db. Got %d, want 0", count)
	}

	_, err = repo.DB.Exec("DELETE FROM users where username = $1 AND password_hash = $2", username, passwordHash)
	if err != nil {
		t.Fatalf("Could not clear users table after TestUserRepository_DeleteUser test case: %v", err)
	}
}
