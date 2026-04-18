package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/ritanshupatel/openrelay/worker-service/config"
	"github.com/ritanshupatel/openrelay/worker-service/deadletter"
	"github.com/ritanshupatel/openrelay/worker-service/store"
	"github.com/ritanshupatel/openrelay/worker-service/telemetry"
	"github.com/ritanshupatel/openrelay/worker-service/worker"
)

func main() {
	ctx := context.Background()
	shutdown := telemetry.InitTracer(ctx)
	defer shutdown()

	cfg := config.Load()

	db := store.NewDB(cfg.DBUrl)
	defer db.Close()

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal("failed to parse redis URL:", err)
	}
	redisClient := redis.NewClient(opts)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("redis ping failed:", err)
	}
	log.Println("worker: connected to redis")

	eventStore := store.NewEventStore(db)
	dlHandler := deadletter.NewHandler(db)

	pool := worker.NewPool(
		redisClient,
		eventStore,
		dlHandler,
		cfg.WorkerCount,
		cfg.MaxRetryAttempts,
		cfg.RequestTimeoutSeconds,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down worker...")
	cancel()
}