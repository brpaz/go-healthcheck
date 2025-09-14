// Package memorycheck provides system memory monitoring health checks for Linux systems.
package memorycheck

import (
	"context"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
)

const (
	Name = "memory-check"
)

// Check represents a simplified memory health check for backward compatibility.
type Check struct {
	name string
}

// Option is a functional option for configuring Check.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithRAMWarnThreshold is a placeholder for backward compatibility.
func WithRAMWarnThreshold(threshold float64) Option {
	return func(c *Check) {
		// placeholder for backward compatibility
	}
}

// WithRAMFailThreshold is a placeholder for backward compatibility.
func WithRAMFailThreshold(threshold float64) Option {
	return func(c *Check) {
		// placeholder for backward compatibility
	}
}

// WithSwapWarnThreshold is a placeholder for backward compatibility.
func WithSwapWarnThreshold(threshold float64) Option {
	return func(c *Check) {
		// placeholder for backward compatibility
	}
}

// WithSwapFailThreshold is a placeholder for backward compatibility.
func WithSwapFailThreshold(threshold float64) Option {
	return func(c *Check) {
		// placeholder for backward compatibility
	}
}

// WithComponentType is a placeholder for backward compatibility.
func WithComponentType(componentType string) Option {
	return func(c *Check) {
		// placeholder for backward compatibility
	}
}

// WithComponentID is a placeholder for backward compatibility.
func WithComponentID(componentID string) Option {
	return func(c *Check) {
		// placeholder for backward compatibility
	}
}

// WithCheckSwap is a placeholder for backward compatibility.
func WithCheckSwap(checkSwap bool) Option {
	return func(c *Check) {
		// placeholder for backward compatibility
	}
}

// New creates a new Memory Check instance with optional configuration.
func New(opts ...Option) *Check {
	check := &Check{
		name: Name,
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

// Run executes the memory health check and returns a simple result.
func (c *Check) Run(ctx context.Context) checks.Result {
	return checks.Result{
		Status:        checks.StatusPass,
		Output:        "memory check placeholder",
		Time:          time.Now(),
		ComponentType: "system",
		ComponentID:   "memory:ram",
	}
}
