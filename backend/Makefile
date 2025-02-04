include .env

build:
	@go build -o ./bin/main ./cmd/api/main.go

run: build
	@./bin/main

build-seed:
	@go build -o bin/seed ./cmd/seed/main.go

seed: build-seed
	@./bin/seed -seed

seed-me: build-seed
	@./bin/seed -seed-me

create-migration:
	@GOOSE_MIGRATION_DIR=./internal/database/migrations goose postgres ${DB_MIGRATION_STRING} create $(name) sql

migrate:
	@GOOSE_MIGRATION_DIR=./internal/database/migrations goose postgres ${DB_MIGRATION_STRING} up

migrate-down:
	@GOOSE_MIGRATION_DIR=./internal/database/migrations goose postgres ${DB_MIGRATION_STRING} down

redo:
	@GOOSE_MIGRATION_DIR=./internal/database/migrations goose postgres ${DB_MIGRATION_STRING} redo
