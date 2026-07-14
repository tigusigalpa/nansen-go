package nansen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.nansen.ai"
	defaultTimeout = 30 * time.Second
	userAgent      = "nansen-go/1.0 (+https://github.com/tigusigalpa/nansen-go)"
)

type retryConfig struct {
	maxAttempts    int
	initialDelay   time.Duration
	maxDelay       time.Duration
	retryRateLimit bool
}

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	retry      retryConfig

	SmartMoney   *SmartMoneyService
	TokenGodMode *TokenGodModeService
	Profiler     *ProfilerService
	Portfolio    *PortfolioService
	Historical   *HistoricalService
}

func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("nansen: API key is required")
	}

	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		timeout:    defaultTimeout,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	c.SmartMoney = &SmartMoneyService{client: c}
	c.TokenGodMode = &TokenGodModeService{client: c}
	c.Profiler = &ProfilerService{client: c}
	c.Portfolio = &PortfolioService{client: c}
	c.Historical = &HistoricalService{client: c}

	return c, nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, body, out interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}

	var bodyBytes []byte
	var marshalErr error
	if body != nil {
		bodyBytes, marshalErr = json.Marshal(body)
		if marshalErr != nil {
			return fmt.Errorf("nansen: failed to marshal request body: %w", marshalErr)
		}
	}

	url := c.baseURL + path

	baseCtx := ctx
	if _, ok := ctx.Deadline(); !ok && c.timeout > 0 {
		var cancel context.CancelFunc
		baseCtx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	for attempt := 0; attempt <= c.retry.maxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(baseCtx, method, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return fmt.Errorf("nansen: failed to create request: %w", err)
		}
		req.Header.Set("apiKey", c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", userAgent)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if baseCtx.Err() != nil {
				return baseCtx.Err()
			}
			if attempt == c.retry.maxAttempts {
				return fmt.Errorf("nansen: request failed: %w", err)
			}
			if err := c.wait(baseCtx, attempt); err != nil {
				return err
			}
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if out == nil {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				return nil
			}
			if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
				resp.Body.Close()
				return fmt.Errorf("nansen: failed to decode response: %w", err)
			}
			resp.Body.Close()
			return nil
		}

		apiErr := newAPIError(resp)
		resp.Body.Close()

		if c.retry.retryRateLimit && resp.StatusCode == http.StatusTooManyRequests && attempt < c.retry.maxAttempts {
			wait := parseRetryAfter(resp, c.backoff(attempt))
			if err := c.sleep(baseCtx, wait); err != nil {
				return err
			}
			continue
		}

		if resp.StatusCode >= 500 && attempt < c.retry.maxAttempts {
			if err := c.wait(baseCtx, attempt); err != nil {
				return err
			}
			continue
		}

		return apiErr
	}

	return fmt.Errorf("nansen: request exceeded maximum retry attempts")
}

func (c *Client) backoff(attempt int) time.Duration {
	delay := c.retry.initialDelay * (1 << attempt)
	if delay > c.retry.maxDelay || delay <= 0 {
		return c.retry.maxDelay
	}
	return delay
}

func (c *Client) wait(ctx context.Context, attempt int) error {
	return c.sleep(ctx, c.backoff(attempt))
}

func (c *Client) sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
