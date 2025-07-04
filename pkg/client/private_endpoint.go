package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

// GetPrivateEndpointService retrieves the private endpoint service information for a cluster.
// Private endpoint services enable secure private network access to TiDB Cloud clusters
// through AWS PrivateLink or Google Cloud Private Service Connect.
//
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - projectID: The ID of the project containing the cluster
//   - clusterID: The ID of the cluster
//
// Returns:
//   - *models.OpenapiGetPrivateEndpointServiceResp: The private endpoint service details
//   - error: An error if the request fails or parameters are invalid
func (c *Client) GetPrivateEndpointService(ctx context.Context, projectID, clusterID string) (*models.OpenapiGetPrivateEndpointServiceResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/private_endpoint_service", c.baseURL, APIVersion, projectID, clusterID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequestWithRetry(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var serviceResp models.OpenapiGetPrivateEndpointServiceResp
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &serviceResp, nil
}

// CreatePrivateEndpointService creates a private endpoint service for a cluster.
// This enables the cluster to accept private endpoint connections from your VPC.
// Note: This operation is only available for TiDB Cloud Dedicated clusters.
//
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - projectID: The ID of the project containing the cluster
//   - clusterID: The ID of the cluster
//
// Returns:
//   - *models.OpenapiGetPrivateEndpointServiceResp: The created service details
//   - error: An error if the request fails or parameters are invalid
func (c *Client) CreatePrivateEndpointService(ctx context.Context, projectID, clusterID string) (*models.OpenapiGetPrivateEndpointServiceResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/private_endpoint_service", c.baseURL, APIVersion, projectID, clusterID)
	
	// According to the API spec, the body should be an empty object
	reqBody := map[string]interface{}{}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequestWithRetry(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var serviceResp models.OpenapiGetPrivateEndpointServiceResp
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &serviceResp, nil
}

// ListPrivateEndpoints lists all private endpoints for a cluster.
// Private endpoints represent the connection points from your VPC to the TiDB Cloud cluster.
//
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - projectID: The ID of the project containing the cluster
//   - clusterID: The ID of the cluster
//
// Returns:
//   - *models.OpenapiListPrivateEndpointsResp: A list of private endpoints
//   - error: An error if the request fails or parameters are invalid
func (c *Client) ListPrivateEndpoints(ctx context.Context, projectID, clusterID string) (*models.OpenapiListPrivateEndpointsResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/private_endpoints", c.baseURL, APIVersion, projectID, clusterID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequestWithRetry(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var endpointsResp models.OpenapiListPrivateEndpointsResp
	if err := json.NewDecoder(resp.Body).Decode(&endpointsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &endpointsResp, nil
}

// CreatePrivateEndpoint creates a private endpoint for a cluster.
// This establishes a private network connection from your VPC to the TiDB Cloud cluster.
// The endpoint name format varies by cloud provider (e.g., 'vpce-xxxxxx' for AWS).
//
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - projectID: The ID of the project containing the cluster
//   - clusterID: The ID of the cluster
//   - req: The private endpoint creation request
//
// Returns:
//   - *models.OpenapiCreatePrivateEndpointResp: The created endpoint details
//   - error: An error if the request fails or parameters are invalid
func (c *Client) CreatePrivateEndpoint(ctx context.Context, projectID, clusterID string, req *models.OpenapiCreatePrivateEndpointReq) (*models.OpenapiCreatePrivateEndpointResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/private_endpoints", c.baseURL, APIVersion, projectID, clusterID)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequestWithRetry(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var createResp models.OpenapiCreatePrivateEndpointResp
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createResp, nil
}

// DeletePrivateEndpoint deletes a private endpoint for a cluster.
// This will disconnect the private network connection from your VPC to the cluster.
// The operation is irreversible and may cause connection interruption.
//
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - projectID: The ID of the project containing the cluster
//   - clusterID: The ID of the cluster
//   - endpointID: The ID of the private endpoint to delete
//
// Returns:
//   - error: An error if the request fails or parameters are invalid
func (c *Client) DeletePrivateEndpoint(ctx context.Context, projectID, clusterID, endpointID string) error {
	if projectID == "" {
		return fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return fmt.Errorf("cluster ID is required")
	}
	if endpointID == "" {
		return fmt.Errorf("endpoint ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/private_endpoints/%s", c.baseURL, APIVersion, projectID, clusterID, endpointID)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequestWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

// ListPrivateEndpointsOfProject lists all private endpoints in a project.
// This provides a project-wide view of all private endpoint connections
// across all clusters in the project.
//
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - projectID: The ID of the project
//
// Returns:
//   - *models.OpenapiListPrivateEndpointsResp: A list of all private endpoints in the project
//   - error: An error if the request fails or parameters are invalid
func (c *Client) ListPrivateEndpointsOfProject(ctx context.Context, projectID string) (*models.OpenapiListPrivateEndpointsResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/private_endpoints", c.baseURL, APIVersion, projectID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequestWithRetry(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var endpointsResp models.OpenapiListPrivateEndpointsResp
	if err := json.NewDecoder(resp.Body).Decode(&endpointsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &endpointsResp, nil
}