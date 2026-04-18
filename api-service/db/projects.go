package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

func ListProjects(ctx context.Context, pool *pgxpool.Pool) ([]Project, error) {
	rows, err := pool.Query(ctx,
		`SELECT id, name, api_key, created_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.APIKey, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func CreateProject(ctx context.Context, pool *pgxpool.Pool, name string) (*Project, error) {
	var p Project
	err := pool.QueryRow(ctx,
		`INSERT INTO projects (name) VALUES ($1)
		RETURNING id, name, api_key, created_at`, name).
		Scan(&p.ID, &p.Name, &p.APIKey, &p.CreatedAt)
	return &p, err
}

func GetProjectByAPIKey(ctx context.Context, pool *pgxpool.Pool, apiKey string) (*Project, error) {
	var p Project
	err := pool.QueryRow(ctx,
		`SELECT id, name, api_key, created_at FROM projects WHERE api_key=$1`, apiKey).
		Scan(&p.ID, &p.Name, &p.APIKey, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}