package models

import "time"

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type Endpoint struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	DestinationURL string    `json:"destination_url"`
	Description    string    `json:"description"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type Event struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
EndpointID  *string `json:"endpoint_id"`
	IdempotencyKey string     `json:"idempotency_key"`
	SourceIP       string     `json:"source_ip"`
	Method         string     `json:"method"`
	Path           string     `json:"path"`
	Headers        []byte     `json:"headers"`
	Payload        []byte     `json:"payload"`
	Status         string     `json:"status"`
	AttemptCount   int        `json:"attempt_count"`
	CreatedAt      time.Time  `json:"created_at"`
	DeliveredAt    *time.Time `json:"delivered_at"`
}