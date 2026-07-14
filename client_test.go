package nansen

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func newTestClient(t *testing.T, srv *httptest.Server, opts ...Option) *Client {
	t.Helper()
	allOpts := append([]Option{WithBaseURL(srv.URL)}, opts...)
	c, err := New("test-api-key", allOpts...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return c
}

func TestDoRequest_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("apiKey"); got != "test-api-key" {
			t.Errorf("apiKey header = %q, want %q", got, "test-api-key")
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)

	var out map[string]string
	if err := c.doRequest(context.Background(), http.MethodPost, "/ping", nil, &out); err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	if out["status"] != "ok" {
		t.Fatalf("out = %v, want status=ok", out)
	}
}

func TestDoRequest_RetriesOn429ThenSucceeds(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"detail":"rate limited"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv, WithRetry(5, 5*time.Millisecond, 50*time.Millisecond))

	var out map[string]string
	if err := c.doRequest(context.Background(), http.MethodPost, "/ping", nil, &out); err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("attempts = %d, want 3", got)
	}
}

func TestDoRequest_ExhaustsRetriesReturnsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"detail":"rate limited"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv, WithRetry(2, 2*time.Millisecond, 10*time.Millisecond))

	err := c.doRequest(context.Background(), http.MethodPost, "/ping", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("err type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusTooManyRequests)
	}
	if apiErr.RateLimitRemaining == nil || *apiErr.RateLimitRemaining != 0 {
		t.Errorf("RateLimitRemaining = %v, want pointer to 0", apiErr.RateLimitRemaining)
	}
}

func TestDoRequest_NoRetryByDefault(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)

	err := c.doRequest(context.Background(), http.MethodPost, "/ping", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Fatalf("attempts = %d, want 1 (no retries without WithRetry)", got)
	}
}

func TestDoRequest_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := newTestClient(t, srv, WithRetry(5, 1*time.Second, 5*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := c.doRequest(ctx, http.MethodPost, "/ping", nil, nil)
	if err == nil {
		t.Fatal("expected error due to context deadline, got nil")
	}
}
