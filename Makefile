.PHONY: build run clean test dev

BINARY_NAME=board-api.exe

build:
	go build -o $(BINARY_NAME) main.go

build-all:
	./build.sh

run:
	go run main.go

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f board-api-linux board-api-windows.exe board-api-freebsd

dev:
	air

test:
	go test ./...

api-test:
	chmod +x test_api.sh
	./test_api.sh
