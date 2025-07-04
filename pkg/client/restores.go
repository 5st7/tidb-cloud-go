package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

// ListRestores lists all restore tasks in a project
func (c *Client) ListRestores(projectID string) (*models.OpenapiListRestoreOfProjectResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/restores", c.baseURL, APIVersion, projectID)
	
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

	var restores models.OpenapiListRestoreOfProjectResp
	if err := json.NewDecoder(resp.Body).Decode(&restores); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &restores, nil
}

// GetRestore gets a restore task by ID
func (c *Client) GetRestore(projectID, restoreID string) (*models.OpenapiGetRestoreResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if restoreID == "" {
		return nil, fmt.Errorf("restore ID is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/restores/%s", c.baseURL, APIVersion, projectID, restoreID)
	
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

	var restore models.OpenapiGetRestoreResp
	if err := json.NewDecoder(resp.Body).Decode(&restore); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &restore, nil
}

// CreateRestore creates a new restore task
func (c *Client) CreateRestore(projectID string, req *models.OpenapiCreateRestoreReq) (*models.OpenapiCreateRestoreResp, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects/%s/restores", c.baseURL, APIVersion, projectID)
	
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

	var createResp models.OpenapiCreateRestoreResp
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createResp, nil
}