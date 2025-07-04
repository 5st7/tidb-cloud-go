// Package errors provides error types and utilities for the TiDB Cloud SDK.
// It includes comprehensive error handling for API responses with support
// for retryable error detection and specific error type checking.
package errors

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the TiDB Cloud API.
// It includes the HTTP status code, TiDB Cloud specific error code,
// error message, and optional additional details.
type APIError struct {
	StatusCode int           `json:"-"`
	Code       int64         `json:"code,omitempty"`
	Message    string        `json:"message,omitempty"`
	Details    []interface{} `json:"details,omitempty"`
}

// Error implements the error interface and returns a formatted error message
// that includes the HTTP status code, error message, and TiDB Cloud error code.
func (e APIError) Error() string {
	return fmt.Sprintf("TiDB Cloud API error (%d): %s (code: %d)", e.StatusCode, e.Message, e.Code)
}

// IsRateLimitError returns true if this is a rate limit error.
// TiDB Cloud enforces a rate limit of 100 requests per minute per API key.
func (e APIError) IsRateLimitError() bool {
	return e.StatusCode == http.StatusTooManyRequests && e.Code == 49900007
}

// IsRetryable returns true if this error should be retried.
// Retryable errors include rate limits, server errors, and temporary network issues.
func (e APIError) IsRetryable() bool {
	switch e.StatusCode {
	case http.StatusTooManyRequests,
		 http.StatusInternalServerError,
		 http.StatusBadGateway,
		 http.StatusServiceUnavailable,
		 http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// IsAuthenticationError returns true if this is an authentication error (401).
// This typically indicates invalid or missing API credentials.
func (e APIError) IsAuthenticationError() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsAuthorizationError returns true if this is an authorization error (403).
// This indicates the API key lacks permission for the requested operation.
func (e APIError) IsAuthorizationError() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsNotFoundError returns true if this is a not found error (404).
// This indicates the requested resource does not exist or is not accessible.
func (e APIError) IsNotFoundError() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsBadRequestError returns true if this is a bad request error (400).
// This indicates invalid request parameters or malformed request data.
func (e APIError) IsBadRequestError() bool {
	return e.StatusCode == http.StatusBadRequest
}