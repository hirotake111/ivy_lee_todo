EXECUTABLE_NAME="./ivy"
SOURCE_FILE="./cmd/tui/main.go"

cli:
	@go run ./cmd/cli/main.go

tui:
	@go run ./cmd/tui/main.go

build:
	@go build -ldflags "-s -w" -o $(EXECUTABLE_NAME) $(SOURCE_FILE)


