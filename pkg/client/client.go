// Package client provides a Go SDK for the TiDB Cloud API.
// It supports HTTP Digest Authentication, automatic retries with exponential backoff,
// and comprehensive error handling for all TiDB Cloud operations.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/5st7/tidb-cloud-go/pkg/auth"
	"github.com/5st7/tidb-cloud-go/pkg/errors"
	"github.com/5st7/tidb-cloud-go/pkg/models"
	"github.com/5st7/tidb-cloud-go/pkg/retry"
)

const (
	// DefaultBaseURL is the default TiDB Cloud API base URL
	DefaultBaseURL = "https://api.tidbcloud.com"
	// APIVersion is the current API version used by the client
	APIVersion = "v1beta"
)

// Client represents a TiDB Cloud API client.
// It handles authentication, retries, and error handling for all API operations.
type Client struct {
	baseURL       string
	httpClient    *http.Client
	publicKey     string
	privateKey    string
	digestAuth    *auth.DigestAuth
	retryExecutor *retry.RetryExecutor
}

// NewClient creates a new TiDB Cloud API client with the provided credentials.
// The client is configured with default settings including a 30-second timeout,
// automatic retry with exponential backoff, and HTTP Digest Authentication.
//
// Parameters:
//   - publicKey: Your TiDB Cloud API public key
//   - privateKey: Your TiDB Cloud API private key
//
// Returns:
//   - *Client: A configured TiDB Cloud client
//   - error: An error if the credentials are invalid
func NewClient(publicKey, privateKey string) (*Client, error) {
	if publicKey == "" {
		return nil, fmt.Errorf("public key is required")
	}
	if privateKey == "" {
		return nil, fmt.Errorf("private key is required")
	}

	retryPolicy := retry.NewRetryPolicy()
	retryExecutor := retry.NewRetryExecutor(retryPolicy)

	return &Client{
		baseURL:       DefaultBaseURL,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		publicKey:     publicKey,
		privateKey:    privateKey,
		digestAuth:    auth.NewDigestAuth(),
		retryExecutor: retryExecutor,
	}, nil
}

// ListProjects retrieves a list of all projects in your organization.
// Each project contains clusters, users, and other resources.
//
// Returns:
//   - *models.OpenapiListProjectsResp: A list of projects with their details
//   - error: An error if the request fails
func (c *Client) ListProjects() (*models.OpenapiListProjectsResp, error) {
	url := fmt.Sprintf("%s/api/%s/projects", c.baseURL, APIVersion)

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

	var projects models.OpenapiListProjectsResp
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &projects, nil
}

// CreateProject creates a new project in your organization.
// A project is a logical container for clusters and other resources.
//
// Parameters:
//   - req: The project creation request containing the project name
//
// Returns:
//   - *models.OpenapiCreateProjectResp: The created project details
//   - error: An error if the request fails or validation fails
func (c *Client) CreateProject(req *models.OpenapiCreateProjectReq) (*models.OpenapiCreateProjectResp, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	url := fmt.Sprintf("%s/api/%s/projects", c.baseURL, APIVersion)

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

	var createResp models.OpenapiCreateProjectResp
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createResp, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	return c.doRequestWithRetry(context.Background(), req)
}

func (c *Client) doRequestWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Store request body for potential retry
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body.Close()
	}

	var finalResp *http.Response
	var finalErr error

	operation := func() error {
		// Restore request body for each attempt
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		resp, err := c.executeHTTPRequest(req)
		if err != nil {
			finalErr = err
			return err
		}

		// Check for API errors
		if resp.StatusCode >= 400 {
			apiErr := c.parseAPIError(resp)
			finalErr = apiErr
			resp.Body.Close()
			return apiErr
		}

		finalResp = resp
		return nil
	}

	err := c.retryExecutor.Execute(ctx, operation)
	if err != nil {
		return nil, finalErr
	}

	return finalResp, nil
}

func (c *Client) executeHTTPRequest(req *http.Request) (*http.Response, error) {
	// Store request body before making the request
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// First attempt without auth
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// If we get 401, try with digest auth
	if resp.StatusCode == http.StatusUnauthorized {
		authHeader := resp.Header.Get("WWW-Authenticate")
		if authHeader != "" {
			resp.Body.Close()

			// Parse the challenge
			if err := c.digestAuth.ParseChallenge(authHeader); err != nil {
				return nil, fmt.Errorf("failed to parse auth challenge: %w", err)
			}

			// Create new request with auth header and restored body
			var newBody io.ReadCloser
			if bodyBytes != nil {
				newBody = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			newReq, err := http.NewRequest(req.Method, req.URL.String(), newBody)
			if err != nil {
				return nil, fmt.Errorf("failed to create auth request: %w", err)
			}

			// Copy headers from original request
			for k, v := range req.Header {
				newReq.Header[k] = v
			}

			// Add digest auth header
			authValue := c.digestAuth.GenerateAuthHeader(c.publicKey, c.privateKey, req.Method, req.URL.Path)
			newReq.Header.Set("Authorization", authValue)

			// Retry the request
			return c.httpClient.Do(newReq)
		}
	}

	return resp, nil
}

func (c *Client) parseAPIError(resp *http.Response) errors.APIError {
	apiError := errors.APIError{
		StatusCode: resp.StatusCode,
	}

	// Try to decode error response
	var errorResp models.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
		if errorResp.Code != nil {
			apiError.Code = *errorResp.Code
		}
		if errorResp.Message != nil {
			apiError.Message = *errorResp.Message
		}
		apiError.Details = errorResp.Details
	} else {
		// Fallback to HTTP status text
		apiError.Message = http.StatusText(resp.StatusCode)
	}

	// Add rate limit information if available
	if resp.StatusCode == http.StatusTooManyRequests {
		if resetHeader := resp.Header.Get("X-Ratelimit-Reset"); resetHeader != "" {
			if resetTime, err := strconv.ParseInt(resetHeader, 10, 64); err == nil {
				apiError.Details = append(apiError.Details, map[string]interface{}{
					"rate_limit_reset": resetTime,
				})
			}
		}
	}

	return apiError
}
