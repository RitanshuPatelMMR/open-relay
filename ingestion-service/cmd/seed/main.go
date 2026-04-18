package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("connect error:", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), `
		INSERT INTO projects (id, name, api_key)
		VALUES ('6c410a13-68b3-4832-a254-e2b377e403f2', 'test', 'a6debfc9-974e-45f8-bd2f-003be2c006c3')
		ON CONFLICT DO NOTHING;

		INSERT INTO endpoints (id, project_id, destination_url, description)
VALUES ('6cf7ea53-55e5-471c-b1a0-3473768e46fc', '6c410a13-68b3-4832-a254-e2b377e403f2', 'https://webhook.site/db140908-6f2a-4689-a690-c81789e8440a', 'test endpoint')
ON CONFLICT (id) DO UPDATE SET description = 'test endpoint';
	`)
	if err != nil {
		log.Fatal("seed error:", err)
	}
	log.Println("seeded successfully")
}