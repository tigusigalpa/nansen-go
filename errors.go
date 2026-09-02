package nansen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ErrNotFound through ErrInternal identify standard Nansen API error categories.
var (
	ErrNotFound        = errors.New("nansen: resource not found")
	ErrUnauthorized    = errors.New("nansen: unauthorized")
	ErrRateLimited     = errors.New("nansen: rate limited")
	ErrBadRequest      = errors.New("nansen: bad request")
	ErrPaymentRequired = errors.New("nansen: payment required")
	ErrForbidden       = errors.New("nansen: forbidden")
	ErrValidation      = errors.New("nansen: validation error")
	ErrInternal        = errors.New("nansen: internal server error")
)

// APIError describes an unsuccessful response returned by the Nansen API.
type APIError struct {
	StatusCode int
	Message    string
	RawBody    []byte
	Headers    http.Header

	// RetryAfter is populated when the response included a Retry-After
	// header, expressing how long the caller should wait before retrying.
	RetryAfter *time.Duration

	// RateLimitRemaining reflects the RateLimit-Remaining or
	// X-RateLimit-Remaining header, when present, indicating how many
	// requests remain in the current rate-limit window.
	RateLimitRemaining *int
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("nansen API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("nansen API error %d", e.StatusCode)
}

// Is reports whether APIError corresponds to a known Nansen error category.
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrNotFound:
		return e.StatusCode == http.StatusNotFound
	case ErrUnauthorized:
		return e.StatusCode == http.StatusUnauthorized
	case ErrRateLimited:
		return e.StatusCode == http.StatusTooManyRequests
	case ErrBadRequest:
		return e.StatusCode == http.StatusBadRequest
	case ErrPaymentRequired:
		return e.StatusCode == http.StatusPaymentRequired
	case ErrForbidden:
		return e.StatusCode == http.StatusForbidden
	case ErrValidation:
		return e.StatusCode == http.StatusUnprocessableEntity
	case ErrInternal:
		return e.StatusCode >= http.StatusInternalServerError && e.StatusCode <= 599
	}
	return false
}

func newAPIError(resp *http.Response) *APIError {
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		raw = []byte("")
	}
	raw = bytes.TrimSpace(raw)

	errBody := struct {
		Detail json.RawMessage `json:"detail"`
	}{}
	message := ""
	if len(raw) > 0 {
		if decodeErr := json.Unmarshal(raw, &errBody); decodeErr == nil && len(errBody.Detail) > 0 {
			if errBody.Detail[0] == '"' {
				message = strings.Trim(string(errBody.Detail), `"`)
			} else {
				message = string(raw)
			}
		} else {
			message = string(raw)
		}
	}
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Message:    message,
		RawBody:    raw,
		Headers:    resp.Header.Clone(),
	}

	if d, ok := parseRetryAfterHeader(resp); ok {
		apiErr.RetryAfter = &d
	}
	if n, ok := parseRateLimitRemaining(resp); ok {
		apiErr.RateLimitRemaining = &n
	}

	return apiErr
}

// parseRetryAfterHeader extracts the Retry-After header as a duration,
// supporting both delay-seconds and HTTP-date formats (RFC 7231).
func parseRetryAfterHeader(resp *http.Response) (time.Duration, bool) {
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0, false
	}
	if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second, true
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d, true
		}
		return 0, true
	}
	return 0, false
}

// parseRateLimitRemaining reads RateLimit-Remaining or its legacy
// X-RateLimit-Remaining alias.
func parseRateLimitRemaining(resp *http.Response) (int, bool) {
	for _, header := range []string{"RateLimit-Remaining", "X-RateLimit-Remaining"} {
		if v := resp.Header.Get(header); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				return n, true
			}
		}
	}
	return 0, false
}

// parseRetryAfter determines how long to wait before the next retry attempt.
// Priority order: Retry-After header (seconds or HTTP-date), then
// RateLimit-Reset / X-RateLimit-Reset, then the caller-supplied fallback
// (exponential backoff). When RateLimit-Remaining / X-RateLimit-Remaining
// indicates requests are still available, the server-reported reset time is
// still honored since a 429 was already returned for this request.
func parseRetryAfter(resp *http.Response, fallback time.Duration) time.Duration {
	if d, ok := parseRetryAfterHeader(resp); ok {
		if d > 0 {
			return d
		}
		return fallback
	}
	if v := resp.Header.Get("RateLimit-Reset"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
	}
	if v := resp.Header.Get("X-RateLimit-Reset"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
	}
	// RateLimit-Remaining / X-RateLimit-Remaining are informational only
	// (no explicit reset time provided); fall back to exponential backoff.
	return fallback
}
