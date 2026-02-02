@echo off
echo Building Board API for Windows...
go build -o board-api.exe main.go
echo Build complete!
pause
