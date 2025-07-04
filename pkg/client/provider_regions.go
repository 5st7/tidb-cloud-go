package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/5st7/tidb-cloud-go/pkg/models"
)

// ListProviderRegions lists all available cloud providers, regions and specifications
func (c *Client) ListProviderRegions() (*models.OpenapiListProviderRegionsResp, error) {
	url := fmt.Sprintf("%s/api/%s/clusters/provider/regions", c.baseURL, APIVersion)

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

	var regions models.OpenapiListProviderRegionsResp
	if err := json.NewDecoder(resp.Body).Decode(&regions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &regions, nil
}
