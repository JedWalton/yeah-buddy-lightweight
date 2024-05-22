package uptimechecker

import (
	"i-couldve-got-six-reps/api/uptimechecker/types"
	"log"
)

// Endpoint Management
func (r *Repository) AddEndpoint(applicationId int, url string, monitoringInterval int) (int, error) {
	var endpointId int
	query := `INSERT INTO Endpoints (application_id, url, monitoring_interval) VALUES ($1, $2, $3) RETURNING endpoint_id`
	err := r.DB.QueryRow(query, applicationId, url, monitoringInterval).Scan(&endpointId)
	if err != nil {
		return 0, err
	}
	return endpointId, nil
}

func (r *Repository) UpdateEndpoint(endpointId int, url string, monitoringInterval int, isActive bool) error {
	query := `UPDATE Endpoints SET url = $1, monitoring_interval = $2, is_active = $3, updated_at = NOW() WHERE endpoint_id = $4`
	_, err := r.DB.Exec(query, url, monitoringInterval, isActive, endpointId)
	return err
}

func (r *Repository) DeleteEndpoint(endpointId int) error {
	query := `DELETE FROM Endpoints WHERE endpoint_id = $1`
	_, err := r.DB.Exec(query, endpointId)
	return err
}

func (r *Repository) GetEndpoint(endpointId int) (*types.Endpoint, error) {
	var endpoint types.Endpoint
	query := `SELECT endpoint_id, application_id, url, monitoring_interval, is_active FROM Endpoints WHERE endpoint_id = $1`
	err := r.DB.QueryRow(query, endpointId).Scan(&endpoint.EndpointID, &endpoint.ApplicationID, &endpoint.URL, &endpoint.MonitoringInterval, &endpoint.IsActive)
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

// Method to list all active endpoints
func (r *Repository) ListActiveEndpoints() ([]types.Endpoint, error) {
	var endpoints []types.Endpoint
	query := `SELECT endpoint_id, application_id, url, monitoring_interval, is_active FROM Endpoints WHERE is_active = TRUE`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ep types.Endpoint
		err := rows.Scan(&ep.EndpointID, &ep.ApplicationID, &ep.URL, &ep.MonitoringInterval, &ep.IsActive)
		if err != nil {
			log.Printf("Failed to scan endpoint: %v", err)
			continue
		}
		endpoints = append(endpoints, ep)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return endpoints, nil
}
