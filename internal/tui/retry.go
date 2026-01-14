package tui

import (
	"context"
	"errors"
	"math"
	"net"
	"strings"
	"time"
)

const (
	maxRetries    = 3
	baseDelay     = 500 * time.Millisecond
	maxDelay      = 5 * time.Second
	backoffFactor = 2.0
)

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	if strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "502") ||
		strings.Contains(errStr, "503") ||
		strings.Contains(errStr, "504") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "temporary failure") {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	return false
}

func withRetry[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	var result T
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(float64(baseDelay) * math.Pow(backoffFactor, float64(attempt-1)))
			if delay > maxDelay {
				delay = maxDelay
			}

			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(delay):
			}
		}

		result, lastErr = fn()
		if lastErr == nil {
			return result, nil
		}

		if !isRetryableError(lastErr) {
			return result, lastErr
		}
	}

	return result, lastErr
}

func friendlyError(err error) error {
	if err == nil {
		return nil
	}

	errStr := strings.ToLower(err.Error())

	switch {
	case strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "429"):
		return errors.New("too many requests - please wait a moment")
	case strings.Contains(errStr, "401") || strings.Contains(errStr, "unauthorized"):
		return errors.New("session expired - please restart the app")
	case strings.Contains(errStr, "403") || strings.Contains(errStr, "forbidden"):
		return errors.New("access denied - check your Spotify Premium status")
	case strings.Contains(errStr, "404") || strings.Contains(errStr, "not found"):
		return errors.New("content not found")
	case strings.Contains(errStr, "timeout"):
		return errors.New("request timed out - check your connection")
	case strings.Contains(errStr, "no such host") || strings.Contains(errStr, "connection refused"):
		return errors.New("cannot connect - check your internet")
	case strings.Contains(errStr, "no active device") || strings.Contains(errStr, "no devices"):
		return errors.New("no Spotify device active - open Spotify on a device")
	case strings.Contains(errStr, "premium"):
		return errors.New("Spotify Premium required for this feature")
	default:
		return err
	}
}
