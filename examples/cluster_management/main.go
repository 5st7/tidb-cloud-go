// Package main demonstrates cluster management operations using the TiDB Cloud Go SDK.
// This example shows how to create, update, and delete clusters with proper error handling.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/5st7/tidb-cloud-go/pkg/client"
	"github.com/5st7/tidb-cloud-go/pkg/models"
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

	// Get the first project for demonstration
	projects, err := client.ListProjects()
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	if len(projects.Items) == 0 {
		log.Fatal("No projects found. Please create a project first.")
	}

	projectID := *projects.Items[0].ID
	fmt.Printf("Using project: %s (%s)\n", *projects.Items[0].Name, projectID)

	// Example 1: Create a new cluster
	fmt.Println("\n=== Creating a New Cluster ===")
	
	createReq := &models.OpenapiCreateClusterReq{
		Name:          stringPtr("sdk-demo-cluster"),
		ClusterType:   stringPtr("DEDICATED"),
		CloudProvider: stringPtr("AWS"),
		Region:        stringPtr("us-west-2"),
		Config: &models.OpenapiClusterConfig{
			RootPassword: stringPtr("YourSecurePassword123!"),
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

	cluster, err := client.CreateCluster(projectID, createReq)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}

	clusterID := *cluster.ClusterID
	fmt.Printf("Cluster created successfully! ID: %s\n", clusterID)
	fmt.Printf("Cluster status: %s\n", *cluster.Status.ClusterStatus)

	// Example 2: Wait for cluster to be available (in a real scenario, you'd want to poll)
	fmt.Println("\n=== Monitoring Cluster Status ===")
	for i := 0; i < 10; i++ {
		time.Sleep(30 * time.Second)
		
		clusterInfo, err := client.GetCluster(projectID, clusterID)
		if err != nil {
			log.Printf("Failed to get cluster info: %v", err)
			continue
		}

		status := *clusterInfo.Status.ClusterStatus
		fmt.Printf("Cluster status: %s\n", status)
		
		if status == "AVAILABLE" {
			fmt.Println("Cluster is now available!")
			break
		}
		
		if status == "UNAVAILABLE" || status == "FAILED" {
			fmt.Printf("Cluster creation failed with status: %s\n", status)
			return
		}
		
		fmt.Printf("Waiting for cluster to be ready... (attempt %d/10)\n", i+1)
	}

	// Example 3: Update cluster configuration
	fmt.Println("\n=== Updating Cluster Configuration ===")
	
	updateReq := &models.OpenapiUpdateClusterReq{
		Config: &models.OpenapiClusterConfig{
			Components: &models.OpenapiClusterConfigComponents{
				TiDB: &models.OpenapiUpdateTiDBComponent{
					NodeSize:     stringPtr("8C16G"),
					NodeQuantity: int64Ptr(2), // Scale up TiDB nodes
				},
			},
		},
	}

	updatedCluster, err := client.UpdateCluster(projectID, clusterID, updateReq)
	if err != nil {
		log.Printf("Failed to update cluster: %v", err)
	} else {
		fmt.Printf("Cluster update initiated. Status: %s\n", *updatedCluster.Status.ClusterStatus)
	}

	// Example 4: Create a backup
	fmt.Println("\n=== Creating a Backup ===")
	
	backupReq := &models.OpenapiCreateBackupReq{
		Name:        stringPtr("sdk-demo-backup"),
		Description: stringPtr("Backup created by SDK demo"),
	}

	backup, err := client.CreateBackup(projectID, clusterID, backupReq)
	if err != nil {
		log.Printf("Failed to create backup: %v", err)
	} else {
		fmt.Printf("Backup created successfully! ID: %s\n", *backup.BackupID)
	}

	// Example 5: List cluster backups
	fmt.Println("\n=== Listing Cluster Backups ===")
	
	backups, err := client.ListBackups(projectID, clusterID)
	if err != nil {
		log.Printf("Failed to list backups: %v", err)
	} else {
		fmt.Printf("Found %d backups:\n", len(backups.Items))
		for _, b := range backups.Items {
			fmt.Printf("- ID: %s, Name: %s, Status: %s\n",
				safeString(b.ID),
				safeString(b.Name),
				safeString(b.Status.BackupStatus))
		}
	}

	// Example 6: Set up private endpoint (if supported)
	fmt.Println("\n=== Setting up Private Endpoint ===")
	
	ctx := context.Background()
	
	// First, create the private endpoint service
	service, err := client.CreatePrivateEndpointService(ctx, projectID, clusterID)
	if err != nil {
		log.Printf("Failed to create private endpoint service: %v", err)
	} else {
		fmt.Printf("Private endpoint service created. Status: %s\n", *service.Status)
		fmt.Printf("Service name: %s\n", *service.Name)
		
		// In a real scenario, you would create the VPC endpoint in your cloud provider
		// and then create the private endpoint connection
		fmt.Println("Next steps:")
		fmt.Println("1. Create a VPC endpoint in your AWS/GCP console")
		fmt.Printf("2. Use service name: %s\n", *service.Name)
		fmt.Println("3. Call CreatePrivateEndpoint with your endpoint ID")
	}

	// Cleanup example (commented out for safety)
	// fmt.Println("\n=== Cleanup (commented out for safety) ===")
	// 
	// // Delete the cluster
	// err = client.DeleteCluster(projectID, clusterID)
	// if err != nil {
	// 	log.Printf("Failed to delete cluster: %v", err)
	// } else {
	// 	fmt.Printf("Cluster deletion initiated\n")
	// }

	fmt.Println("\n=== Cluster management example completed ===")
	fmt.Printf("Cluster ID: %s (remember to clean up manually if needed)\n", clusterID)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}