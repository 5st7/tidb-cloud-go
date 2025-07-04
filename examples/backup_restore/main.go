// Package main demonstrates backup and restore operations using the TiDB Cloud Go SDK.
// This example shows how to create backups, monitor their status, and restore clusters from backups.
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

	// Get the first project and cluster for demonstration
	projects, err := client.ListProjects()
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	if len(projects.Items) == 0 {
		log.Fatal("No projects found. Please create a project first.")
	}

	projectID := *projects.Items[0].ID
	fmt.Printf("Using project: %s (%s)\n", *projects.Items[0].Name, projectID)

	clusters, err := client.ListClusters(projectID)
	if err != nil {
		log.Fatalf("Failed to list clusters: %v", err)
	}

	if len(clusters.Items) == 0 {
		log.Fatal("No clusters found. Please create a cluster first.")
	}

	clusterID := *clusters.Items[0].ID
	clusterName := *clusters.Items[0].Name
	fmt.Printf("Using cluster: %s (%s)\n", clusterName, clusterID)

	// Example 1: List existing backups
	fmt.Println("\n=== Listing Existing Backups ===")

	backups, err := client.ListBackups(projectID, clusterID)
	if err != nil {
		log.Fatalf("Failed to list backups: %v", err)
	}

	fmt.Printf("Found %d existing backups:\n", len(backups.Items))
	for _, backup := range backups.Items {
		fmt.Printf("- ID: %s\n", safeString(backup.ID))
		fmt.Printf("  Name: %s\n", safeString(backup.Name))
		fmt.Printf("  Type: %s\n", safeString(backup.Type))
		fmt.Printf("  Status: %s\n", safeString(backup.Status.BackupStatus))
		fmt.Printf("  Size: %d bytes\n", safeInt64(backup.BackupSizeBytes))
		fmt.Printf("  Created: %s\n", safeString(backup.CreateTimestamp))
		if backup.ExpiryTime != nil {
			fmt.Printf("  Expires: %s\n", safeString(backup.ExpiryTime))
		}
		fmt.Println()
	}

	// Example 2: Create a new backup
	fmt.Println("=== Creating a New Backup ===")

	backupName := fmt.Sprintf("sdk-demo-backup-%d", time.Now().Unix())
	createBackupReq := &models.OpenapiCreateBackupReq{
		Name:        stringPtr(backupName),
		Description: stringPtr("Backup created by SDK demo for testing restore functionality"),
	}

	newBackup, err := client.CreateBackup(projectID, clusterID, createBackupReq)
	if err != nil {
		log.Fatalf("Failed to create backup: %v", err)
	}

	backupID := *newBackup.BackupID
	fmt.Printf("Backup created successfully! ID: %s\n", backupID)

	// Example 3: Monitor backup progress
	fmt.Println("\n=== Monitoring Backup Progress ===")

	for i := 0; i < 20; i++ {
		time.Sleep(30 * time.Second)

		backupInfo, err := client.GetBackup(projectID, clusterID, backupID)
		if err != nil {
			log.Printf("Failed to get backup info: %v", err)
			continue
		}

		status := safeString(backupInfo.Status.BackupStatus)
		fmt.Printf("Backup status: %s", status)

		if backupInfo.BackupSizeBytes != nil && *backupInfo.BackupSizeBytes > 0 {
			fmt.Printf(" (Size: %d bytes)", *backupInfo.BackupSizeBytes)
		}
		fmt.Println()

		if status == "SUCCESS" {
			fmt.Println("Backup completed successfully!")
			fmt.Printf("Final size: %d bytes\n", safeInt64(backupInfo.BackupSizeBytes))
			break
		}

		if status == "FAILED" {
			fmt.Printf("Backup failed!\n")
			return
		}

		fmt.Printf("Waiting for backup to complete... (attempt %d/20)\n", i+1)
	}

	// Example 4: List all restores in the project
	fmt.Println("\n=== Listing Existing Restores ===")

	restores, err := client.ListRestores(projectID)
	if err != nil {
		log.Printf("Failed to list restores: %v", err)
	} else {
		fmt.Printf("Found %d existing restores:\n", len(restores.Items))
		for _, restore := range restores.Items {
			fmt.Printf("- ID: %s\n", safeString(restore.ID))
			fmt.Printf("  Name: %s\n", safeString(restore.Name))
			fmt.Printf("  Backup ID: %s\n", safeString(restore.BackupID))
			fmt.Printf("  Status: %s\n", safeString(restore.Status.RestoreStatus))
			if restore.ClusterInfo != nil {
				fmt.Printf("  Cluster: %s (%s)\n",
					safeString(restore.ClusterInfo.Name),
					safeString(restore.ClusterInfo.ID))
			}
			fmt.Printf("  Created: %s\n", safeString(restore.CreateTimestamp))
			if restore.FinishedTimestamp != nil {
				fmt.Printf("  Finished: %s\n", safeString(restore.FinishedTimestamp))
			}
			fmt.Println()
		}
	}

	// Example 5: Create a restore (create a new cluster from backup)
	fmt.Println("=== Creating a Restore (New Cluster from Backup) ===")

	restoreName := fmt.Sprintf("sdk-demo-restore-%d", time.Now().Unix())
	createRestoreReq := &models.OpenapiCreateRestoreReq{
		BackupID: stringPtr(backupID),
		Name:     stringPtr(restoreName),
		Config: &models.OpenapiClusterConfig{
			RootPassword: stringPtr("YourSecurePassword123!"),
			Port:         int64Ptr(4000),
			Components: &models.OpenapiClusterConfigComponents{
				TiDB: &models.OpenapiUpdateTiDBComponent{
					NodeSize:     stringPtr("8C16G"),
					NodeQuantity: int64Ptr(1),
				},
				TiKV: &models.OpenapiUpdateTiKVComponent{
					NodeSize:       stringPtr("8C32G"),
					NodeQuantity:   int64Ptr(3),
					StorageSizeGib: int64Ptr(500),
				},
			},
		},
	}

	restore, err := client.CreateRestore(projectID, createRestoreReq)
	if err != nil {
		log.Printf("Failed to create restore: %v", err)
	} else {
		restoreID := *restore.RestoreID
		fmt.Printf("Restore created successfully! ID: %s\n", restoreID)

		// Example 6: Monitor restore progress
		fmt.Println("\n=== Monitoring Restore Progress ===")

		for i := 0; i < 20; i++ {
			time.Sleep(60 * time.Second) // Restores take longer than backups

			restoreInfo, err := client.GetRestore(projectID, restoreID)
			if err != nil {
				log.Printf("Failed to get restore info: %v", err)
				continue
			}

			status := safeString(restoreInfo.Status.RestoreStatus)
			fmt.Printf("Restore status: %s\n", status)

			if status == "SUCCESS" {
				fmt.Println("Restore completed successfully!")
				if restoreInfo.ClusterInfo != nil {
					fmt.Printf("New cluster created: %s (%s)\n",
						safeString(restoreInfo.ClusterInfo.Name),
						safeString(restoreInfo.ClusterInfo.ID))
				}
				break
			}

			if status == "FAILED" {
				fmt.Printf("Restore failed!\n")
				break
			}

			fmt.Printf("Waiting for restore to complete... (attempt %d/20)\n", i+1)
		}
	}

	// Example 7: Clean up - delete the backup we created
	fmt.Println("\n=== Cleanup: Deleting Test Backup ===")

	err = client.DeleteBackup(projectID, clusterID, backupID)
	if err != nil {
		log.Printf("Failed to delete backup: %v", err)
	} else {
		fmt.Printf("Backup %s deletion initiated\n", backupID)
	}

	fmt.Println("\n=== Backup and Restore example completed ===")
	fmt.Println("Remember to clean up any restored clusters manually if they're no longer needed")
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

func safeInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
