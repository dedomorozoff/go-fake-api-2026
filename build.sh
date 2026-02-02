#!/bin/bash
echo "Building Board API..."
GOOS=linux GOARCH=amd64 go build -o board-api-linux main.go
GOOS=windows GOARCH=amd64 go build -o board-api-windows.exe main.go
echo "Build complete!"
