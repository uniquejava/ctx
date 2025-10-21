.PHONY: build clean install test help

# Build the binary to bin directory
build:
	@echo "Building ctx..."
	@mkdir -p bin
	@go build -o bin/ctx .
	@echo "✅ Built successfully: bin/ctx"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f *.test *.out
	@echo "✅ Cleaned successfully"

# Install to system PATH (requires sudo)
install: build
	@echo "Installing ctx to /usr/local/bin..."
	@sudo cp bin/ctx /usr/local/bin/
	@echo "✅ Installed successfully: /usr/local/bin/ctx"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Show help
help:
	@echo "Available commands:"
	@echo "  build    - Build the ctx binary to bin/ctx"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install ctx to /usr/local/bin (requires sudo)"
	@echo "  test     - Run tests"
	@echo "  help     - Show this help message"

# Quick development build
dev: clean build
	@echo "✅ Development build complete"