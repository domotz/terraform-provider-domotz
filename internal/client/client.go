package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout  = 30 * time.Second
	defaultPageSize = 100
	maxRetries      = 3
	Version         = "1.0.0"
)

// NotFoundError represents a 404 API response
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string { return e.Message }

// Client represents the Domotz API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new Domotz API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// doRequest executes an HTTP request with authentication, error handling, and retries
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second // 1s, 2s, 4s
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := c.doRequestOnce(ctx, method, path, body, result)
		if err == nil {
			return nil
		}

		if !isRetryableError(err) {
			return err
		}
		lastErr = err
	}

	return fmt.Errorf("max retries exceeded for %s %s: %w", method, path, lastErr)
}

// doRequestOnce executes a single HTTP request
func (c *Client) doRequestOnce(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("terraform-provider-domotz/%s", Version))

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle 404 Not Found specifically
	if resp.StatusCode == 404 {
		return &NotFoundError{
			Message: fmt.Sprintf("resource not found: %s %s", method, path),
		}
	}

	// Handle other non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Message)
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// isRetryableError checks if an error is retryable (rate limiting or transient errors)
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "status 429") ||
		strings.Contains(errMsg, "status 502") ||
		strings.Contains(errMsg, "status 503") ||
		strings.Contains(errMsg, "status 504")
}

// doRequestNoContent executes a request that expects no response body (e.g., DELETE)
func (c *Client) doRequestNoContent(ctx context.Context, method, path string, body interface{}) error {
	return c.doRequest(ctx, method, path, body, nil)
}
