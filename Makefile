.PHONY: dev build run clean test migrate-up migrate-down migrate-status migrate-create sqlc docs

MIGRATIONS_DIR=db/migrations
DB_URL=$(shell grep DB_URL .env | cut -d '=' -f2-)

dev:
	air

build:
	go build -o ./tmp/main ./cmd

run:
	go run ./cmd

clean:
	rm -rf ./tmp

test:
	go test -v ./internal/...

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

migrate-create:
	@read -p "Migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) -s create $$name sql

sqlc:
	sqlc generate

docs:
	swag init -g cmd/api/main.go --output docs
