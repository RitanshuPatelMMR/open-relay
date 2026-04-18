package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
	EndpointID     *string    `json:"endpoint_id"`
	IdempotencyKey *string    `json:"idempotency_key"`
	SourceIP       *string    `json:"source_ip"`
	Method         string     `json:"method"`
	Path           *string    `json:"path"`
	Headers        *string    `json:"headers"`
	Payload        *string    `json:"payload"`
	Status         string     `json:"status"`
	AttemptCount   int        `json:"attempt_count"`
	CreatedAt      time.Time  `json:"created_at"`
	DeliveredAt    *time.Time `json:"delivered_at"`
}

type DeliveryAttempt struct {
	ID           string     `json:"id"`
	EventID      string     `json:"event_id"`
	AttemptedAt  time.Time  `json:"attempted_at"`
	StatusCode   *int       `json:"status_code"`
	ResponseBody *string    `json:"response_body"`
	DurationMs   *int       `json:"duration_ms"`
	ErrorMessage *string    `json:"error_message"`
	NextRetryAt  *time.Time `json:"next_retry_at"`
}

func ListEvents(ctx context.Context, pool *pgxpool.Pool, projectID, status string, limit, offset int) ([]Event, error) {
	query := `SELECT id, project_id, endpoint_id, idempotency_key, source_ip, method, path,
		headers::text, payload::text, status, attempt_count, created_at, delivered_at
		FROM events WHERE 1=1`
	args := []any{}
	i := 1

	if projectID != "" {
		query += ` AND project_id = $` + itoa(i)
		args = append(args, projectID)
		i++
	}
	if status != "" {
		query += ` AND status = $` + itoa(i)
		args = append(args, status)
		i++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + itoa(i) + ` OFFSET $` + itoa(i+1)
	args = append(args, limit, offset)

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.EndpointID, &e.IdempotencyKey,
			&e.SourceIP, &e.Method, &e.Path, &e.Headers, &e.Payload,
			&e.Status, &e.AttemptCount, &e.CreatedAt, &e.DeliveredAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func GetEventByID(ctx context.Context, pool *pgxpool.Pool, id string) (*Event, error) {
	var e Event
	err := pool.QueryRow(ctx, `SELECT id, project_id, endpoint_id, idempotency_key, source_ip, method, path,
		headers::text, payload::text, status, attempt_count, created_at, delivered_at
		FROM events WHERE id = $1`, id).
		Scan(&e.ID, &e.ProjectID, &e.EndpointID, &e.IdempotencyKey,
			&e.SourceIP, &e.Method, &e.Path, &e.Headers, &e.Payload,
			&e.Status, &e.AttemptCount, &e.CreatedAt, &e.DeliveredAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetDeliveryAttempts(ctx context.Context, pool *pgxpool.Pool, eventID string) ([]DeliveryAttempt, error) {
	rows, err := pool.Query(ctx, `SELECT id, event_id, attempted_at, status_code, response_body,
		duration_ms, error_message, next_retry_at FROM delivery_attempts WHERE event_id = $1
		ORDER BY attempted_at ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []DeliveryAttempt
	for rows.Next() {
		var a DeliveryAttempt
		if err := rows.Scan(&a.ID, &a.EventID, &a.AttemptedAt, &a.StatusCode,
			&a.ResponseBody, &a.DurationMs, &a.ErrorMessage, &a.NextRetryAt); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}

func RequeueEvent(ctx context.Context, pool *pgxpool.Pool, id string) error {
	_, err := pool.Exec(ctx,
		`UPDATE events SET status='pending', attempt_count=0, delivered_at=NULL WHERE id=$1`, id)
	return err
}

func itoa(i int) string {
	return string(rune('0' + i))
}