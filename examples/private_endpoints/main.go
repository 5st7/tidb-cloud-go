// Package main demonstrates private endpoint management using the TiDB Cloud Go SDK.
// This example shows how to set up and manage private network connections to TiDB Cloud clusters.
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

	// Set up context with timeout for API calls
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
		log.Fatal("No clusters found. Please create a Dedicated cluster first.")
	}

	// Find a Dedicated cluster (private endpoints are only supported for Dedicated clusters)
	var clusterID, clusterName string
	for _, cluster := range clusters.Items {
		if cluster.ClusterType != nil && *cluster.ClusterType == "DEDICATED" {
			clusterID = *cluster.ID
			clusterName = *cluster.Name
			break
		}
	}

	if clusterID == "" {
		log.Fatal("No Dedicated clusters found. Private endpoints are only supported for Dedicated clusters.")
	}

	fmt.Printf("Using cluster: %s (%s)\n", clusterName, clusterID)

	// Example 1: List all private endpoints in the project
	fmt.Println("\n=== Listing All Private Endpoints in Project ===")
	
	allEndpoints, err := client.ListPrivateEndpointsOfProject(ctx, projectID)
	if err != nil {
		log.Printf("Failed to list project private endpoints: %v", err)
	} else {
		fmt.Printf("Found %d private endpoints in the project:\n", len(allEndpoints.Items))
		for _, endpoint := range allEndpoints.Items {
			fmt.Printf("- ID: %s\n", safeString(endpoint.ID))
			fmt.Printf("  Cluster: %s\n", safeString(endpoint.ClusterID))
			fmt.Printf("  Provider: %s\n", safeString(endpoint.CloudProvider))
			fmt.Printf("  Region: %s\n", safeString(endpoint.Region))
			fmt.Printf("  Endpoint Name: %s\n", safeString(endpoint.EndpointName))
			fmt.Printf("  Status: %s\n", safeString(endpoint.Status))
			fmt.Printf("  Service Name: %s\n", safeString(endpoint.ServiceName))
			fmt.Println()
		}
	}

	// Example 2: Check if private endpoint service exists for the cluster
	fmt.Println("=== Checking Private Endpoint Service ===")
	
	service, err := client.GetPrivateEndpointService(ctx, projectID, clusterID)
	if err != nil {
		fmt.Printf("No private endpoint service found, creating one...\n")
		
		// Example 3: Create private endpoint service
		fmt.Println("\n=== Creating Private Endpoint Service ===")
		
		service, err = client.CreatePrivateEndpointService(ctx, projectID, clusterID)
		if err != nil {
			log.Fatalf("Failed to create private endpoint service: %v", err)
		}
		
		fmt.Println("Private endpoint service created successfully!")
	} else {
		fmt.Println("Private endpoint service already exists!")
	}

	// Display service information
	fmt.Printf("Service Details:\n")
	fmt.Printf("  Cloud Provider: %s\n", safeString(service.CloudProvider))
	fmt.Printf("  Service Name: %s\n", safeString(service.Name))
	fmt.Printf("  Status: %s\n", safeString(service.Status))
	fmt.Printf("  DNS Name: %s\n", safeString(service.DNSName))
	fmt.Printf("  Port: %d\n", safeInt64(service.Port))
	if len(service.AzIDs) > 0 {
		fmt.Printf("  Availability Zones: %v\n", service.AzIDs)
	}

	// Example 4: List existing private endpoints for the cluster
	fmt.Println("\n=== Listing Private Endpoints for Cluster ===")
	
	endpoints, err := client.ListPrivateEndpoints(ctx, projectID, clusterID)
	if err != nil {
		log.Printf("Failed to list private endpoints: %v", err)
	} else {
		fmt.Printf("Found %d private endpoints for this cluster:\n", len(endpoints.Items))
		for _, endpoint := range endpoints.Items {
			fmt.Printf("- ID: %s\n", safeString(endpoint.ID))
			fmt.Printf("  Endpoint Name: %s\n", safeString(endpoint.EndpointName))
			fmt.Printf("  Status: %s\n", safeString(endpoint.Status))
			fmt.Printf("  Message: %s\n", safeString(endpoint.Message))
			fmt.Println()
		}
	}

	// Example 5: Instructions for creating VPC endpoint (this requires manual steps in cloud console)
	fmt.Println("\n=== Instructions for Creating VPC Endpoint ===")
	fmt.Println("To create a private endpoint connection, you need to:")
	fmt.Println("1. Go to your cloud provider console (AWS/GCP)")
	fmt.Println("2. Create a VPC endpoint with the following details:")
	fmt.Printf("   - Service Name: %s\n", safeString(service.Name))
	if service.CloudProvider != nil {
		switch *service.CloudProvider {
		case "AWS":
			fmt.Println("   - Service Type: Interface")
			fmt.Println("   - Policy: Full Access (or custom policy)")
			fmt.Println("   - VPC: Your target VPC")
			fmt.Println("   - Subnets: Select appropriate subnets")
			fmt.Println("   - Security Groups: Configure as needed")
		case "GCP":
			fmt.Println("   - Network: Your target VPC network")
			fmt.Println("   - Subnetwork: Select appropriate subnetwork")
		}
	}
	fmt.Println("3. Copy the VPC endpoint ID (e.g., vpce-xxxxxxxxx for AWS)")
	fmt.Println("4. Use CreatePrivateEndpoint API with the endpoint ID")

	// Example 6: Simulate creating a private endpoint (with placeholder endpoint ID)
	// NOTE: In real usage, you would get this ID after creating the VPC endpoint in your cloud console
	if len(os.Args) > 1 && os.Args[1] == "create-endpoint" {
		endpointName := os.Getenv("VPC_ENDPOINT_ID") // e.g., "vpce-1234567890abcdef0"
		if endpointName == "" {
			fmt.Println("\nTo create a private endpoint, set VPC_ENDPOINT_ID environment variable")
			fmt.Println("and run with 'create-endpoint' argument")
		} else {
			fmt.Printf("\n=== Creating Private Endpoint with ID: %s ===\n", endpointName)
			
			createReq := &models.OpenapiCreatePrivateEndpointReq{
				EndpointName: stringPtr(endpointName),
			}

			endpoint, err := client.CreatePrivateEndpoint(ctx, projectID, clusterID, createReq)
			if err != nil {
				log.Printf("Failed to create private endpoint: %v", err)
			} else {
				fmt.Printf("Private endpoint created successfully!\n")
				fmt.Printf("  ID: %s\n", safeString(endpoint.ID))
				fmt.Printf("  Status: %s\n", safeString(endpoint.Status))
				fmt.Printf("  Message: %s\n", safeString(endpoint.Message))
				
				// Monitor connection status
				fmt.Println("\n=== Monitoring Private Endpoint Status ===")
				for i := 0; i < 10; i++ {
					time.Sleep(30 * time.Second)
					
					endpoints, err := client.ListPrivateEndpoints(ctx, projectID, clusterID)
					if err != nil {
						log.Printf("Failed to check endpoint status: %v", err)
						continue
					}
					
					for _, ep := range endpoints.Items {
						if ep.ID != nil && *ep.ID == *endpoint.ID {
							status := safeString(ep.Status)
							fmt.Printf("Endpoint status: %s\n", status)
							
							if status == "ACTIVE" {
								fmt.Println("Private endpoint is now active!")
								return
							}
							
							if status == "FAILED" {
								fmt.Printf("Private endpoint failed: %s\n", safeString(ep.Message))
								return
							}
						}
					}
					
					fmt.Printf("Waiting for endpoint to be active... (attempt %d/10)\n", i+1)
				}
			}
		}
	}

	// Example 7: Display connection information
	fmt.Println("\n=== Connection Information ===")
	fmt.Println("Once your private endpoint is active, you can connect using:")
	if service.DNSName != nil && service.Port != nil {
		fmt.Printf("Host: %s\n", *service.DNSName)
		fmt.Printf("Port: %d\n", *service.Port)
	}
	fmt.Println("Username: root (or your database user)")
	fmt.Println("Password: <your cluster password>")
	fmt.Println("SSL: Required")

	fmt.Println("\n=== Private Endpoint example completed ===")
	fmt.Println("Run with 'create-endpoint' argument and VPC_ENDPOINT_ID env var to create an endpoint")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
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