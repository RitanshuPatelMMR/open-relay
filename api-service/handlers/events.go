package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	appdb "github.com/ritanshupatel/openrelay/api-service/db"
	mw "github.com/ritanshupatel/openrelay/api-service/middleware"
)

type EventsHandler struct {
	Pool *pgxpool.Pool
	Rdb  *redis.Client
}

func (h *EventsHandler) List(w http.ResponseWriter, r *http.Request) {
	project := r.Context().Value(mw.ProjectKey).(*appdb.Project)
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	events, err := appdb.ListEvents(r.Context(), h.Pool, project.ID, status, limit, offset)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	if events == nil {
		events = []appdb.Event{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	event, err := appdb.GetEventByID(r.Context(), h.Pool, id)
	if err != nil {
		http.Error(w, `{"error":"event not found"}`, http.StatusNotFound)
		return
	}
	attempts, _ := appdb.GetDeliveryAttempts(r.Context(), h.Pool, id)
	if attempts == nil {
		attempts = []appdb.DeliveryAttempt{}
	}
	resp := map[string]any{"event": event, "delivery_attempts": attempts}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *EventsHandler) Replay(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := appdb.RequeueEvent(r.Context(), h.Pool, id); err != nil {
		http.Error(w, `{"error":"failed to requeue"}`, http.StatusInternalServerError)
		return
	}
	// push back to Redis stream
	h.Rdb.XAdd(r.Context(), &redis.XAddArgs{
		Stream: "events",
		Values: map[string]any{"event_id": id},
	})
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"requeued"}`))
}