package uptimechecker

import "i-couldve-got-six-reps/api/uptimechecker/types"

// Application Management
func (r *Repository) CreateApplication(name, description string) (int, error) {
	var applicationId int
	query := `INSERT INTO Applications (name, description) VALUES ($1, $2) RETURNING application_id`
	err := r.DB.QueryRow(query, name, description).Scan(&applicationId)
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
