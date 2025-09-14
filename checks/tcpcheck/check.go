// Package tcpcheck provides TCP/UDP port connectivity health checks.
// It verifies that specific network ports are open and accepting connections.
package tcpcheck

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
)

const (
	Name           = "tcp-check"
	defaultTimeout = 5 * time.Second
)

// NetworkType represents the type of network connection (tcp, udp, etc.)
type NetworkType string

const (
	TCP NetworkType = "tcp"
	UDP NetworkType = "udp"
)

// Check represents a TCP/UDP port health check that verifies connectivity.
type Check struct {
	name    string
	host    string
	port    int
	network NetworkType
	timeout time.Duration
	dialer  Dialer
}

// Dialer interface allows for custom dialers (useful for testing)
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// DefaultDialer wraps the standard net.Dialer
type DefaultDialer struct {
	*net.Dialer
}

func (d *DefaultDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.Dialer.DialContext(ctx, network, address)
}

// Option is a functional option for configuring Check.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithHost sets the host to connect to.
func WithHost(host string) Option {
	return func(c *Check) {
		c.host = host
	}
}

// WithPort sets the port to connect to.
func WithPort(port int) Option {
	return func(c *Check) {
		c.port = port
	}
}

// WithNetwork sets the network type (tcp or udp).
func WithNetwork(network NetworkType) Option {
	return func(c *Check) {
		c.network = network
	}
}

// WithTimeout sets the timeout for the connection attempt.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Check) {
		c.timeout = timeout
	}
}

// WithDialer sets a custom dialer for the connection.
func WithDialer(dialer Dialer) Option {
	return func(c *Check) {
		c.dialer = dialer
	}
}

// New creates a new TCP/UDP Check instance with optional configuration.
func New(opts ...Option) *Check {
	check := &Check{
		name:    Name,
		host:    "",
		port:    0,
		network: TCP,
		timeout: defaultTimeout,
		dialer:  &DefaultDialer{&net.Dialer{}},
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

// Run executes the TCP/UDP health check and returns the result.
func (c *Check) Run(ctx context.Context) checks.Result {
	result := checks.Result{
		Status: checks.StatusPass,
		Time:   time.Now(),
	}

	// Validate configuration
	if c.host == "" {
		result.Status = checks.StatusFail
		result.Output = "host is required"
		return result
	}

	if c.port <= 0 || c.port > 65535 {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("invalid port: %d (must be 1-65535)", c.port)
		return result
	}

	// Create timeout context for the connection attempt
	connCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	address := fmt.Sprintf("%s:%d", c.host, c.port)
	startTime := time.Now()

	// Attempt to establish connection
	conn, err := c.dialer.DialContext(connCtx, string(c.network), address)
	if err != nil {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("failed to connect to %s://%s: %v", c.network, address, err)
		return result
	}

	// Close connection immediately since we only need to verify connectivity
	if closeErr := conn.Close(); closeErr != nil {
		// Log close error but don't fail the check
		result.Output = fmt.Sprintf("connection successful but failed to close: %v", closeErr)
	}

	duration := time.Since(startTime)
	result.ObservedUnit = "ms"
	result.ObservedValue = duration.Milliseconds()

	return result
}

// Address returns the full address string for this check
func (c *Check) Address() string {
	return fmt.Sprintf("%s://%s:%d", c.network, c.host, c.port)
}
