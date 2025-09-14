// Package sqlcheck provides SQL database health checks.
// It verifies database connectivity by performing ping operations and optionally provides metrics.
package sqlcheck

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
)

const (
	defaultTimeout = 5 * time.Second
)

type database interface {
	PingContext(ctx context.Context) error
	Stats() sql.DBStats // Add Stats method for metrics
}

// ConnectivityCheck represents a SQL database connectivity health check that verifies connectivity through ping operations.
type ConnectivityCheck struct {
	name          string
	db            database
	timeout       time.Duration
	componentType string
	componentID   string
}

// ConnectivityOption is a functional option for configuring ConnectivityCheck.
type ConnectivityOption func(*ConnectivityCheck)

// WithConnectivityName sets the name of the connectivity check.
func WithConnectivityName(name string) ConnectivityOption {
	return func(c *ConnectivityCheck) {
		c.name = name
	}
}

// WithConnectivityDB sets the database connection to use for the health check.
func WithConnectivityDB(db database) ConnectivityOption {
	return func(c *ConnectivityCheck) {
		c.db = db
	}
}

// WithConnectivityTimeout sets the timeout for the database ping operation.
func WithConnectivityTimeout(timeout time.Duration) ConnectivityOption {
	return func(c *ConnectivityCheck) {
		c.timeout = timeout
	}
}

// WithConnectivityComponentType sets the component type for the check result.
func WithConnectivityComponentType(componentType string) ConnectivityOption {
	return func(c *ConnectivityCheck) {
		c.componentType = componentType
	}
}

// WithConnectivityComponentID sets the component ID for the check result.
func WithConnectivityComponentID(componentID string) ConnectivityOption {
	return func(c *ConnectivityCheck) {
		c.componentID = componentID
	}
}

// NewConnectivityCheck creates a new SQL connectivity Check instance with optional configuration.
func NewConnectivityCheck(opts ...ConnectivityOption) *ConnectivityCheck {
	check := &ConnectivityCheck{
		name:          "sql-check",
		db:            nil,
		timeout:       defaultTimeout,
		componentType: "database",
		componentID:   "sql-check",
	}

	for _, opt := range opts {
		opt(check)
	}

	return check
}

// GetName returns the name of the check.
func (c *ConnectivityCheck) GetName() string {
	return c.name
}

// Run executes the SQL connectivity health check and returns the result.
func (c *ConnectivityCheck) Run(ctx context.Context) checks.Result {
	if c.db == nil {
		return checks.Result{
			Status:        checks.StatusFail,
			Output:        "database connection is required",
			Time:          time.Now(),
			ComponentType: c.componentType,
			ComponentID:   c.componentID,
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
			Status:        checks.StatusFail,
			Output:        "database ping failed: " + err.Error(),
			Time:          now,
			ComponentType: c.componentType,
			ComponentID:   c.componentID,
		}
	}

	duration := time.Since(startTime)

	return checks.Result{
		Status:        checks.StatusPass,
		Output:        "",
		Time:          now,
		ComponentType: c.componentType,
		ComponentID:   c.componentID,
		ObservedUnit:  "ms",
		ObservedValue: duration.Milliseconds(),
	}
}

// MetricCheck represents a SQL database metrics health check that provides specific database metrics.
type MetricCheck struct {
	name          string
	db            database
	metricType    string // e.g., "open-connections", "in-use-connections"
	componentType string
	componentID   string
}

// MetricOption is a functional option for configuring MetricCheck.
type MetricOption func(*MetricCheck)

// WithMetricName sets the name of the metric check.
func WithMetricName(name string) MetricOption {
	return func(c *MetricCheck) {
		c.name = name
	}
}

// WithMetricDB sets the database connection to use for the metric check.
func WithMetricDB(db database) MetricOption {
	return func(c *MetricCheck) {
		c.db = db
	}
}

// WithMetricType sets the type of metric to collect.
func WithMetricType(metricType string) MetricOption {
	return func(c *MetricCheck) {
		c.metricType = metricType
	}
}

// WithMetricComponentType sets the component type for the check result.
func WithMetricComponentType(componentType string) MetricOption {
	return func(c *MetricCheck) {
		c.componentType = componentType
	}
}

// WithMetricComponentID sets the component ID for the check result.
func WithMetricComponentID(componentID string) MetricOption {
	return func(c *MetricCheck) {
		c.componentID = componentID
	}
}

// NewMetricCheck creates a new SQL metric Check instance with optional configuration.
func NewMetricCheck(opts ...MetricOption) *MetricCheck {
	check := &MetricCheck{
		name:          "sql-check:metric",
		db:            nil,
		metricType:    "open-connections",
		componentType: "database",
		componentID:   "database",
	}

	for _, opt := range opts {
		opt(check)
	}

	return check
}

// GetName returns the name of the check.
func (c *MetricCheck) GetName() string {
	return c.name
}

// Run executes the SQL metric health check and returns the result.
func (c *MetricCheck) Run(ctx context.Context) checks.Result {
	if c.db == nil {
		return checks.Result{
			Status:        checks.StatusFail,
			Output:        "database connection is required",
			Time:          time.Now(),
			ComponentType: c.componentType,
			ComponentID:   c.componentID,
		}
	}

	now := time.Now()
	stats := c.db.Stats()

	var value int64
	var unit string

	switch c.metricType {
	case "open-connections":
		value = int64(stats.OpenConnections)
		unit = ""
	case "in-use-connections":
		value = int64(stats.InUse)
		unit = ""
	case "idle-connections":
		value = int64(stats.Idle)
		unit = ""
	case "max-open-connections":
		value = int64(stats.MaxOpenConnections)
		unit = ""
	case "wait-count":
		value = stats.WaitCount
		unit = ""
	case "wait-duration":
		value = stats.WaitDuration.Milliseconds()
		unit = "ms"
	default:
		return checks.Result{
			Status:        checks.StatusFail,
			Output:        fmt.Sprintf("unknown metric type: %s", c.metricType),
			Time:          now,
			ComponentType: c.componentType,
			ComponentID:   c.componentID,
		}
	}

	return checks.Result{
		Status:        checks.StatusPass,
		Output:        "",
		Time:          now,
		ComponentType: c.componentType,
		ComponentID:   c.componentID,
		ObservedUnit:  unit,
		ObservedValue: value,
	}
}

// Helper functions to create commonly used metric checks

// NewOpenConnectionsCheck creates a check for monitoring open database connections.
func NewOpenConnectionsCheck(baseName string, db database, opts ...MetricOption) *MetricCheck {
	defaultOpts := []MetricOption{
		WithMetricName(fmt.Sprintf("%s:open-connections", baseName)),
		WithMetricDB(db),
		WithMetricType("open-connections"),
	}
	allOpts := append(defaultOpts, opts...)
	return NewMetricCheck(allOpts...)
}

// NewInUseConnectionsCheck creates a check for monitoring in-use database connections.
func NewInUseConnectionsCheck(baseName string, db database, opts ...MetricOption) *MetricCheck {
	defaultOpts := []MetricOption{
		WithMetricName(fmt.Sprintf("%s:in-use-connections", baseName)),
		WithMetricDB(db),
		WithMetricType("in-use-connections"),
	}
	allOpts := append(defaultOpts, opts...)
	return NewMetricCheck(allOpts...)
}

// NewIdleConnectionsCheck creates a check for monitoring idle database connections.
func NewIdleConnectionsCheck(baseName string, db database, opts ...MetricOption) *MetricCheck {
	defaultOpts := []MetricOption{
		WithMetricName(fmt.Sprintf("%s:idle-connections", baseName)),
		WithMetricDB(db),
		WithMetricType("idle-connections"),
	}
	allOpts := append(defaultOpts, opts...)
	return NewMetricCheck(allOpts...)
}

// NewMaxOpenConnectionsCheck creates a check for monitoring max open database connections.
func NewMaxOpenConnectionsCheck(baseName string, db database, opts ...MetricOption) *MetricCheck {
	defaultOpts := []MetricOption{
		WithMetricName(fmt.Sprintf("%s:max-open-connections", baseName)),
		WithMetricDB(db),
		WithMetricType("max-open-connections"),
	}
	allOpts := append(defaultOpts, opts...)
	return NewMetricCheck(allOpts...)
}

// NewWaitCountCheck creates a check for monitoring database connection wait count.
func NewWaitCountCheck(baseName string, db database, opts ...MetricOption) *MetricCheck {
	defaultOpts := []MetricOption{
		WithMetricName(fmt.Sprintf("%s:wait-count", baseName)),
		WithMetricDB(db),
		WithMetricType("wait-count"),
	}
	allOpts := append(defaultOpts, opts...)
	return NewMetricCheck(allOpts...)
}

// NewWaitDurationCheck creates a check for monitoring database connection wait duration.
func NewWaitDurationCheck(baseName string, db database, opts ...MetricOption) *MetricCheck {
	defaultOpts := []MetricOption{
		WithMetricName(fmt.Sprintf("%s:wait-duration", baseName)),
		WithMetricDB(db),
		WithMetricType("wait-duration"),
	}
	allOpts := append(defaultOpts, opts...)
	return NewMetricCheck(allOpts...)
}

// Legacy compatibility functions

const Name = "sql-check"

// Check represents the old interface for backward compatibility.
// Deprecated: Use ConnectivityCheck and MetricCheck directly instead.
type Check struct {
	name           string
	db             database
	timeout        time.Duration
	includeMetrics bool // Flag to include connection pool metrics
}

// Option is a functional option for configuring Check.
// Deprecated: Use ConnectivityOption and MetricOption directly instead.
type Option func(*Check)

// WithName sets the name of the check.
// Deprecated: Use WithConnectivityName or WithMetricName instead.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithDB sets the database connection to use for the health check.
// Deprecated: Use WithConnectivityDB or WithMetricDB instead.
func WithDB(db database) Option {
	return func(c *Check) {
		c.db = db
	}
}

// WithTimeout sets the timeout for the database ping operation.
// Deprecated: Use WithConnectivityTimeout instead.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Check) {
		c.timeout = timeout
	}
}

// WithMetrics enables collection of database connection pool metrics.
// Deprecated: Create metric checks directly using NewMetricCheck and helper functions.
func WithMetrics(includeMetrics bool) Option {
	return func(c *Check) {
		c.includeMetrics = includeMetrics
	}
}

// New creates a new SQL Check instance with optional configuration.
// Deprecated: Use NewConnectivityCheck and metric check constructors instead.
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
// Deprecated: Use ConnectivityCheck.GetName() instead.
func (c *Check) GetName() string {
	return Name
}

// Run executes the SQL health check and returns the result.
// Deprecated: This method violates the new RFC-compliant structure.
// Use ConnectivityCheck.Run() and individual metric checks instead.
func (c *Check) Run(ctx context.Context) checks.Result {
	// For backward compatibility, just return the connectivity check result
	connectivityCheck := NewConnectivityCheck(
		WithConnectivityName(c.name),
		WithConnectivityDB(c.db),
		WithConnectivityTimeout(c.timeout),
		WithConnectivityComponentID(c.name),
	)
	return connectivityCheck.Run(ctx)
}
