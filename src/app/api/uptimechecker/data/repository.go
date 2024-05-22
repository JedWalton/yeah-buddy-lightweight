package uptimechecker

import "database/sql"

type Repository struct {
	DB *sql.DB
}

func NewUptimeCheckerRepository(db *sql.DB) Repository {
	return Repository{DB: db}
}
