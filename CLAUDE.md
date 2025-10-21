# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based CLI tool called `ctx` designed to simplify Kubernetes context and namespace management. The tool provides an easier interface for switching between Kubernetes contexts and setting default namespaces, with interactive selection capabilities.

## Tech Stack

- **Language**: Go 1.25.1
- **Module**: `github.com/uniquejava/ctx`
- **Target**: Command-line tool for Kubernetes context management
- **Dependencies**: None currently (fresh Go module)

## Development Environment Setup

### Prerequisites
- Go 1.25.1 or later
- kubectl installed and configured
- Access to `~/.kube/config` file

### Initial Setup
```bash
# Initialize Go module (already done)
go mod init github.com/uniquejava/ctx

# Download dependencies when added
go mod tidy

# Build the project
go build -o ctx .
```

## Essential Commands

### Building and Running
```bash
# Build the binary
go build -o ctx .

# Run directly during development
go run . [command]

# Install to system (optional)
go install .
```

### Testing
```bash
# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestFunctionName
```

### Development Workflow
```bash
# Format code
go fmt ./...

# Lint code
go vet ./...

# Tidy dependencies
go mod tidy

# Clean build cache
go clean -cache
```

## Project Structure

```
ctx/
├── go.mod              # Go module definition
├── REQUIREMENTS.md     # Project requirements (in Chinese)
├── CLAUDE.md          # This file - development guide
└── .idea/             # IntelliJ IDEA configuration
    ├── .gitignore
    ├── ctx.iml
    ├── encodings.xml
    ├── modules.xml
    └── workspace.xml
```

## Requirements Summary

The tool should implement the following commands:

1. `ctx ls` - List all contexts from `~/.kube/config`, mark current with `*`
2. `ctx use xxx` - Switch to context `xxx`
3. `ctx use xxx:default-ns` - Switch context and set default namespace
4. `ctx rm xxx` - Remove context `xxx`
5. Interactive selection with arrow keys and Enter (if possible)

## Key Implementation Notes

### Kubernetes Configuration
- Primary file location: `~/.kube/config`
- File permissions are typically restrictive (600)
- Contains contexts, clusters, and users configuration
- Current context is tracked in the file

### Expected Go Packages to Use
- `os` and `io/ioutil` for file operations
- `flag` package for command-line argument parsing
- `gopkg.in/yaml.v2` or similar for YAML parsing
- Possibly `github.com/spf13/cobra` for CLI framework (optional)
- Terminal UI libraries for interactive selection (if implemented)

### Code Organization Recommendations
- `main.go` - Entry point and command routing
- `config.go` - Kubernetes config file handling
- `context.go` - Context management operations
- `ui.go` - Interactive interface (if implemented)

## Development Guidelines

### Go Best Practices
- Follow standard Go project layout
- Use meaningful variable and function names
- Handle errors appropriately
- Use proper Go formatting (`go fmt`)
- Write tests for core functionality

### CLI Design
- Provide clear error messages
- Handle edge cases (invalid context names, permission issues)
- Support both programmatic and interactive use
- Follow standard CLI conventions

### Security Considerations
- Handle kube config file permissions carefully
- Validate user inputs
- Don't expose sensitive information unnecessarily

## Testing Strategy

### Unit Tests
- Test config file parsing
- Test context switching logic
- Test error handling scenarios
- Mock file system operations where appropriate

### Integration Testing
- Test with actual kube config files
- Test context switching with kubectl verification
- Test interactive mode functionality

## Build and Distribution

### Building for Multiple Platforms
```bash
# Build for different OS/architectures
GOOS=linux GOARCH=amd64 go build -o ctx-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o ctx-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o ctx-windows-amd64.exe .
```

### Version Management
- Consider using semantic versioning
- Add version command (`ctx version`)
- Tag releases in Git