package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
	ID             string
	ProjectID      string
	EndpointID     *string
	IdempotencyKey string
	Method         string
	Headers        []byte
	Payload        []byte
	Status         string
	AttemptCount   int
}

type Endpoint struct {
	ID             string
	DestinationURL string
}

type EventStore struct {
	db *pgxpool.Pool
}

func NewEventStore(db *pgxpool.Pool) *EventStore {
	return &EventStore{db: db}
}

func (s *EventStore) GetEventWithEndpoint(ctx context.Context, eventID string) (*Event, *Endpoint, error) {
	var e Event
	var ep Endpoint
	var epID *string
	var epURL *string

	err := s.db.QueryRow(ctx,
		`SELECT e.id, e.project_id, e.endpoint_id, e.idempotency_key,
		        e.method, e.headers, e.payload, e.status, e.attempt_count,
		        en.id, en.destination_url
		 FROM events e
		 LEFT JOIN endpoints en ON en.id = e.endpoint_id
		 WHERE e.id = $1`,
		eventID,
	).Scan(
		&e.ID, &e.ProjectID, &e.EndpointID, &e.IdempotencyKey,
		&e.Method, &e.Headers, &e.Payload, &e.Status, &e.AttemptCount,
		&epID, &epURL,
	)
	if err != nil {
		return nil, nil, err
	}

	if epID != nil {
		ep.ID = *epID
	}
	if epURL != nil {
		ep.DestinationURL = *epURL
	}

	return &e, &ep, nil
}

func (s *EventStore) MarkDelivered(ctx context.Context, eventID string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE events SET status='delivered', delivered_at=NOW(), attempt_count=attempt_count+1 WHERE id=$1`,
		eventID,
	)
	return err
}

func (s *EventStore) MarkFailed(ctx context.Context, eventID string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE events SET status='failed', attempt_count=attempt_count+1 WHERE id=$1`,
		eventID,
	)
	return err
}

func (s *EventStore) IncrementAttempt(ctx context.Context, eventID string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE events SET attempt_count=attempt_count+1 WHERE id=$1`,
		eventID,
	)
	return err
}

func (s *EventStore) LogAttempt(ctx context.Context, eventID string, statusCode int, body string, durationMs int, errMsg string, nextRetry *time.Time) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO delivery_attempts (event_id, status_code, response_body, duration_ms, error_message, next_retry_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		eventID, statusCode, body, durationMs, errMsg, nextRetry,
	)
	return err
}