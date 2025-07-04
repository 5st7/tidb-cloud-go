package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

// ListBackups lists all backups for a cluster
func (c *Client) ListBackups(projectID, clusterID string) (*models.OpenapiListBackupOfClusterResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/backups", c.baseURL, APIVersion, projectID, clusterID)
	
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

	var backups models.OpenapiListBackupOfClusterResp
	if err := json.NewDecoder(resp.Body).Decode(&backups); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &backups, nil
}

// GetBackup gets a backup by ID
func (c *Client) GetBackup(projectID, clusterID, backupID string) (*models.OpenapiGetBackupOfClusterResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is required")
	}
	if backupID == "" {
		return nil, fmt.Errorf("backup ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/backups/%s", c.baseURL, APIVersion, projectID, clusterID, backupID)
	
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

	var backup models.OpenapiGetBackupOfClusterResp
	if err := json.NewDecoder(resp.Body).Decode(&backup); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &backup, nil
}

// CreateBackup creates a new backup
func (c *Client) CreateBackup(projectID, clusterID string, req *models.OpenapiCreateBackupReq) (*models.OpenapiCreateBackupResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/backups", c.baseURL, APIVersion, projectID, clusterID)
	
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

	var createResp models.OpenapiCreateBackupResp
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createResp, nil
}

// DeleteBackup deletes a backup
func (c *Client) DeleteBackup(projectID, clusterID, backupID string) error {
	if projectID == "" {
		return fmt.Errorf("project ID is required")
	}
	if clusterID == "" {
		return fmt.Errorf("cluster ID is required")
	}
	if backupID == "" {
		return fmt.Errorf("backup ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/clusters/%s/backups/%s", c.baseURL, APIVersion, projectID, clusterID, backupID)
	
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