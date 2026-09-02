package nansen

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newResponse(t *testing.T, status int, body string, headers map[string]string) *http.Response {
	t.Helper()
	rec := httptest.NewRecorder()
	for k, v := range headers {
		rec.Header().Set(k, v)
	}
	rec.WriteHeader(status)
	if body != "" {
		_, _ = rec.WriteString(body)
	}
	resp := rec.Result()
	resp.Body = io.NopCloser(strings.NewReader(body))
	return resp
}

func TestAPIError_IsSentinelMapping(t *testing.T) {
	cases := []struct {
		status int
		target error
	}{
		{http.StatusNotFound, ErrNotFound},
		{http.StatusUnauthorized, ErrUnauthorized},
		{http.StatusTooManyRequests, ErrRateLimited},
		{http.StatusBadRequest, ErrBadRequest},
		{http.StatusPaymentRequired, ErrPaymentRequired},
		{http.StatusForbidden, ErrForbidden},
		{http.StatusUnprocessableEntity, ErrValidation},
		{http.StatusInternalServerError, ErrInternal},
		{http.StatusBadGateway, ErrInternal},
	}
	for _, tc := range cases {
		apiErr := &APIError{StatusCode: tc.status}
		if !errors.Is(apiErr, tc.target) {
			t.Errorf("status %d: errors.Is(apiErr, %v) = false, want true", tc.status, tc.target)
		}
	}

	apiErr := &APIError{StatusCode: http.StatusOK}
	if errors.Is(apiErr, ErrNotFound) {
		t.Errorf("status 200 should not match ErrNotFound")
	}
}

func TestNewAPIError_ParsesDetailMessage(t *testing.T) {
	resp := newResponse(t, http.StatusNotFound, `{"detail":"wallet not found"}`, nil)
	apiErr := newAPIError(resp)

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusNotFound)
	}
	if apiErr.Message != "wallet not found" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "wallet not found")
	}
}

func TestAPIError_ErrorUsesStatusText(t *testing.T) {
	err := &APIError{StatusCode: http.StatusBadGateway}
	if got, want := err.Error(), "nansen API error 502"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestNewAPIError_ParsesHeaders(t *testing.T) {
	resp := newResponse(t, http.StatusTooManyRequests, `{"detail":"slow down"}`, map[string]string{
		"Retry-After":           "30",
		"RateLimit-Remaining":   "0",
		"X-RateLimit-Remaining": "0",
	})
	apiErr := newAPIError(resp)

	if apiErr.RetryAfter == nil || *apiErr.RetryAfter != 30*time.Second {
		t.Errorf("RetryAfter = %v, want 30s", apiErr.RetryAfter)
	}
	if apiErr.RateLimitRemaining == nil || *apiErr.RateLimitRemaining != 0 {
		t.Errorf("RateLimitRemaining = %v, want pointer to 0", apiErr.RateLimitRemaining)
	}
	if apiErr.Headers.Get("Retry-After") != "30" {
		t.Errorf("Headers[Retry-After] = %q, want %q", apiErr.Headers.Get("Retry-After"), "30")
	}
}

func TestParseRetryAfter_FallsBackToExponentialBackoff(t *testing.T) {
	resp := newResponse(t, http.StatusTooManyRequests, "", nil)
	got := parseRetryAfter(resp, 2*time.Second)
	if got != 2*time.Second {
		t.Errorf("parseRetryAfter() = %v, want fallback 2s", got)
	}
}

func TestParseRetryAfter_PrefersRetryAfterHeader(t *testing.T) {
	resp := newResponse(t, http.StatusTooManyRequests, "", map[string]string{"Retry-After": "7"})
	got := parseRetryAfter(resp, 2*time.Second)
	if got != 7*time.Second {
		t.Errorf("parseRetryAfter() = %v, want 7s", got)
	}
}

func TestParseRetryAfterHeaderHTTPDateAndInvalidValue(t *testing.T) {
	future := time.Now().Add(time.Hour).UTC().Format(http.TimeFormat)
	if d, ok := parseRetryAfterHeader(newResponse(t, http.StatusTooManyRequests, "", map[string]string{"Retry-After": future})); !ok || d <= 0 {
		t.Errorf("parseRetryAfterHeader() = (%v, %v), want positive duration", d, ok)
	}
	if _, ok := parseRetryAfterHeader(newResponse(t, http.StatusTooManyRequests, "", map[string]string{"Retry-After": "invalid"})); ok {
		t.Error("parseRetryAfterHeader() accepted invalid value")
	}
}

func TestParseRetryAfter_UsesRateLimitResetWhenNoRetryAfter(t *testing.T) {
	resp := newResponse(t, http.StatusTooManyRequests, "", map[string]string{"RateLimit-Reset": "15"})
	got := parseRetryAfter(resp, 2*time.Second)
	if got != 15*time.Second {
		t.Errorf("parseRetryAfter() = %v, want 15s", got)
	}
}
