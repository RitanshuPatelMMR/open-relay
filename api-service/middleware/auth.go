package middleware

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	appdb "github.com/ritanshupatel/openrelay/api-service/db"
)

type contextKey string
const ProjectKey contextKey = "project"

func Auth(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, `{"error":"missing X-API-Key header"}`, http.StatusUnauthorized)
				return
			}
			project, err := appdb.GetProjectByAPIKey(r.Context(), pool, apiKey)
			if err != nil {
				http.Error(w, `{"error":"invalid api key"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ProjectKey, project)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}