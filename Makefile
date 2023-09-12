ifeq ($(wildcard .env),)
    # If .env doesn't exist, use environment variables
    DB_URI ?= $(shell echo $$DB_URI)
else
    # If .env exists, load variables from it
    include .env
endif

run:
	go run cmd/main.go

tidy:
	go mod tidy

migrate-up:
	migrate -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_ADDRESS}:${DB_PORT})/${DB_NAME}" -path database/migrations up

migrate-down:
	migrate -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_ADDRESS}:${DB_PORT})/${DB_NAME}" -path database/migrations down

test-all:
	go test -v --failfast ./test/...