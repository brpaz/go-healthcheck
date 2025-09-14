package dbcheck

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/brpaz/go-healthcheck/v2/checks"
)

const (
	defaultConnectionsWarnThreshold = 80.0
	defaultConnectionsFailThreshold = 100.0
)

type DatabaseStatsProvider interface {
	Stats() sql.DBStats
}

// OpenConnectionsCheck represents a database connections health check that verifies
// the number of open connections against the configured maximum.
type OpenConnectionsCheck struct {
	name          string
	db            DatabaseStatsProvider
	timeout       time.Duration
	warnThreshold float64
	failThreshold float64
}

// OpenConnectionsOption is a functional option for configuring ConnectionsCheck.
type OpenConnectionsOption func(*OpenConnectionsCheck)

// WithOpenConnectionsName sets the name of the connections check.
func WithOpenConnectionsName(name string) OpenConnectionsOption {
	return func(c *OpenConnectionsCheck) {
		c.name = name
	}
}

// WithOpenConnectionsDB sets the database connection to use for the connections health check.
func WithOpenConnectionsDB(db DatabaseStatsProvider) OpenConnectionsOption {
	return func(c *OpenConnectionsCheck) {
		c.db = db
	}
}

// WithOpenConnectionsTimeout sets the timeout for the connections check operation.
func WithOpenConnectionsTimeout(timeout time.Duration) OpenConnectionsOption {
	return func(c *OpenConnectionsCheck) {
		c.timeout = timeout
	}
}

// WithOpenConnectionsWarnThreshold sets the warning threshold as a percentage (0-100) of max connections.
func WithOpenConnectionsWarnThreshold(threshold float64) OpenConnectionsOption {
	return func(c *OpenConnectionsCheck) {
		c.warnThreshold = threshold
	}
}

// WithOpenConnectionsFailThreshold sets the failure threshold as a percentage (0-100) of max connections.
func WithOpenConnectionsFailThreshold(threshold float64) OpenConnectionsOption {
	return func(c *OpenConnectionsCheck) {
		c.failThreshold = threshold
	}
}

// NewOpenConnections creates a new Database Connections Check instance with optional configuration.
func NewOpenConnections(opts ...OpenConnectionsOption) *OpenConnectionsCheck {
	check := &OpenConnectionsCheck{
		name:          "database:connections",
		db:            nil,
		timeout:       defaultTimeout,
		warnThreshold: defaultConnectionsWarnThreshold,
		failThreshold: defaultConnectionsFailThreshold,
	}

	for _, opt := range opts {
		opt(check)
	}

	return check
}

// GetName returns the name of the connections check.
func (c *OpenConnectionsCheck) GetName() string {
	return c.name
}

// Run executes the database connections health check and returns the result.
func (c *OpenConnectionsCheck) Run(ctx context.Context) checks.Result {
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
			Output: "operation timed out",
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

	// Calculate actual threshold values - convert percentage to ratio
	failThresholdConnections := int(float64(maxConnections) * c.failThreshold / 100.0)
	warnThresholdConnections := int(float64(maxConnections) * c.warnThreshold / 100.0)

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
		Time:          now,
		ObservedValue: openConnections,
	}
}
