package store

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("connected to database")
	return pool
}