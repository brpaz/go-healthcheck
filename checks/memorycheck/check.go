// Package memorycheck provides system memory monitoring health checks for Linux systems.
package memorycheck

import (
	"context"
	"fmt"
	"time"

	"github.com/brpaz/go-healthcheck/v2/checks"
)

// MemoryStats represents memory statistics
type MemoryStats struct {
	Total     uint64
	Available uint64
	Used      uint64
	UsedPct   float64
}

// Check represents a memory health check that monitors system memory usage.
type Check struct {
	name          string
	warnThreshold float64 // Percentage of memory usage that triggers warning
	failThreshold float64 // Percentage of memory usage that triggers failure
	reader        MemoryReader
}

// Option is a functional option for configuring Check.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithWarnThreshold sets the memory usage percentage threshold to trigger a warning status.
func WithWarnThreshold(threshold float64) Option {
	return func(c *Check) {
		c.warnThreshold = threshold
	}
}

// WithFailThreshold sets the memory usage percentage threshold to trigger a failure status.
func WithFailThreshold(threshold float64) Option {
	return func(c *Check) {
		c.failThreshold = threshold
	}
}

// WithMemoryReader sets a custom memory reader (useful for testing).
func WithMemoryReader(reader MemoryReader) Option {
	return func(c *Check) {
		c.reader = reader
	}
}

// NewCheck creates a new Memory Check instance with optional configuration.
func NewCheck(opts ...Option) *Check {
	check := &Check{
		name:          "memory",
		warnThreshold: 80.0,
		failThreshold: 95.0,
		reader:        &DefaultMemoryReader{},
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

// Run executes the memory health check and returns the result.
func (c *Check) Run(ctx context.Context) checks.Result {
	result := checks.Result{
		Status: checks.StatusPass,
		Time:   time.Now(),
	}

	// Read memory statistics
	memStats, err := c.reader.ReadMemoryStats()
	if err != nil {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("failed to read memory stats: %v", err)
		return result
	}

	result.ObservedValue = memStats.UsedPct
	result.ObservedUnit = "%"

	// Check thresholds
	if memStats.UsedPct >= c.failThreshold {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("memory usage critical: %.1f%% used (threshold: %.1f%%)",
			memStats.UsedPct, c.failThreshold)
	} else if memStats.UsedPct >= c.warnThreshold {
		result.Status = checks.StatusWarn
		result.Output = fmt.Sprintf("memory usage high: %.1f%% used (threshold: %.1f%%)",
			memStats.UsedPct, c.warnThreshold)
	} else {
		result.Status = checks.StatusPass
	}

	return result
}

// GetMemoryInfo returns current memory statistics
func (c *Check) GetMemoryInfo() (*MemoryStats, error) {
	return c.reader.ReadMemoryStats()
}
