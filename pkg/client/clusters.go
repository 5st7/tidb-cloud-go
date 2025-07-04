package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

// ListClusters lists all clusters in a project
func (c *Client) ListClusters(projectID string) (*models.OpenapiListClustersOfProjectResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters", c.baseURL, APIVersion, projectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var clusters models.OpenapiListClustersOfProjectResp
	if err := json.NewDecoder(resp.Body).Decode(&clusters); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &clusters, nil
}

// GetCluster gets a cluster by ID
func (c *Client) GetCluster(projectID, clusterID string) (*models.OpenapiClusterItem, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s", c.baseURL, APIVersion, projectID, clusterID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var cluster models.OpenapiClusterItem
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &cluster, nil
}

// CreateCluster creates a new cluster
func (c *Client) CreateCluster(projectID string, req *models.OpenapiCreateClusterReq) (*models.OpenapiCreateClusterResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters", c.baseURL, APIVersion, projectID)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var createResp models.OpenapiCreateClusterResp
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createResp, nil
}

// UpdateCluster updates an existing cluster
func (c *Client) UpdateCluster(projectID, clusterID string, req *models.OpenapiUpdateClusterReq) error {
	if projectID == "" {
		return fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return fmt.Errorf("cluster ID is required")
	}
	if req == nil {
		return fmt.Errorf("request is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s", c.baseURL, APIVersion, projectID, clusterID)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

// DeleteCluster deletes a cluster
func (c *Client) DeleteCluster(projectID, clusterID string) error {
	if projectID == "" {
		return fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return fmt.Errorf("cluster ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s", c.baseURL, APIVersion, projectID, clusterID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}
