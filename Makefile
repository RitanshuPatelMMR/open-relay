dev:
	docker compose up -d

stop:
	docker compose down

test:
	go test ./...

build:
	docker compose build

logs:
	docker compose logs -f