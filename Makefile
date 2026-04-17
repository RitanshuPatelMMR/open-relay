dev:
	docker compose up -d

stop:
	docker compose down

migrate:
	cd ingestion-service && go run cmd/migrate/main.go

run-ingestion:
	cd ingestion-service && go run main.go

run-worker:
	cd worker-service && go run main.go

test:
	go test ./...

build:
	docker compose build

logs:
	docker compose logs -f