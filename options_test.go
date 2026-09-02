package nansen

import (
	"net/http"
	"testing"
	"time"
)

func TestWithRetry_ZeroAttemptsAllowsZeroDelays(t *testing.T) {
	c, err := New("key", WithRetry(0, 0, 0))
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}
	if c.retry.maxAttempts != 0 {
		t.Errorf("maxAttempts = %d, want 0", c.retry.maxAttempts)
	}
}

func TestWithRetry_PositiveAttemptsRequirePositiveDelays(t *testing.T) {
	if _, err := New("key", WithRetry(3, 0, time.Second)); err == nil {
		t.Error("expected error for zero initialDelay with maxAttempts > 0")
	}
	if _, err := New("key", WithRetry(3, time.Second, 0)); err == nil {
		t.Error("expected error for zero maxDelay with maxAttempts > 0")
	}
	if _, err := New("key", WithRetry(-1, time.Second, time.Second)); err == nil {
		t.Error("expected error for negative maxAttempts")
	}
}

func TestWithBaseURL_RejectsEmpty(t *testing.T) {
	if _, err := New("key", WithBaseURL("")); err == nil {
		t.Error("expected error for empty base URL")
	}
}

func TestWithBaseURL_ValidatesAndNormalizes(t *testing.T) {
	c, err := New("key", WithBaseURL("https://example.com/"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if c.baseURL != "https://example.com" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "https://example.com")
	}

	for _, baseURL := range []string{"example.com", "ftp://example.com", "https://example.com?foo=bar"} {
		if _, err := New("key", WithBaseURL(baseURL)); err == nil {
			t.Errorf("New() with base URL %q returned nil error", baseURL)
		}
	}
}

func TestNew_RejectsNilOption(t *testing.T) {
	if _, err := New("key", nil); err == nil {
		t.Error("expected error for nil option")
	}
}

func TestWithTimeout_RejectsNonPositive(t *testing.T) {
	if _, err := New("key", WithTimeout(0)); err == nil {
		t.Error("expected error for zero timeout")
	}
}

func TestClientOptionsApplyValues(t *testing.T) {
	httpClient := &http.Client{}
	c, err := New("key", WithHTTPClient(httpClient), WithTimeout(time.Second))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if c.httpClient != httpClient || c.timeout != time.Second {
		t.Error("client options were not applied")
	}
	if _, err := New("key", WithHTTPClient(nil)); err == nil {
		t.Error("expected error for nil HTTP client")
	}
}

func TestNew_RequiresAPIKey(t *testing.T) {
	if _, err := New(""); err == nil {
		t.Error("expected error for empty API key")
	}
}
