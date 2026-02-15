.PHONY: dev build test lint fmt migrate-up migrate-down migrate-create sqlc docker-build docker-up-hotreload

# Default DB URL for local development
DB_URL ?= postgres://studytrack:studytrack@localhost:5432/studytrack?sslmode=disable
MIGRATIONS_DIR := db/migrations

dev:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -v -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	go tool golangci-lint run ./...

fmt:
	gofmt -w .
	go tool goimports -w .

migrate-up:
	go run ./cmd/migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	go run ./cmd/migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-create:
	go run ./cmd/migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq

sqlc:
	go tool sqlc generate

docker-build:
	docker build -t studytrack-api .

docker-up-hotreload:
	docker compose up --build

docker-up:
	docker compose up -d

docker-down:
	docker compose down
