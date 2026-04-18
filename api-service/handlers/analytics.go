package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	appdb "github.com/ritanshupatel/openrelay/api-service/db"
	mw "github.com/ritanshupatel/openrelay/api-service/middleware"
)

type AnalyticsHandler struct {
	Pool *pgxpool.Pool
}

type HourlyBucket struct {
	Hour      string `json:"hour"`
	Total     int    `json:"total"`
	Delivered int    `json:"delivered"`
	Failed    int    `json:"failed"`
}

func (h *AnalyticsHandler) Get(w http.ResponseWriter, r *http.Request) {
	project := r.Context().Value(mw.ProjectKey).(*appdb.Project)

	rows, err := h.Pool.Query(r.Context(), `
		SELECT
			date_trunc('hour', created_at) AS hour,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'delivered') AS delivered,
			COUNT(*) FILTER (WHERE status = 'failed') AS failed
		FROM events
		WHERE project_id = $1
		  AND created_at > NOW() - INTERVAL '24 hours'
		GROUP BY hour
		ORDER BY hour ASC
	`, project.ID)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var buckets []HourlyBucket
	for rows.Next() {
		var b HourlyBucket
		if err := rows.Scan(&b.Hour, &b.Total, &b.Delivered, &b.Failed); err != nil {
			continue
		}
		buckets = append(buckets, b)
	}
	if buckets == nil {
		buckets = []HourlyBucket{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"hourly": buckets})
}