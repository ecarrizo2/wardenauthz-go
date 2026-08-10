package wardenauth

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

const (
	defaultMaxRetries = 3
	defaultRetryDelay = 500 * time.Millisecond
	maxRetryDelay     = 30 * time.Second
)

func (c *Client) doWithRetry(ctx context.Context, method, path string, body any, result any) error {
	maxRetries := c.maxRetries
	retryDelay := c.retryDelay

	for attempt := 0; ; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("request cancelled: %w", ctx.Err())
		default:
		}

		err := c.doRaw(ctx, method, path, body, result)
		if err == nil {
			return nil
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			return err
		}

		if !isRetryable(apiErr.StatusCode) || attempt >= maxRetries {
			return err
		}

		delay := time.Duration(math.Pow(2, float64(attempt))) * retryDelay
		if delay > maxRetryDelay {
			delay = maxRetryDelay
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("request cancelled during retry backoff: %w", ctx.Err())
		case <-time.After(delay):
		}
	}
}

func isRetryable(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= 500
}
