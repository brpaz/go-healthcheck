// Package httpcheck provides HTTP url health checks.
// It requests HTTP urls and verifies their availability based on status codes and response times.
package httpcheck

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
)

const defaultTimeout = 5 * time.Second

// Check represents an HTTP health check that verifies url availability.
type Check struct {
	name           string
	componentType  string
	componentID    string
	url            string
	timeout        time.Duration
	exceptedStatus []int
	client         *http.Client
}

// Option is a functional option for configuring Check.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithURL sets the url of the check.
func WithURL(url string) Option {
	return func(c *Check) {
		c.url = url
	}
}

// WithTimeout sets the timeout of the check.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Check) {
		c.timeout = timeout
	}
}

// WithHTTPClient specifies a custom HTTP client to use for the health check.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Check) {
		c.client = client
	}
}

// WithExpectedStatus sets the status codes that will be considered as successful.
// By default, any status code less than 400 is considered a success.
func WithExpectedStatus(codes ...int) Option {
	return func(c *Check) {
		c.exceptedStatus = codes
	}
}

// WithComponentType sets the component type for the check.
func WithComponentType(componentType string) Option {
	return func(c *Check) {
		c.componentType = componentType
	}
}

func WithComponentID(componentID string) Option {
	return func(c *Check) {
		c.componentID = componentID
	}
}

// New creates a new HTTP Check instance with optional configuration.
func New(opts ...Option) *Check {
	check := &Check{
		name:           "http-check",
		componentType:  "http",
		componentID:    "",
		url:            "",
		timeout:        defaultTimeout,
		exceptedStatus: nil, // Use default behavior (< 400)
		client:         http.DefaultClient,
	}

	for _, opt := range opts {
		opt(check)
	}

	return check
}

// GetName returns the name of the check.
func (c *Check) GetName() string {
	return c.name
}

// Run executes the HTTP health check and returns the result.
func (c *Check) Run(ctx context.Context) checks.Result {
	// Validate configuration
	if c.url == "" {
		return checks.Result{
			Status:        checks.StatusFail,
			Output:        "URL is required for HTTP health check",
			Time:          time.Now(),
			ComponentType: c.componentType,
			ComponentID:   c.componentID,
		}
	}

	result := checks.Result{
		Status:        checks.StatusPass,
		Time:          time.Now(),
		ComponentType: c.componentType,
		ComponentID:   c.componentID,
	}

	// Create timeout context for the HTTP request
	requestCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(requestCtx, "GET", c.url, nil)
	if err != nil {
		result.Output = "failed to create request: " + err.Error()
		result.Status = checks.StatusFail
		return result
	}

	startTime := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		result.Output = "failed to execute request: " + err.Error()
		result.Status = checks.StatusFail
		return result
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	duration := time.Since(startTime)
	result.ObservedUnit = "ms"
	result.ObservedValue = duration.Milliseconds()

	// Evaluate response status
	if c.isExpectedStatusCode(resp.StatusCode) {
		result.Status = checks.StatusPass
	} else {
		result.Status = checks.StatusFail
		result.Output = "unexpected status code: " + resp.Status
	}

	return result
}

// isExpectedStatusCode determines if the given status code indicates success.
func (c *Check) isExpectedStatusCode(statusCode int) bool {
	// If specific success codes are configured, use them
	if len(c.exceptedStatus) > 0 {
		return slices.Contains(c.exceptedStatus, statusCode)
	}

	return statusCode >= 200 && statusCode < 400
}
