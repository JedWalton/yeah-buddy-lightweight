package uptimechecker

import (
	"i-couldve-got-six-reps/api/uptimechecker/types"
)

// Application Management
func (r *Repository) CreateApplication(userId int, name, description string) (int, error) {
	var applicationId int
	// Include the userId in the SQL INSERT statement
	query := `INSERT INTO Applications (user_id, name, description) VALUES ($1, $2, $3) RETURNING application_id`
	// Pass the userId along with name and description to the QueryRow function
	err := r.DB.QueryRow(query, userId, name, description).Scan(&applicationId)
	if err != nil {
		return 0, err
	}
	return applicationId, nil
}

func (r *Repository) UpdateApplication(applicationId int, name, description string) error {
	query := `UPDATE Applications SET name = $1, description = $2, updated_at = NOW() WHERE application_id = $3`
	_, err := r.DB.Exec(query, name, description, applicationId)
	return err
}

func (r *Repository) DeleteApplication(applicationId int) error {
	query := `DELETE FROM Applications WHERE application_id = $1`
	_, err := r.DB.Exec(query, applicationId)
	return err
}

func (r *Repository) FindApplication(applicationId int) (*types.Application, error) {
	var app types.Application
	query := `SELECT application_id, name, description FROM Applications WHERE application_id = $1`
	err := r.DB.QueryRow(query, applicationId).Scan(&app.ApplicationID, &app.Name, &app.Description)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
