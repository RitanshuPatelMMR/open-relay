package deadletter

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Handle(ctx context.Context, eventID string, reason string) {
	log.Printf("💀 DEAD LETTER: event %s failed all retries. Reason: %s", eventID, reason)

	_, err := h.db.Exec(ctx,
		`UPDATE events SET status='failed', attempt_count=attempt_count+1 WHERE id=$1`,
		eventID,
	)
	if err != nil {
		log.Println("failed to mark event as failed:", err)
	}
}