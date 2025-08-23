ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

DB_URL=postgres://$(DATABASE_USERNAME):$(DATABASE_PASSWORD)@$(DATABASE_ADDRESS):$(DATABASE_PORT)/$(DATABASE_DBNAME)?sslmode=disable

migrate-create:
	@read -p "Enter migration name: " name; \
	go run github.com/pressly/goose/v3/cmd/goose -dir ./migration create $$name sql

migrate-up:
	go run github.com/pressly/goose/v3/cmd/goose -dir ./migration postgres "$(DB_URL)" up

migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose -dir ./migration postgres "$(DB_URL)" down

migrate-status:
	go run github.com/pressly/goose/v3/cmd/goose -dir ./migration postgres "$(DB_URL)" status

vendor:
	@echo "Running go mod vendor..."
	@go mod vendor

run:
	@echo "Running app..."
	@go run cmd/gql/main.go

generate:
	@echo "Running generate file..."
	@go generate ./...

gql-gen:
	@echo "Running gql generator..."
	@go run github.com/99designs/gqlgen generate

test:
	@echo "Running tests..."
	@go test -v $$(go list ./... | grep -v /vendor/ | grep -v /graph/ | grep -v /cmd/)