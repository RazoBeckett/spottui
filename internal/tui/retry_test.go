package tui

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"rate limit", errors.New("rate limit exceeded"), true},
		{"429 status", errors.New("HTTP 429 Too Many Requests"), true},
		{"502 bad gateway", errors.New("502 Bad Gateway"), true},
		{"503 unavailable", errors.New("503 Service Unavailable"), true},
		{"504 timeout", errors.New("504 Gateway Timeout"), true},
		{"timeout", errors.New("request timeout"), true},
		{"connection refused", errors.New("connection refused"), true},
		{"connection reset", errors.New("connection reset by peer"), true},
		{"no such host", errors.New("no such host"), true},
		{"temporary failure", errors.New("temporary failure in name resolution"), true},
		{"auth error", errors.New("401 Unauthorized"), false},
		{"not found", errors.New("404 Not Found"), false},
		{"generic error", errors.New("something went wrong"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFriendlyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{"nil error", nil, ""},
		{"rate limit", errors.New("429 rate limit"), "too many requests"},
		{"unauthorized", errors.New("401 unauthorized"), "session expired"},
		{"forbidden", errors.New("403 forbidden"), "access denied"},
		{"not found", errors.New("404 not found"), "not found"},
		{"timeout", errors.New("request timeout"), "timed out"},
		{"no host", errors.New("no such host"), "cannot connect"},
		{"connection refused", errors.New("connection refused"), "cannot connect"},
		{"no device", errors.New("no active device"), "no Spotify device"},
		{"premium", errors.New("premium required"), "Premium required"},
		{"generic", errors.New("unknown error"), "unknown error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := friendlyError(tt.err)
			if tt.err == nil {
				assert.Nil(t, result)
			} else {
				assert.Contains(t, result.Error(), tt.contains)
			}
		})
	}
}

func TestWithRetry_Success(t *testing.T) {
	ctx := context.Background()
	callCount := 0

	result, err := withRetry(ctx, func() (string, error) {
		callCount++
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 1, callCount)
}

func TestWithRetry_ImmediateFailure(t *testing.T) {
	ctx := context.Background()
	callCount := 0

	_, err := withRetry(ctx, func() (string, error) {
		callCount++
		return "", errors.New("401 unauthorized")
	})

	assert.Error(t, err)
	assert.Equal(t, 1, callCount)
}

func TestWithRetry_RetryThenSuccess(t *testing.T) {
	ctx := context.Background()
	callCount := 0

	result, err := withRetry(ctx, func() (string, error) {
		callCount++
		if callCount < 3 {
			return "", errors.New("503 service unavailable")
		}
		return "success after retry", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success after retry", result)
	assert.Equal(t, 3, callCount)
}

func TestWithRetry_MaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	callCount := 0

	_, err := withRetry(ctx, func() (string, error) {
		callCount++
		return "", errors.New("503 service unavailable")
	})

	assert.Error(t, err)
	assert.Equal(t, maxRetries+1, callCount)
}

func TestWithRetry_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	callCount := 0

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := withRetry(ctx, func() (string, error) {
		callCount++
		return "", errors.New("503 service unavailable")
	})

	assert.Error(t, err)
	assert.True(t, callCount >= 1 && callCount <= maxRetries+1)
}
