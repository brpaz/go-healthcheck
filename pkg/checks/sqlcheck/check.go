// Package sqlcheck provides SQL database health checks.
// It verifies database connectivity and optionally executes queries to validate database health.
package sqlcheck

import (
	"context"
	"database/sql"
	"time"

	"github.com/brpaz/go-healthcheck/pkg/checks"
)

const (
	defaultTimeout = 5 * time.Second
	defaultQuery   = "SELECT 1"
)

// Check represents a SQL database health check that verifies connectivity and optionally executes a query.
type Check struct {
	name     string
	db       *sql.DB
	query    string
	timeout  time.Duration
	expected interface{}
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
func WithDB(db *sql.DB) Option {
	return func(c *Check) {
		c.db = db
	}
}

// WithQuery sets a custom query to execute. If not set, "SELECT 1" is used.
func WithQuery(query string) Option {
	return func(c *Check) {
		c.query = query
	}
}

// WithTimeout sets the timeout for the database query.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Check) {
		c.timeout = timeout
	}
}

// WithExpectedResult sets the expected result for the query.
// If set, the query result will be compared against this value.
func WithExpectedResult(expected interface{}) Option {
	return func(c *Check) {
		c.expected = expected
	}
}

// New creates a new SQL Check instance with optional configuration.
func New(opts ...Option) *Check {
	check := &Check{
		name:     "sql-check",
		db:       nil,
		query:    defaultQuery,
		timeout:  defaultTimeout,
		expected: nil,
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

	result := checks.Result{
		Status:        checks.StatusPass,
		Time:          time.Now(),
		ComponentType: "database",
		ComponentID:   c.name,
	}

	// Create timeout context for the database query
	queryCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	startTime := time.Now()

	// First, check if the database is reachable with Ping
	if err := c.db.PingContext(queryCtx); err != nil {
		result.Status = checks.StatusFail
		result.Output = "database ping failed: " + err.Error()
		return []checks.Result{result}
	}

	// Execute the query if specified
	if c.query != "" {
		if err := c.executeQuery(queryCtx, &result); err != nil {
			result.Status = checks.StatusFail
			result.Output = "query execution failed: " + err.Error()
			return []checks.Result{result}
		}
	}

	duration := time.Since(startTime)
	result.ObservedUnit = "ms"
	result.ObservedValue = duration.Milliseconds()

	if result.Status == checks.StatusPass {
		result.Output = "database is healthy"
	}

	return []checks.Result{result}
}

// executeQuery executes the configured query and optionally validates the result.
func (c *Check) executeQuery(ctx context.Context, result *checks.Result) error {
	row := c.db.QueryRowContext(ctx, c.query)

	// If we have an expected result, validate it
	if c.expected != nil {
		var actual interface{}
		if err := row.Scan(&actual); err != nil {
			return err
		}

		if actual != c.expected {
			result.Status = checks.StatusWarn
			result.Output = "query result does not match expected value"
			return nil
		}
	} else {
		// Just check if the query executes without error
		var dummy interface{}
		if err := row.Scan(&dummy); err != nil {
			return err
		}
	}

	return nil
}
