# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go SDK for TiDB Cloud API, generated from the OpenAPI specification at https://docs-download.pingcap.com/api/tidbcloud-oas.json. The SDK provides programmatic access to TiDB Cloud resources including Projects, Clusters, Backups, and Restores.

## Architecture

### Directory Structure (Go Best Practices)
```
tidb-cloud-go/
├── cmd/                    # CLI applications
├── pkg/                    # Public library code
│   ├── client/            # Main SDK client
│   ├── models/            # Generated models from OpenAPI spec
│   └── auth/              # Authentication helpers
├── internal/              # Private application code
├── api/                   # API definitions and generated code
├── examples/              # Example usage
├── docs/                  # Documentation
└── testdata/              # Test fixtures
```

### Core Components

1. **Client Package** (`pkg/client/`)
   - Main TiDB Cloud client with HTTP Digest Authentication
   - Rate limiting (100 requests/minute)
   - Base URL: https://api.tidbcloud.com/api/v1-beta

2. **Models Package** (`pkg/models/`)
   - Generated from OpenAPI spec using swagger-codegen or oapi-codegen
   - Resource types: Project, Cluster, Backup, Restore
   - Supports AWS and GCP cloud providers

3. **Authentication** (`pkg/auth/`)
   - HTTP Digest Authentication implementation
   - API key management (public/private key pairs)
   - Secure handling of credentials

## Development Commands

### Code Generation
```bash
# Generate models from OpenAPI spec
go generate ./...

# Alternative: Using oapi-codegen
oapi-codegen -generate types -package models api/tidbcloud-oas.json > pkg/models/types.go
oapi-codegen -generate client -package client api/tidbcloud-oas.json > pkg/client/client.go
```

### Testing (TDD with t-wada style)
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestClusterCreate ./pkg/client

# Run tests with race detection
go test -race ./...

# Benchmark tests
go test -bench=. ./...
```

### Build and Lint
```bash
# Build
go build ./...

# Lint (requires golangci-lint)
golangci-lint run

# Format code
go fmt ./...

# Vet code
go vet ./...

# Tidy modules
go mod tidy
```

## TDD Implementation Guidelines

Following t-wada's TDD approach:

1. **Red**: Write failing tests first
2. **Green**: Write minimal code to make tests pass  
3. **Refactor**: Improve code while keeping tests green

### Test Structure
```go
// Table-driven tests for comprehensive coverage
func TestClusterOperations(t *testing.T) {
    tests := []struct {
        name     string
        input    ClusterCreateRequest
        expected ClusterResponse
        wantErr  bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## API Client Design

### HTTP Digest Authentication
- Implement digest authentication for API security
- Handle challenge-response flow
- Cache authentication state appropriately

### Rate Limiting
- Implement client-side rate limiting (100 req/min)
- Provide backoff and retry mechanisms
- Expose rate limit status to users

### Error Handling
- Wrap API errors with context
- Provide typed error responses
- Handle network timeouts and retries

## Key API Resources

### Projects
- List and create projects
- Project-level operations and permissions

### Clusters  
- Support both Serverless and Dedicated types
- Multi-cloud support (AWS, GCP)
- Cluster lifecycle management (create, modify, delete)
- Customer-Managed Encryption Keys support

### Backups
- Automated and manual backup operations
- Backup lifecycle management
- Cross-region backup support

### Restores
- Point-in-time recovery
- Restore from backups
- Cluster restoration workflows

## Development Notes

- OpenAPI spec is maintained externally - regenerate models when spec updates
- Follow Go naming conventions for public APIs
- Use context.Context for cancellation and timeouts
- Implement proper resource cleanup in tests
- Consider implementing SDK middleware for logging, metrics, etc.