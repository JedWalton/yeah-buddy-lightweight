package auth

import (
	"database/sql"
	_ "github.com/lib/pq"
	"i-couldve-got-six-reps/app/auth/dtos"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) GetUserByUsername(username string) (*dtos.User, error) {
	var user dtos.User
	err := repo.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
