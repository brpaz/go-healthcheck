// Package sqlcheck provides SQL database health checks.
// It verifies database connectivity by performing ping operations.
package sqlcheck

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
)

const (
	Name           = "sql-check"
	defaultTimeout = 5 * time.Second
)

type database interface {
	PingContext(ctx context.Context) error
	Stats() sql.DBStats // Add Stats method for metrics
}

// Check represents a SQL database health check that verifies connectivity through ping operations.
type Check struct {
	name           string
	db             database
	timeout        time.Duration
	includeMetrics bool // Flag to include connection pool metrics
}

// Option is a functional option for configuring Check.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithDB sets the database connection to use for the health check.
func WithDB(db database) Option {
	return func(c *Check) {
		c.db = db
	}
}

// WithTimeout sets the timeout for the database ping operation.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Check) {
		c.timeout = timeout
	}
}

// WithMetrics enables collection of database connection pool metrics.
func WithMetrics(includeMetrics bool) Option {
	return func(c *Check) {
		c.includeMetrics = includeMetrics
	}
}

// New creates a new SQL Check instance with optional configuration.
func New(opts ...Option) *Check {
	check := &Check{
		name:           Name,
		db:             nil,
		timeout:        defaultTimeout,
		includeMetrics: false,
	}

	for _, opt := range opts {
		opt(check)
	}

	return check
}

// GetName returns the name of the check.
func (c *Check) GetName() string {
	return Name
}

// Run executes the SQL health check and returns the result.
func (c *Check) Run(ctx context.Context) []checks.Result {
	if c.db == nil {
		return []checks.Result{{
			Status:        checks.StatusFail,
			Output:        "database connection is required",
			Time:          time.Now(),
			ComponentType: "database",
			ComponentID:   c.name,
		}}
	}

	now := time.Now()
	var results []checks.Result

	// Create timeout context for the database query
	queryCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	startTime := time.Now()

	// First, check if the database is reachable with Ping
	if err := c.db.PingContext(queryCtx); err != nil {
		result := checks.Result{
			Status:        checks.StatusFail,
			Output:        "database ping failed: " + err.Error(),
			Time:          now,
			ComponentType: "database",
			ComponentID:   c.name,
		}
		return []checks.Result{result}
	}

	duration := time.Since(startTime)

	// Main connectivity check
	connectivityResult := checks.Result{
		Status:        checks.StatusPass,
		Output:        "",
		Time:          now,
		ComponentType: "database",
		ComponentID:   c.name,
		ObservedUnit:  "ms",
		ObservedValue: duration.Milliseconds(),
	}
	results = append(results, connectivityResult)

	// Add separate metric checks if enabled
	if c.includeMetrics {
		stats := c.db.Stats()

		// Open connections check
		results = append(results, checks.Result{
			Status:        checks.StatusPass,
			Output:        "",
			Time:          now,
			ComponentType: "database",
			ComponentID:   fmt.Sprintf("%s:open-connections", c.name),
			ObservedUnit:  "",
			ObservedValue: int64(stats.OpenConnections),
		})

		// In-use connections check
		results = append(results, checks.Result{
			Status:        checks.StatusPass,
			Output:        "",
			Time:          now,
			ComponentType: "database",
			ComponentID:   fmt.Sprintf("%s:in-use-connections", c.name),
			ObservedUnit:  "",
			ObservedValue: int64(stats.InUse),
		})

		// Idle connections check
		results = append(results, checks.Result{
			Status:        checks.StatusPass,
			Output:        "",
			Time:          now,
			ComponentType: "database",
			ComponentID:   fmt.Sprintf("%s:idle-connections", c.name),
			ObservedUnit:  "",
			ObservedValue: int64(stats.Idle),
		})

		// Max open connections check
		results = append(results, checks.Result{
			Status:        checks.StatusPass,
			Output:        "",
			Time:          now,
			ComponentType: "database",
			ComponentID:   fmt.Sprintf("%s:max-open-connections", c.name),
			ObservedUnit:  "",
			ObservedValue: int64(stats.MaxOpenConnections),
		})

		// Wait count check
		results = append(results, checks.Result{
			Status:        checks.StatusPass,
			Output:        "",
			Time:          now,
			ComponentType: "database",
			ComponentID:   fmt.Sprintf("%s:wait-count", c.name),
			ObservedUnit:  "",
			ObservedValue: stats.WaitCount,
		})

		// Wait duration check
		results = append(results, checks.Result{
			Status:        checks.StatusPass,
			Output:        "",
			Time:          now,
			ComponentType: "database",
			ComponentID:   fmt.Sprintf("%s:wait-duration", c.name),
			ObservedUnit:  "ms",
			ObservedValue: stats.WaitDuration.Milliseconds(),
		})
	}

	return results
}
