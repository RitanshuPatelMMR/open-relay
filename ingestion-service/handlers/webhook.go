package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ritanshupatel/openrelay/ingestion-service/models"
	"github.com/ritanshupatel/openrelay/ingestion-service/queue"
	"github.com/ritanshupatel/openrelay/ingestion-service/store"
)

type WebhookHandler struct {
	events *store.EventStore
	queue  *queue.RedisQueue
}

func NewWebhookHandler(events *store.EventStore, q *queue.RedisQueue) *WebhookHandler {
	return &WebhookHandler{events: events, queue: q}
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := chi.URLParam(r, "projectID")

	// 1. check project exists
	project, err := h.events.GetProject(ctx, projectID)
	if err != nil {
		log.Println("db error fetching project:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if project == nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	// 2. get active endpoint
	endpoint, err := h.events.GetActiveEndpoint(ctx, projectID)
	if err != nil {
		log.Println("db error fetching endpoint:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 3. check idempotency
	idempotencyKey := r.Header.Get("X-Idempotency-Key")
	if idempotencyKey != "" {
		dup, err := h.events.IsDuplicate(ctx, projectID, idempotencyKey)
		if err != nil {
			log.Println("db error checking duplicate:", err)
		}
		if dup {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "already processed"})
			return
		}
	}

	// 4. read body (max 5MB)
	body, err := io.ReadAll(io.LimitReader(r.Body, 5*1024*1024))
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	// 5. headers as JSON
	headers, _ := json.Marshal(r.Header)

	// 6. build event
	var endpointID *string
	if endpoint != nil {
		endpointID = &endpoint.ID
	}

	sourceIP := r.RemoteAddr
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		sourceIP = strings.Split(fwd, ",")[0]
	}

	event := &models.Event{
		ProjectID:      projectID,
		EndpointID:     endpointID,
		IdempotencyKey: idempotencyKey,
		SourceIP:       sourceIP,
		Method:         r.Method,
		Path:           r.URL.Path,
		Headers:        headers,
		Payload:        body,
	}

	// 7. save to DB
	eventID, err := h.events.InsertEvent(ctx, event)
	if err != nil {
		log.Println("db error inserting event:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 8. push to Redis Streams
	if err := h.queue.PushEvent(ctx, eventID); err != nil {
		log.Println("redis push error:", err)
		// don't fail — event is saved in DB, can be recovered
	}

	log.Printf("event queued: %s for project: %s", eventID, projectID)

	// 9. return 200 immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "accepted",
		"event_id": eventID,
	})
}