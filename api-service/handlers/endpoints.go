package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	appdb "github.com/ritanshupatel/openrelay/api-service/db"
	mw "github.com/ritanshupatel/openrelay/api-service/middleware"
)

type EndpointsHandler struct {
	Pool *pgxpool.Pool
}

func (h *EndpointsHandler) List(w http.ResponseWriter, r *http.Request) {
	project := r.Context().Value(mw.ProjectKey).(*appdb.Project)
	eps, err := appdb.ListEndpoints(r.Context(), h.Pool, project.ID)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	if eps == nil {
		eps = []appdb.Endpoint{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eps)
}

func (h *EndpointsHandler) Create(w http.ResponseWriter, r *http.Request) {
	project := r.Context().Value(mw.ProjectKey).(*appdb.Project)
	var body struct {
		DestinationURL string `json:"destination_url"`
		Description    string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.DestinationURL == "" {
		http.Error(w, `{"error":"destination_url required"}`, http.StatusBadRequest)
		return
	}
	ep, err := appdb.CreateEndpoint(r.Context(), h.Pool, project.ID, body.DestinationURL, body.Description)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ep)
}

func (h *EndpointsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		DestinationURL string `json:"destination_url"`
		IsActive       bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	ep, err := appdb.UpdateEndpoint(r.Context(), h.Pool, id, body.DestinationURL, body.IsActive)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ep)
}

func (h *EndpointsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := appdb.DeleteEndpoint(r.Context(), h.Pool, id); err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}