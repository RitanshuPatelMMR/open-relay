package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Endpoint struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	DestinationURL string    `json:"destination_url"`
	Description    *string   `json:"description"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

func ListEndpoints(ctx context.Context, pool *pgxpool.Pool, projectID string) ([]Endpoint, error) {
	rows, err := pool.Query(ctx,
		`SELECT id, project_id, destination_url, description, is_active, created_at
		FROM endpoints WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eps []Endpoint
	for rows.Next() {
		var e Endpoint
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.DestinationURL,
			&e.Description, &e.IsActive, &e.CreatedAt); err != nil {
			return nil, err
		}
		eps = append(eps, e)
	}
	return eps, nil
}

func CreateEndpoint(ctx context.Context, pool *pgxpool.Pool, projectID, destinationURL, description string) (*Endpoint, error) {
	var e Endpoint
	err := pool.QueryRow(ctx,
		`INSERT INTO endpoints (project_id, destination_url, description)
		VALUES ($1, $2, $3) RETURNING id, project_id, destination_url, description, is_active, created_at`,
		projectID, destinationURL, nullStr(description)).
		Scan(&e.ID, &e.ProjectID, &e.DestinationURL, &e.Description, &e.IsActive, &e.CreatedAt)
	return &e, err
}

func UpdateEndpoint(ctx context.Context, pool *pgxpool.Pool, id, destinationURL string, isActive bool) (*Endpoint, error) {
	var e Endpoint
	err := pool.QueryRow(ctx,
		`UPDATE endpoints SET destination_url=$1, is_active=$2 WHERE id=$3
		RETURNING id, project_id, destination_url, description, is_active, created_at`,
		destinationURL, isActive, id).
		Scan(&e.ID, &e.ProjectID, &e.DestinationURL, &e.Description, &e.IsActive, &e.CreatedAt)
	return &e, err
}

func DeleteEndpoint(ctx context.Context, pool *pgxpool.Pool, id string) error {
	_, err := pool.Exec(ctx, `DELETE FROM endpoints WHERE id=$1`, id)
	return err
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}