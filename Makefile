run:
	cd cmd && go run main.go
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
build:
	cd cmd && go build -o Binary main.go
run-build:
	cd cmd && ./Binary
