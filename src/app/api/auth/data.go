package auth

import (
	"database/sql"
	_ "github.com/lib/pq"
	"i-couldve-got-six-reps/api/auth/dtos"
)

type Repository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) Repository {
	return Repository{DB: db}
}

func (repo *Repository) GetUserByUsername(username string) (*dtos.User, error) {
	var user dtos.User
	err := repo.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *Repository) CreateUser(username, passwordHash string) (int, error) {
	var userId int
	// The RETURNING clause returns the id of the inserted user
	err := repo.DB.QueryRow("INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", username, passwordHash).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (repo *Repository) DeleteUser(username string) error {
	_, err := repo.DB.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) DeleteUserById(id int) error {
	_, err := repo.DB.Exec("DELETE FROM users where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
