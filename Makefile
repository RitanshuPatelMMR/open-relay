dev:
	docker compose up -d

stop:
	docker compose down

migrate:
	cd ingestion-service && go run cmd/migrate/main.go

test:
	go test ./...

build:
	docker compose build

logs:
	docker compose logs -f