// Package connectionscheck provides database connections health check implementation.
package connectionscheck

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
)

const (
	defaultTimeout       = 5 * time.Second
	defaultWarnThreshold = 0.8 // 80% of max connections
	defaultFailThreshold = 1.0 // 100% of max connections
)

// DatabaseStatsProvider interface for database instances that provide connection statistics.
type DatabaseStatsProvider interface {
	Stats() sql.DBStats
}

// ConnectionsCheck represents a database connections health check that verifies
// the number of open connections against the configured maximum.
type ConnectionsCheck struct {
	name          string
	db            DatabaseStatsProvider
	timeout       time.Duration
	warnThreshold float64
	failThreshold float64
}

// Option is a functional option for configuring ConnectionsCheck.
type Option func(*ConnectionsCheck)

// WithName sets the name of the connections check.
func WithName(name string) Option {
	return func(c *ConnectionsCheck) {
		c.name = name
	}
}

// WithDB sets the database connection to use for the health check.
func WithDB(db DatabaseStatsProvider) Option {
	return func(c *ConnectionsCheck) {
		c.db = db
	}
}

// WithTimeout sets the timeout for the connections check operation.
func WithTimeout(timeout time.Duration) Option {
	return func(c *ConnectionsCheck) {
		c.timeout = timeout
	}
}

// WithWarnThreshold sets the warning threshold as a percentage (0.0-1.0) of max connections.
func WithWarnThreshold(threshold float64) Option {
	return func(c *ConnectionsCheck) {
		c.warnThreshold = threshold
	}
}

// WithFailThreshold sets the failure threshold as a percentage (0.0-1.0) of max connections.
func WithFailThreshold(threshold float64) Option {
	return func(c *ConnectionsCheck) {
		c.failThreshold = threshold
	}
}

// New creates a new Database Connections Check instance with optional configuration.
func New(opts ...Option) *ConnectionsCheck {
	check := &ConnectionsCheck{
		name:          "db-connections-check",
		db:            nil,
		timeout:       defaultTimeout,
		warnThreshold: defaultWarnThreshold,
		failThreshold: defaultFailThreshold,
	}

	for _, opt := range opts {
		opt(check)
	}

	return check
}

// GetName returns the name of the check.
func (c *ConnectionsCheck) GetName() string {
	return c.name
}

// Run executes the database connections health check and returns the result.
func (c *ConnectionsCheck) Run(ctx context.Context) checks.Result {
	now := time.Now()
	if c.db == nil {
		return checks.Result{
			Status: checks.StatusFail,
			Output: "database connection is required",
			Time:   now,
		}
	}

	// Create timeout context for the check
	checkCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Check if context is cancelled
	select {
	case <-checkCtx.Done():
		return checks.Result{
			Status: checks.StatusFail,
			Output: "connections check timeout",
			Time:   now,
		}
	default:
	}

	// Get database statistics
	stats := c.db.Stats()
	openConnections := stats.OpenConnections
	maxConnections := stats.MaxOpenConnections

	if maxConnections <= 0 {
		return checks.Result{
			Status: checks.StatusPass,
			Time:   now,
		}
	}

	// Calculate actual threshold values
	failThresholdConnections := int(float64(maxConnections) * c.failThreshold)
	warnThresholdConnections := int(float64(maxConnections) * c.warnThreshold)

	// Check if open connections exceed the failure threshold
	if openConnections >= failThresholdConnections {
		return checks.Result{
			Status:        checks.StatusFail,
			Output:        fmt.Sprintf("open connections (%d) exceed failure threshold (%d)", openConnections, failThresholdConnections),
			Time:          now,
			ObservedValue: openConnections,
		}
	}

	// Check if we're approaching the limit (warn threshold)
	if openConnections >= warnThresholdConnections {
		return checks.Result{
			Status:        checks.StatusWarn,
			Output:        fmt.Sprintf("open connections (%d) approaching maximum (%d)", openConnections, maxConnections),
			Time:          now,
			ObservedValue: openConnections,
		}
	}

	return checks.Result{
		Status:        checks.StatusPass,
		Output:        fmt.Sprintf("open connections: %d/%d", openConnections, maxConnections),
		Time:          now,
		ObservedValue: openConnections,
	}
}
