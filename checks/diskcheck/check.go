// Package diskcheck provides disk space monitoring health checks.
// It monitors disk usage and alerts when thresholds are exceeded.
package diskcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/brpaz/go-healthcheck/v2/checks"
)

const (
	Name = "disk-check"
)

// DiskInfo represents disk usage information
type DiskInfo struct {
	Path     string
	Total    uint64
	Free     uint64
	Used     uint64
	UsedPct  float64
	AvailPct float64
}

// Check represents a disk space health check that monitors disk usage.
type Check struct {
	name          string
	path          string
	warnThreshold float64 // Percentage of disk usage that triggers warning
	failThreshold float64 // Percentage of disk usage that triggers failure
	stater        FileSystemStater
}

// Option is a functional option for configuring Check.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithPath sets the paths to monitor, replacing any existing paths.
func WithPath(path string) Option {
	return func(c *Check) {
		c.path = path
	}
}

// WithWarnThreshold sets the disk usage percentage that triggers a warning (default: 80%).
func WithWarnThreshold(threshold float64) Option {
	return func(c *Check) {
		c.warnThreshold = threshold
	}
}

// WithFailThreshold sets the disk usage percentage that triggers a failure (default: 90%).
func WithFailThreshold(threshold float64) Option {
	return func(c *Check) {
		c.failThreshold = threshold
	}
}

// WithFileSystemStater sets a custom filesystem stater (useful for testing).
func WithFileSystemStater(stater FileSystemStater) Option {
	return func(c *Check) {
		c.stater = stater
	}
}

// New creates a new Disk Check instance with optional configuration.
func New(opts ...Option) *Check {
	check := &Check{
		name:          Name,
		path:          "/",
		warnThreshold: 80.0,
		failThreshold: 90.0,
		stater:        &DefaultFileSystemStater{},
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

// Run executes the disk space health check and returns results for each monitored path.
// Note: This implementation temporarily returns only the first path check for RFC compliance.
// TODO: Split into separate checks per path.
func (c *Check) Run(ctx context.Context) checks.Result {
	result := checks.Result{
		Status: checks.StatusPass,
		Time:   time.Now(),
	}

	diskInfo, err := c.stater.Statfs(c.path)
	if err != nil {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("failed to get disk stats for %s: %v", c.path, err)
		return result
	}

	result.Status = checks.StatusPass
	result.ObservedValue = diskInfo.UsedPct
	result.ObservedUnit = "%"

	// Check thresholds
	if diskInfo.UsedPct >= c.failThreshold {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("disk usage critical: %.1f%% used (threshold: %.1f%%)",
			diskInfo.UsedPct, c.failThreshold)
	} else if diskInfo.UsedPct >= c.warnThreshold {
		result.Status = checks.StatusWarn
		result.Output = fmt.Sprintf("disk usage high: %.1f%% used (threshold: %.1f%%)",
			diskInfo.UsedPct, c.warnThreshold)
	}
	return result
}

// GetDiskInfo returns disk information for all monitored paths
func (c *Check) GetDiskInfo() ([]*DiskInfo, error) {
	info, err := c.stater.Statfs(c.path)
	if err != nil {
		return nil, err
	}
	return []*DiskInfo{info}, nil
}
