package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ritanshupatel/openrelay/ingestion-service/config"
	"github.com/ritanshupatel/openrelay/ingestion-service/handlers"
	"github.com/ritanshupatel/openrelay/ingestion-service/store"
)

func main() {
	cfg := config.Load()

	db := store.NewDB(cfg.DBUrl)
	defer db.Close()

	eventStore := store.NewEventStore(db)
	webhookHandler := handlers.NewWebhookHandler(eventStore)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handlers.HealthHandler)
	r.Post("/in/{projectID}", webhookHandler.Handle)

	log.Printf("ingestion-service starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}