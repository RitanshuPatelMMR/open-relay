package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritanshupatel/openrelay/ingestion-service/models"
)

type EventStore struct {
	db *pgxpool.Pool
}

func NewEventStore(db *pgxpool.Pool) *EventStore {
	return &EventStore{db: db}
}

// GetProject fetches project by ID
func (s *EventStore) GetProject(ctx context.Context, projectID string) (*models.Project, error) {
	var p models.Project
	err := s.db.QueryRow(ctx,
		`SELECT id, name, api_key, created_at FROM projects WHERE id = $1`,
		projectID,
	).Scan(&p.ID, &p.Name, &p.APIKey, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// GetActiveEndpoint fetches active endpoint for project
func (s *EventStore) GetActiveEndpoint(ctx context.Context, projectID string) (*models.Endpoint, error) {
	var e models.Endpoint
	err := s.db.QueryRow(ctx,
		`SELECT id, project_id, destination_url, description, is_active, created_at
		 FROM endpoints WHERE project_id = $1 AND is_active = true LIMIT 1`,
		projectID,
	).Scan(&e.ID, &e.ProjectID, &e.DestinationURL, &e.Description, &e.IsActive, &e.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// IsDuplicate checks if event with same idempotency key exists
func (s *EventStore) IsDuplicate(ctx context.Context, projectID, key string) (bool, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM events WHERE project_id = $1 AND idempotency_key = $2`,
		projectID, key,
	).Scan(&count)
	return count > 0, err
}

// InsertEvent saves new event to DB
func (s *EventStore) InsertEvent(ctx context.Context, e *models.Event) (string, error) {
	var id string
	err := s.db.QueryRow(ctx,
		`INSERT INTO events 
		 (project_id, endpoint_id, idempotency_key, source_ip, method, path, headers, payload, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
		 RETURNING id`,
		e.ProjectID, e.EndpointID, e.IdempotencyKey,
		e.SourceIP, e.Method, e.Path, e.Headers, e.Payload,
	).Scan(&id)
	return id, err
}