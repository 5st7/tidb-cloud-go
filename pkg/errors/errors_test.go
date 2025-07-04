package errors

import (
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		apiError APIError
		expected string
	}{
		{
			name: "with all fields",
			apiError: APIError{
				StatusCode: 400,
				Code:       49900001,
				Message:    "Invalid request",
				Details:    []interface{}{"field validation failed"},
			},
			expected: "TiDB Cloud API error (400): Invalid request (code: 49900001)",
		},
		{
			name: "without details",
			apiError: APIError{
				StatusCode: 404,
				Code:       49900002,
				Message:    "Resource not found",
			},
			expected: "TiDB Cloud API error (404): Resource not found (code: 49900002)",
		},
		{
			name: "rate limit error",
			apiError: APIError{
				StatusCode: 429,
				Code:       49900007,
				Message:    "Rate limit exceeded",
			},
			expected: "TiDB Cloud API error (429): Rate limit exceeded (code: 49900007)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.apiError.Error()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAPIError_IsRateLimitError(t *testing.T) {
	tests := []struct {
		name     string
		apiError APIError
		expected bool
	}{
		{
			name: "rate limit error",
			apiError: APIError{
				StatusCode: 429,
				Code:       49900007,
			},
			expected: true,
		},
		{
			name: "not rate limit error - wrong status",
			apiError: APIError{
				StatusCode: 400,
				Code:       49900007,
			},
			expected: false,
		},
		{
			name: "not rate limit error - wrong code",
			apiError: APIError{
				StatusCode: 429,
				Code:       49900001,
			},
			expected: false,
		},
		{
			name: "not rate limit error",
			apiError: APIError{
				StatusCode: 404,
				Code:       49900002,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.apiError.IsRateLimitError()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAPIError_IsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		apiError APIError
		expected bool
	}{
		{
			name:     "rate limit error - retryable",
			apiError: APIError{StatusCode: 429},
			expected: true,
		},
		{
			name:     "internal server error - retryable",
			apiError: APIError{StatusCode: 500},
			expected: true,
		},
		{
			name:     "bad gateway - retryable",
			apiError: APIError{StatusCode: 502},
			expected: true,
		},
		{
			name:     "service unavailable - retryable",
			apiError: APIError{StatusCode: 503},
			expected: true,
		},
		{
			name:     "gateway timeout - retryable",
			apiError: APIError{StatusCode: 504},
			expected: true,
		},
		{
			name:     "bad request - not retryable",
			apiError: APIError{StatusCode: 400},
			expected: false,
		},
		{
			name:     "unauthorized - not retryable",
			apiError: APIError{StatusCode: 401},
			expected: false,
		},
		{
			name:     "forbidden - not retryable",
			apiError: APIError{StatusCode: 403},
			expected: false,
		},
		{
			name:     "not found - not retryable",
			apiError: APIError{StatusCode: 404},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.apiError.IsRetryable()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
