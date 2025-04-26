# Project binary name
BINARY_NAME = permcheck

# Default target
.PHONY: all
all: build

# Build the project
.PHONY: build
build:
	go build -o $(BINARY_NAME) ./main.go

# Run the project
.PHONY: run
run:
	go run ./main.go scan ./mock_files

# Run tests
.PHONY: test
test:
	go test -v ./...

# Clean build files
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy
