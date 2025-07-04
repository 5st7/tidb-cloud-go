// Package main demonstrates basic usage of the TiDB Cloud Go SDK.
// This example shows how to authenticate, list projects and clusters,
// and handle errors properly.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/5st7/tidb-cloud-go/pkg/client"
	"github.com/5st7/tidb-cloud-go/pkg/errors"
)

func main() {
	// Get API credentials from environment variables
	publicKey := os.Getenv("TIDB_CLOUD_PUBLIC_KEY")
	privateKey := os.Getenv("TIDB_CLOUD_PRIVATE_KEY")

	if publicKey == "" || privateKey == "" {
		log.Fatal("Please set TIDB_CLOUD_PUBLIC_KEY and TIDB_CLOUD_PRIVATE_KEY environment variables")
	}

	// Create a new TiDB Cloud client
	client, err := client.NewClient(publicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example 1: List all projects
	fmt.Println("=== Listing Projects ===")
	projects, err := client.ListProjects()
	if err != nil {
		handleError("listing projects", err)
		return
	}

	fmt.Printf("Found %d projects:\n", len(projects.Items))
	for _, project := range projects.Items {
		fmt.Printf("- ID: %s, Name: %s, Clusters: %d\n",
			safeString(project.ID),
			safeString(project.Name),
			safeInt64(project.ClusterCount))
	}

	// Example 2: List clusters for the first project
	if len(projects.Items) > 0 {
		projectID := safeString(projects.Items[0].ID)
		fmt.Printf("\n=== Listing Clusters for Project %s ===\n", projectID)

		clusters, err := client.ListClusters(projectID)
		if err != nil {
			handleError("listing clusters", err)
			return
		}

		fmt.Printf("Found %d clusters:\n", len(clusters.Items))
		for _, cluster := range clusters.Items {
			fmt.Printf("- ID: %s, Name: %s, Type: %s, Provider: %s, Region: %s, Status: %s\n",
				safeString(cluster.ID),
				safeString(cluster.Name),
				safeString(cluster.ClusterType),
				safeString(cluster.CloudProvider),
				safeString(cluster.Region),
				safeString(cluster.Status.ClusterStatus))
		}

		// Example 3: List backups for the first cluster
		if len(clusters.Items) > 0 {
			clusterID := safeString(clusters.Items[0].ID)
			fmt.Printf("\n=== Listing Backups for Cluster %s ===\n", clusterID)

			backups, err := client.ListBackups(projectID, clusterID)
			if err != nil {
				handleError("listing backups", err)
				return
			}

			fmt.Printf("Found %d backups:\n", len(backups.Items))
			for _, backup := range backups.Items {
				fmt.Printf("- ID: %s, Name: %s, Type: %s, Status: %s, Size: %d bytes\n",
					safeString(backup.ID),
					safeString(backup.Name),
					safeString(backup.Type),
					safeString(backup.Status.BackupStatus),
					safeInt64(backup.BackupSizeBytes))
			}

			// Example 4: List private endpoints for the cluster
			fmt.Printf("\n=== Listing Private Endpoints for Cluster %s ===\n", clusterID)

			endpoints, err := client.ListPrivateEndpoints(ctx, projectID, clusterID)
			if err != nil {
				handleError("listing private endpoints", err)
				return
			}

			fmt.Printf("Found %d private endpoints:\n", len(endpoints.Items))
			for _, endpoint := range endpoints.Items {
				fmt.Printf("- ID: %s, Name: %s, Provider: %s, Status: %s\n",
					safeString(endpoint.ID),
					safeString(endpoint.EndpointName),
					safeString(endpoint.CloudProvider),
					safeString(endpoint.Status))
			}
		}
	}

	// Example 5: List provider regions
	fmt.Println("\n=== Listing Provider Regions ===")
	regions, err := client.ListProviderRegions()
	if err != nil {
		handleError("listing provider regions", err)
		return
	}

	fmt.Printf("Found %d regions:\n", len(regions.Items))
	for _, region := range regions.Items {
		status := "unavailable"
		if region.Available != nil && *region.Available {
			status = "available"
		}
		fmt.Printf("- Provider: %s, Region: %s, Status: %s\n",
			safeString(region.CloudProvider),
			safeString(region.Region),
			status)
	}

	fmt.Println("\n=== Example completed successfully ===")
}

// handleError demonstrates proper error handling for TiDB Cloud API errors
func handleError(operation string, err error) {
	fmt.Printf("Error %s: %v\n", operation, err)

	// Check if it's a TiDB Cloud API error
	if apiErr, ok := err.(errors.APIError); ok {
		fmt.Printf("API Error Details:\n")
		fmt.Printf("  Status Code: %d\n", apiErr.StatusCode)
		fmt.Printf("  Error Code: %d\n", apiErr.Code)
		fmt.Printf("  Message: %s\n", apiErr.Message)

		// Check specific error types
		switch {
		case apiErr.IsAuthenticationError():
			fmt.Println("  Type: Authentication Error - Check your API credentials")
		case apiErr.IsAuthorizationError():
			fmt.Println("  Type: Authorization Error - Insufficient permissions")
		case apiErr.IsRateLimitError():
			fmt.Println("  Type: Rate Limit Error - Too many requests")
		case apiErr.IsNotFoundError():
			fmt.Println("  Type: Not Found Error - Resource doesn't exist")
		case apiErr.IsBadRequestError():
			fmt.Println("  Type: Bad Request Error - Invalid request parameters")
		case apiErr.IsRetryable():
			fmt.Println("  Type: Retryable Error - The SDK will automatically retry")
		default:
			fmt.Println("  Type: Unknown Error")
		}

		if len(apiErr.Details) > 0 {
			fmt.Printf("  Details: %v\n", apiErr.Details)
		}
	}
}

// Helper functions for safe pointer dereferencing
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
