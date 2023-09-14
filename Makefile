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

test-all:
	go test -v --failfast ./test/...