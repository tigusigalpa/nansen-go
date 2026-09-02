package nansen

import (
	"errors"
	"net/http"
	"time"
)

// Option configures a Client during construction.
type Option func(*Client) error

// WithBaseURL overrides the Nansen API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) error {
		baseURL, err := normalizeBaseURL(url)
		if err != nil {
			return err
		}
		c.baseURL = baseURL
		return nil
	}
}

// WithHTTPClient uses client to make HTTP requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		if client == nil {
			return errors.New("nansen: HTTP client cannot be nil")
		}
		c.httpClient = client
		return nil
	}
}

// WithTimeout sets the default timeout for requests without a deadline.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		if timeout <= 0 {
			return errors.New("nansen: timeout must be positive")
		}
		c.timeout = timeout
		return nil
	}
}

// WithRetry enables retries with exponential backoff for transient failures.
func WithRetry(maxAttempts int, initialDelay, maxDelay time.Duration) Option {
	return func(c *Client) error {
		if maxAttempts < 0 {
			return errors.New("nansen: maxAttempts cannot be negative")
		}
		if maxAttempts > 0 {
			if initialDelay <= 0 {
				return errors.New("nansen: initial delay must be positive")
			}
			if maxDelay <= 0 {
				return errors.New("nansen: max delay must be positive")
			}
		}
		c.retry = retryConfig{
			maxAttempts:    maxAttempts,
			initialDelay:   initialDelay,
			maxDelay:       maxDelay,
			retryRateLimit: true,
		}
		return nil
	}
}
