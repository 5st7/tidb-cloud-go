package retry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/5st7/tidb-cloud-go/pkg/errors"
)

func TestRetryPolicy_ShouldRetry(t *testing.T) {
	policy := NewRetryPolicy()

	tests := []struct {
		name     string
		err      error
		attempt  int
		expected bool
	}{
		{
			name:     "rate limit error within max attempts",
			err:      errors.APIError{StatusCode: 429, Code: 49900007},
			attempt:  1,
			expected: true,
		},
		{
			name:     "rate limit error at max attempts",
			err:      errors.APIError{StatusCode: 429, Code: 49900007},
			attempt:  3,
			expected: false,
		},
		{
			name:     "server error within max attempts",
			err:      errors.APIError{StatusCode: 500},
			attempt:  1,
			expected: true,
		},
		{
			name:     "bad request error",
			err:      errors.APIError{StatusCode: 400},
			attempt:  1,
			expected: false,
		},
		{
			name:     "non-API error",
			err:      fmt.Errorf("network error"),
			attempt:  1,
			expected: true,
		},
		{
			name:     "non-API error at max attempts",
			err:      fmt.Errorf("network error"),
			attempt:  3,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := policy.ShouldRetry(tt.err, tt.attempt)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRetryPolicy_CalculateDelay(t *testing.T) {
	policy := NewRetryPolicy()

	tests := []struct {
		name     string
		attempt  int
		expected time.Duration
	}{
		{
			name:     "first retry",
			attempt:  1,
			expected: 1 * time.Second,
		},
		{
			name:     "second retry",
			attempt:  2,
			expected: 2 * time.Second,
		},
		{
			name:     "third retry",
			attempt:  3,
			expected: 4 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := policy.CalculateDelay(tt.attempt)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRetryExecutor_Execute(t *testing.T) {
	executor := NewRetryExecutor(NewRetryPolicy())

	t.Run("successful on first attempt", func(t *testing.T) {
		callCount := 0
		operation := func() error {
			callCount++
			return nil
		}

		err := executor.Execute(context.Background(), operation)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if callCount != 1 {
			t.Errorf("Expected 1 call, got %d", callCount)
		}
	})

	t.Run("successful on second attempt", func(t *testing.T) {
		callCount := 0
		operation := func() error {
			callCount++
			if callCount == 1 {
				return errors.APIError{StatusCode: 500}
			}
			return nil
		}

		err := executor.Execute(context.Background(), operation)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if callCount != 2 {
			t.Errorf("Expected 2 calls, got %d", callCount)
		}
	})

	t.Run("fails after max attempts", func(t *testing.T) {
		callCount := 0
		operation := func() error {
			callCount++
			return errors.APIError{StatusCode: 500}
		}

		err := executor.Execute(context.Background(), operation)
		if err == nil {
			t.Error("Expected error but got none")
		}
		if callCount != 4 { // initial + 3 retries
			t.Errorf("Expected 4 calls, got %d", callCount)
		}
	})

	t.Run("non-retryable error", func(t *testing.T) {
		callCount := 0
		operation := func() error {
			callCount++
			return errors.APIError{StatusCode: 400}
		}

		err := executor.Execute(context.Background(), operation)
		if err == nil {
			t.Error("Expected error but got none")
		}
		if callCount != 1 {
			t.Errorf("Expected 1 call, got %d", callCount)
		}
	})
}
