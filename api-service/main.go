package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/ritanshupatel/openrelay/api-service/config"
	appdb "github.com/ritanshupatel/openrelay/api-service/db"
	"github.com/ritanshupatel/openrelay/api-service/handlers"
	mw "github.com/ritanshupatel/openrelay/api-service/middleware"
	"github.com/ritanshupatel/openrelay/api-service/telemetry"
	ws "github.com/ritanshupatel/openrelay/api-service/websocket"
)

func main() {
	ctx := context.Background()
	shutdown := telemetry.InitTracer(ctx)
	defer shutdown()

	cfg := config.Load()
	pool := appdb.Connect(cfg.DBUrl)
	defer pool.Close()

	opt, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		log.Fatalf("Redis URL parse failed: %v", err)
	}
	rdb := redis.NewClient(opt)

	hub := ws.NewHub()
	go hub.Run()
	go hub.SubscribeRedis(rdb)

	eventsH := &handlers.EventsHandler{Pool: pool, Rdb: rdb}
	endpointsH := &handlers.EndpointsHandler{Pool: pool}
	projectsH := &handlers.ProjectsHandler{Pool: pool}
	analyticsH := &handlers.AnalyticsHandler{Pool: pool}

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(corsMiddleware)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})

	r.Group(func(r chi.Router) {
		r.Use(mw.Auth(pool))

		r.Get("/api/events", eventsH.List)
		r.Get("/api/events/{id}", eventsH.Get)
		r.Post("/api/events/{id}/replay", eventsH.Replay)

		r.Get("/api/endpoints", endpointsH.List)
		r.Post("/api/endpoints", endpointsH.Create)
		r.Put("/api/endpoints/{id}", endpointsH.Update)
		r.Delete("/api/endpoints/{id}", endpointsH.Delete)

		r.Get("/api/projects", projectsH.List)
		r.Post("/api/projects", projectsH.Create)

		r.Get("/api/analytics", analyticsH.Get)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("api-service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}