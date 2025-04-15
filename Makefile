EXECUTABLE_NAME="./ivy"
SOURCE_FILE="./cmd/tui/main.go"

cli:
	@go run ./cmd/cli/main.go

tui:
	@go run ./cmd/tui/main.go

build:
	@go build -ldflags "-s -w" -o $(EXECUTABLE_NAME) $(SOURCE_FILE)

release-local: # local-only release
	@goreleaser release --snapshot --clean

sql: # Connect to the database
	sqlite3 $(HOME)/Library/Caches/ivy_lee_todo/data.db

