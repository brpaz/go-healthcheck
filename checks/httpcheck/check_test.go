package httpcheck_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/httpcheck"
)

// customRoundTripper is a test helper that adds a custom header to requests
type customRoundTripper struct {
	next http.RoundTripper
}

func (c *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Custom-Header", "test-value")
	return c.next.RoundTrip(req)
}

func TestHTTPCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("successful check with 200 status", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}))
		defer server.Close()

		check := httpcheck.New(
			httpcheck.WithName("test-check"),
			httpcheck.WithURL(server.URL),
			httpcheck.WithTimeout(1*time.Second),
			httpcheck.WithComponentID("test-component"),
			httpcheck.WithComponentType("http"),
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test-component", result.ComponentID)
		assert.Equal(t, "http", result.ComponentType)
		assert.Equal(t, "ms", result.ObservedUnit)
		assert.GreaterOrEqual(t, result.ObservedValue, int64(0))
	})

	t.Run("failed check with 500 status", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		check := httpcheck.New(
			httpcheck.WithName("test-check"),
			httpcheck.WithURL(server.URL),
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "unexpected status code")
	})

	t.Run("failed check with empty endpoint", func(t *testing.T) {
		t.Parallel()

		check := httpcheck.New(
			httpcheck.WithName("test-check"),
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "URL is required for HTTP health check")
	})

	t.Run("successful check with custom success status codes", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated) // 201
		}))
		defer server.Close()

		check := httpcheck.New(
			httpcheck.WithName("test-check"),
			httpcheck.WithURL(server.URL),
			httpcheck.WithExpectedStatus(201, 202),
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusPass, result.Status)
	})

	t.Run("timeout check fails when server is slow", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond) // Delay response
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		check := httpcheck.New(
			httpcheck.WithName("timeout-test"),
			httpcheck.WithURL(server.URL),
			httpcheck.WithTimeout(50*time.Millisecond), // Short timeout
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "failed to execute request")
		assert.Contains(t, result.Output, "context deadline exceeded")
	})

	t.Run("custom HTTP client with custom transport", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if custom header was added by our transport
			if r.Header.Get("X-Custom-Header") == "test-value" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("Custom client worked"))
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}))
		defer server.Close()

		// Create a custom transport that adds a header
		customTransport := &customRoundTripper{
			next: http.DefaultTransport,
		}

		customClient := &http.Client{
			Transport: customTransport,
			Timeout:   5 * time.Second,
		}

		check := httpcheck.New(
			httpcheck.WithName("custom-client-test"),
			httpcheck.WithURL(server.URL),
			httpcheck.WithHTTPClient(customClient),
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Empty(t, result.Output) // Successful checks don't set output
	})

	t.Run("custom HTTP client timeout is respected", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		// Create a custom client with very short timeout
		customClient := &http.Client{
			Timeout: 50 * time.Millisecond,
		}

		check := httpcheck.New(
			httpcheck.WithName("client-timeout-test"),
			httpcheck.WithURL(server.URL),
			httpcheck.WithHTTPClient(customClient),
			httpcheck.WithTimeout(1*time.Second), // This should be overridden by client timeout
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "failed to execute request")
	})
}

func TestHTTPCheck_GetName(t *testing.T) {
	t.Parallel()

	t.Run("returns custom name when set", func(t *testing.T) {
		check := httpcheck.New(
			httpcheck.WithName("my-custom-check"),
			httpcheck.WithURL("http://example.com"),
		)

		assert.Equal(t, "my-custom-check", check.GetName())
	})

	t.Run("returns default name when not set", func(t *testing.T) {
		check := httpcheck.New(
			httpcheck.WithURL("http://example.com"),
		)

		assert.Equal(t, "http-check", check.GetName())
	})
}
