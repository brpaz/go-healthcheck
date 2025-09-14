package mockcheck

import (
	"context"
	"time"

	"github.com/brpaz/go-healthcheck/v2/checks"
)

// Check is a mock implementation of the Check interface for testing purposes.
// It returns a single check result with the specified result status.
type Check struct {
	name   string
	status checks.Status
}

// Option is a functional option for configuring the MockCheck.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithStatus sets the status of the check.
func WithStatus(status checks.Status) Option {
	return func(c *Check) {
		c.status = status
	}
}

// NewCheck creates a new MockCheck instance with optional configuration.
func NewCheck(opts ...Option) *Check {
	m := &Check{
		name:   "mock",
		status: checks.StatusPass,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// GetName returns the name of the check.
func (c *Check) GetName() string {
	return c.name
}

// Execute runs the mock check and returns a single check result based on the configured status.
func (c *Check) Run(ctx context.Context) checks.Result {
	var output string
	if c.status == checks.StatusFail {
		output = "check failed"
	}

	result := checks.Result{
		Status: c.status,
		Output: output,
		Time:   time.Now(),
	}
	return result
}
