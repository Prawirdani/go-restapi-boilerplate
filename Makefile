run:
	go run cmd/main.go
tidy:
	go mod tidy
swag:
	swag init -g cmd/main.go
swag-fmt:
	swag fmt
dev:
	make swag
	make swag-fmt
	make run

test-all:
	go test -v --failfast ./test/...