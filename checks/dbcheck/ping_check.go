package dbcheck

import (
	"context"
	"time"

	"github.com/brpaz/go-healthcheck/v2/checks"
)

type DatabasePinger interface {
	PingContext(ctx context.Context) error
}

const (
	defaultTimeout = 5 * time.Second
)

// PingCheck represents a SQL database ping health check that verifies connectivity through ping operations.
type PingCheck struct {
	name    string
	db      DatabasePinger
	timeout time.Duration
}

// PingOption is a functional option for configuring PingCheck.
type PingOption func(*PingCheck)

// WithPingName sets the name of the ping check.
func WithPingName(name string) PingOption {
	return func(c *PingCheck) {
		c.name = name
	}
}

// WithPingDB sets the database connection to use for the ping health check.
func WithPingDB(db DatabasePinger) PingOption {
	return func(c *PingCheck) {
		c.db = db
	}
}

// WithPingTimeout sets the timeout for the database ping operation.
func WithPingTimeout(timeout time.Duration) PingOption {
	return func(c *PingCheck) {
		c.timeout = timeout
	}
}

// NewPing creates a new SQL Ping Check instance with optional configuration.
func NewPing(opts ...PingOption) *PingCheck {
	check := &PingCheck{
		name:    "database:ping",
		db:      nil,
		timeout: defaultTimeout,
	}

	for _, opt := range opts {
		opt(check)
	}

	return check
}

// GetName returns the name of the ping check.
func (c *PingCheck) GetName() string {
	return c.name
}

// Run executes the SQL Ping health check and returns the result.
func (c *PingCheck) Run(ctx context.Context) checks.Result {
	if c.db == nil {
		return checks.Result{
			Status: checks.StatusFail,
			Output: "database connection is required",
			Time:   time.Now(),
		}
	}

	now := time.Now()

	// Create timeout context for the database query
	queryCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	startTime := time.Now()

	// Check if the database is reachable with Ping
	if err := c.db.PingContext(queryCtx); err != nil {
		return checks.Result{
			Status: checks.StatusFail,
			Output: "database ping failed: " + err.Error(),
			Time:   now,
		}
	}

	duration := time.Since(startTime)

	return checks.Result{
		Status:        checks.StatusPass,
		Time:          now,
		ObservedUnit:  "ms",
		ObservedValue: duration.Milliseconds(),
	}
}
