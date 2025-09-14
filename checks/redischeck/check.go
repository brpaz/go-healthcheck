// Package redischeck provides simple Redis health checks.
// It verifies Redis connectivity and ping operations.
package redischeck

import (
	"context"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
)

const (
	Name           = "redis-check"
	defaultTimeout = 5 * time.Second
)

// RedisClient defines the interface for Redis operations needed for health checks
type RedisClient interface {
	Ping(ctx context.Context) error
	Close() error
}

// Check represents a Redis health check that verifies connectivity and basic operations.
type Check struct {
	name          string
	client        RedisClient
	timeout       time.Duration
	componentType string
	componentID   string
}

// Option is a functional option for configuring Check.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithClient sets the Redis client to use for the health check.
func WithClient(client RedisClient) Option {
	return func(c *Check) {
		c.client = client
	}
}

// WithTimeout sets the timeout for Redis operations.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Check) {
		c.timeout = timeout
	}
}

// WithComponentType sets the component type for the check.
func WithComponentType(componentType string) Option {
	return func(c *Check) {
		c.componentType = componentType
	}
}

// WithComponentID sets the component ID for the check.
func WithComponentID(componentID string) Option {
	return func(c *Check) {
		c.componentID = componentID
	}
}

// New creates a new Redis Check instance with optional configuration.
func New(opts ...Option) *Check {
	check := &Check{
		name:          Name,
		client:        nil,
		timeout:       defaultTimeout,
		componentType: "datastore",
		componentID:   "redis",
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

// Run executes the Redis health check and returns the result.
func (c *Check) Run(ctx context.Context) []checks.Result {
	if c.client == nil {
		return []checks.Result{{
			Status:        checks.StatusFail,
			Output:        "Redis client is required",
			Time:          time.Now(),
			ComponentType: c.componentType,
			ComponentID:   c.componentID,
		}}
	}

	result := checks.Result{
		Status:        checks.StatusPass,
		Time:          time.Now(),
		ComponentType: c.componentType,
		ComponentID:   c.componentID,
	}

	// Create timeout context for Redis operations
	redisCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	startTime := time.Now()

	// Perform ping operation
	if err := c.client.Ping(redisCtx); err != nil {
		result.Status = checks.StatusFail
		result.Output = "Redis ping failed: " + err.Error()
		return []checks.Result{result}
	}

	duration := time.Since(startTime)
	result.ObservedUnit = "ms"
	result.ObservedValue = duration.Milliseconds()

	return []checks.Result{result}
}
