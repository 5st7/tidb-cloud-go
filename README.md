# TiDB Cloud Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/5st7/tidb-cloud-go.svg)](https://pkg.go.dev/github.com/5st7/tidb-cloud-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/5st7/tidb-cloud-go)](https://goreportcard.com/report/github.com/5st7/tidb-cloud-go)

An unofficial Go SDK for the [TiDB Cloud](https://www.pingcap.com/tidb-cloud/) API that provides comprehensive support for managing TiDB Cloud resources including projects, clusters, backups, restores, and private endpoints.

## Features

- **Complete API Coverage**: All TiDB Cloud API endpoints are supported
- **HTTP Digest Authentication**: Secure authentication using API keys
- **Automatic Retry Logic**: Exponential backoff with intelligent retry policies
- **Context Support**: All operations support context for cancellation and timeouts
- **Comprehensive Error Handling**: Detailed error types with helper methods
- **Type Safety**: Full type safety with auto-generated models
- **Rate Limit Handling**: Built-in support for TiDB Cloud's rate limits
- **Private Endpoints**: Support for AWS PrivateLink and GCP Private Service Connect

## Installation

```bash
go get github.com/5st7/tidb-cloud-go
```

## Quick Start

### Prerequisites

1. Get your TiDB Cloud API credentials:
   - Go to [TiDB Cloud Console](https://console.tidbcloud.com/)
   - Navigate to Settings → API Keys
   - Create a new API key and save the public and private keys

2. Set environment variables:
```bash
export TIDB_CLOUD_PUBLIC_KEY="your-public-key"
export TIDB_CLOUD_PRIVATE_KEY="your-private-key"
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/5st7/tidb-cloud-go/pkg/client"
)

func main() {
    // Create client with API credentials
    client, err := client.NewClient(
        os.Getenv("TIDB_CLOUD_PUBLIC_KEY"),
        os.Getenv("TIDB_CLOUD_PRIVATE_KEY"),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // List all projects
    projects, err := client.ListProjects()
    if err != nil {
        log.Fatalf("Failed to list projects: %v", err)
    }

    fmt.Printf("Found %d projects:\n", len(projects.Items))
    for _, project := range projects.Items {
        fmt.Printf("- %s (%s)\n", *project.Name, *project.ID)
    }
}
```

## API Reference

### Projects

```go
// List all projects
projects, err := client.ListProjects()

// Create a new project
req := &models.OpenapiCreateProjectReq{
    Name: stringPtr("My New Project"),
}
project, err := client.CreateProject(req)
```

### Clusters

```go
// List clusters in a project
clusters, err := client.ListClusters(projectID)

// Get cluster details
cluster, err := client.GetCluster(projectID, clusterID)

// Create a new cluster
req := &models.OpenapiCreateClusterReq{
    Name:          stringPtr("my-cluster"),
    ClusterType:   stringPtr("DEDICATED"),
    CloudProvider: stringPtr("AWS"),
    Region:        stringPtr("us-west-2"),
    Config: &models.OpenapiClusterConfig{
        RootPassword: stringPtr("SecurePassword123!"),
        Port:         int64Ptr(4000),
        Components: &models.OpenapiClusterConfigComponents{
            TiDB: &models.OpenapiUpdateTiDBComponent{
                NodeSize:     stringPtr("8C16G"),
                NodeQuantity: int64Ptr(1),
            },
            TiKV: &models.OpenapiUpdateTiKVComponent{
                NodeSize:        stringPtr("8C32G"),
                NodeQuantity:    int64Ptr(3),
                StorageSizeGib:  int64Ptr(500),
            },
        },
    },
}
cluster, err := client.CreateCluster(projectID, req)

// Update cluster configuration
updateReq := &models.OpenapiUpdateClusterReq{
    Config: &models.OpenapiClusterConfig{
        Components: &models.OpenapiClusterConfigComponents{
            TiDB: &models.OpenapiUpdateTiDBComponent{
                NodeQuantity: int64Ptr(2), // Scale up
            },
        },
    },
}
cluster, err := client.UpdateCluster(projectID, clusterID, updateReq)

// Delete a cluster
err := client.DeleteCluster(projectID, clusterID)
```

### Backups

```go
// List backups for a cluster
backups, err := client.ListBackups(projectID, clusterID)

// Get backup details
backup, err := client.GetBackup(projectID, clusterID, backupID)

// Create a backup
req := &models.OpenapiCreateBackupReq{
    Name:        stringPtr("my-backup"),
    Description: stringPtr("Daily backup"),
}
backup, err := client.CreateBackup(projectID, clusterID, req)

// Delete a backup
err := client.DeleteBackup(projectID, clusterID, backupID)
```

### Restores

```go
// List restores in a project
restores, err := client.ListRestores(projectID)

// Get restore details
restore, err := client.GetRestore(projectID, restoreID)

// Create a restore (new cluster from backup)
req := &models.OpenapiCreateRestoreReq{
    BackupID: stringPtr(backupID),
    Name:     stringPtr("restored-cluster"),
    Config: &models.OpenapiClusterConfig{
        RootPassword: stringPtr("NewPassword123!"),
        // ... cluster configuration
    },
}
restore, err := client.CreateRestore(projectID, req)
```

### Private Endpoints

```go
ctx := context.Background()

// Create private endpoint service
service, err := client.CreatePrivateEndpointService(ctx, projectID, clusterID)

// Get private endpoint service details
service, err := client.GetPrivateEndpointService(ctx, projectID, clusterID)

// List private endpoints for a cluster
endpoints, err := client.ListPrivateEndpoints(ctx, projectID, clusterID)

// Create a private endpoint
req := &models.OpenapiCreatePrivateEndpointReq{
    EndpointName: stringPtr("vpce-1234567890abcdef0"), // Your VPC endpoint ID
}
endpoint, err := client.CreatePrivateEndpoint(ctx, projectID, clusterID, req)

// Delete a private endpoint
err := client.DeletePrivateEndpoint(ctx, projectID, clusterID, endpointID)

// List all private endpoints in a project
endpoints, err := client.ListPrivateEndpointsOfProject(ctx, projectID)
```

### Provider Regions

```go
// List available cloud provider regions
regions, err := client.ListProviderRegions()

for _, region := range regions.Items {
    if region.Available != nil && *region.Available {
        fmt.Printf("Available: %s %s\n", *region.CloudProvider, *region.Region)
    }
}
```

## Error Handling

The SDK provides comprehensive error handling with specific error types:

```go
import "github.com/5st7/tidb-cloud-go/pkg/errors"

clusters, err := client.ListClusters(projectID)
if err != nil {
    // Check if it's a TiDB Cloud API error
    if apiErr, ok := err.(errors.APIError); ok {
        fmt.Printf("API Error: %s (Status: %d, Code: %d)\n", 
            apiErr.Message, apiErr.StatusCode, apiErr.Code)
        
        // Check specific error types
        switch {
        case apiErr.IsAuthenticationError():
            fmt.Println("Invalid API credentials")
        case apiErr.IsAuthorizationError():
            fmt.Println("Insufficient permissions")
        case apiErr.IsRateLimitError():
            fmt.Println("Rate limit exceeded - SDK will automatically retry")
        case apiErr.IsNotFoundError():
            fmt.Println("Resource not found")
        case apiErr.IsBadRequestError():
            fmt.Println("Invalid request parameters")
        case apiErr.IsRetryable():
            fmt.Println("Temporary error - SDK will automatically retry")
        }
    }
}
```

## Context and Timeouts

All API operations support context for cancellation and timeouts:

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use context with API calls
service, err := client.GetPrivateEndpointService(ctx, projectID, clusterID)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("Request timed out")
    }
}
```

## Rate Limiting and Retries

The SDK automatically handles rate limiting and retries:

- **Rate Limit**: 100 requests per minute per API key (enforced by TiDB Cloud)
- **Automatic Retries**: Exponential backoff for retryable errors
- **Max Attempts**: 3 attempts (initial + 2 retries)
- **Backoff**: 1s, 2s, 4s, capped at 30s

```go
// The client automatically retries on:
// - Rate limit errors (429)
// - Server errors (500, 502, 503, 504)
// - Network errors

// No retries for:
// - Authentication errors (401)
// - Authorization errors (403)
// - Bad request errors (400)
// - Not found errors (404)
```

