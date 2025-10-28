.PHONY: build install clean test build-all

# Build for current platform
build:
	go build -o bin/noti ./cmd/noti

# Install to GOPATH/bin
install:
	go install ./cmd/noti

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Cross-compile for all platforms
build-all:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/noti-linux-amd64 ./cmd/noti
	GOOS=linux GOARCH=arm64 go build -o dist/noti-linux-arm64 ./cmd/noti
	GOOS=darwin GOARCH=amd64 go build -o dist/noti-darwin-amd64 ./cmd/noti
	GOOS=darwin GOARCH=arm64 go build -o dist/noti-darwin-arm64 ./cmd/noti
	GOOS=windows GOARCH=amd64 go build -o dist/noti-windows-amd64.exe ./cmd/noti

# Run the CLI
run:
	go run ./cmd/noti
