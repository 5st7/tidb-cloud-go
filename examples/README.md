# TiDB Cloud Go SDK Examples

This directory contains example code demonstrating how to use the TiDB Cloud Go SDK.

## Prerequisites

Before running these examples, you need:

1. A TiDB Cloud account
2. API credentials (public and private keys)
3. Go 1.19 or later installed

## Setting Up Authentication

Set your API credentials as environment variables:

```bash
export TIDB_CLOUD_API_PUBLIC_KEY="your-public-key"
export TIDB_CLOUD_API_PRIVATE_KEY="your-private-key"
```

## Available Examples

### list_projects
Lists all projects in your TiDB Cloud organization.

```bash
cd list_projects
go run main.go
```

### basic_usage
Demonstrates basic SDK operations including client initialization and simple API calls.

```bash
cd basic_usage
go run main.go
```

### cluster_management
Shows how to create, modify, and delete TiDB clusters.

```bash
cd cluster_management
go run main.go
```

### backup_restore
Demonstrates backup creation and restore operations.

```bash
cd backup_restore
go run main.go
```

### private_endpoints
Shows how to manage private endpoints for secure cluster access.

```bash
cd private_endpoints
go run main.go
```

## Running Examples

Each example can be run independently:

```bash
cd <example-directory>
go run main.go
```

## Common Patterns

All examples follow these patterns:

1. **Client Initialization**: Create a client with API credentials
2. **Error Handling**: Proper error checking and logging
3. **Resource Cleanup**: Clean up created resources when applicable
4. **Configuration**: Use environment variables for sensitive data

## Notes

- Examples may create real resources in your TiDB Cloud account
- Some operations (like cluster creation) may incur costs
- Always review the code before running to understand what resources will be created
- Remember to clean up any test resources after running examples