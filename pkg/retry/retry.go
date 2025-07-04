// Package retry provides automatic retry functionality with exponential backoff
// for the TiDB Cloud SDK. It supports intelligent retry logic based on error types
// and implements context-aware cancellation.
package retry

import (
	"context"
	"math"
	"time"

	"github.com/5st7/tidb-cloud-go/pkg/errors"
)

// RetryPolicy defines the retry policy for API requests.
// It configures the maximum number of attempts, base delay, and maximum delay
// for exponential backoff retry logic.
type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// NewRetryPolicy creates a new retry policy with default values.
// Default configuration:
//   - MaxAttempts: 3 (initial attempt + 2 retries)
//   - BaseDelay: 1 second
//   - MaxDelay: 30 seconds
func NewRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

// ShouldRetry determines if an error should be retried based on the error type
// and current attempt count. It returns true for retryable errors like rate limits
// and server errors, but false for client errors like authentication failures.
func (p *RetryPolicy) ShouldRetry(err error, attempt int) bool {
	if attempt >= p.MaxAttempts {
		return false
	}

	// Check if it's an API error
	if apiErr, ok := err.(errors.APIError); ok {
		return apiErr.IsRetryable()
	}

	// Retry non-API errors (network errors, etc.)
	return true
}

// CalculateDelay calculates the delay for the given attempt using exponential backoff.
// The delay starts at BaseDelay and doubles with each attempt, capped at MaxDelay.
// Formula: min(BaseDelay * 2^(attempt-1), MaxDelay)
func (p *RetryPolicy) CalculateDelay(attempt int) time.Duration {
	delay := time.Duration(math.Pow(2, float64(attempt-1))) * p.BaseDelay
	if delay > p.MaxDelay {
		delay = p.MaxDelay
	}
	return delay
}

// RetryExecutor executes operations with retry logic.
// It applies the configured retry policy and handles context cancellation.
type RetryExecutor struct {
	policy *RetryPolicy
}

// NewRetryExecutor creates a new retry executor with the specified policy.
// The executor will apply the policy's retry logic to all operations.
func NewRetryExecutor(policy *RetryPolicy) *RetryExecutor {
	return &RetryExecutor{
		policy: policy,
	}
}

// Execute executes an operation with retry logic according to the configured policy.
// It respects context cancellation and applies exponential backoff between retries.
// The operation function is called repeatedly until it succeeds, fails with a
// non-retryable error, or the maximum attempts are reached.
func (e *RetryExecutor) Execute(ctx context.Context, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt <= e.policy.MaxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on the last attempt or if it's not retryable
		if !e.policy.ShouldRetry(err, attempt) {
			break
		}

		// Calculate and wait for delay
		delay := e.policy.CalculateDelay(attempt + 1)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return lastErr
}
