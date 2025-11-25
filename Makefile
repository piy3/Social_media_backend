include .env
MIGRATIONS_PATH=./cmd/migrate/migrations
DB_INIT_SCRIPT=./scripts/db_init.sql

docker-up:
	@docker-compose up --build

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)

migrate-up:
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) up

migrate-down:
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) down 1

seed:
	@go run cmd/migrate/seed/main.go

.PHONY: docker-up migration migrate-up migrate-down seed