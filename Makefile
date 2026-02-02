.PHONY: build run clean test dev

BINARY_NAME=board-api.exe

build:
	go build -o $(BINARY_NAME) main.go

run:
	go run main.go

clean:
	go clean
	rm -f $(BINARY_NAME)

dev:
	air

test:
	go test ./...
