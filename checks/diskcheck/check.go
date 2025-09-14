// Package diskcheck provides disk space monitoring health checks.
// It monitors disk usage and alerts when thresholds are exceeded.
package diskcheck

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
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

// FileSystemStater defines the interface for getting filesystem statistics
type FileSystemStater interface {
	Statfs(path string) (*DiskInfo, error)
}

// DefaultFileSystemStater implements FileSystemStater using syscall.Statfs
type DefaultFileSystemStater struct{}

func (d *DefaultFileSystemStater) Statfs(path string) (*DiskInfo, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, fmt.Errorf("failed to get filesystem stats for %s: %w", path, err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - (stat.Bfree * uint64(stat.Bsize))

	var usedPct, availPct float64
	if total > 0 {
		usedPct = float64(used) / float64(total) * 100
		availPct = float64(free) / float64(total) * 100
	}

	return &DiskInfo{
		Path:     path,
		Total:    total,
		Free:     free,
		Used:     used,
		UsedPct:  usedPct,
		AvailPct: availPct,
	}, nil
}

// Check represents a disk space health check that monitors disk usage.
type Check struct {
	name          string
	paths         []string
	warnThreshold float64 // Percentage of disk usage that triggers warning
	failThreshold float64 // Percentage of disk usage that triggers failure
	componentType string
	componentID   string
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

// WithPaths sets the paths to monitor, replacing any existing paths.
func WithPaths(paths ...string) Option {
	return func(c *Check) {
		c.paths = paths
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

// WithComponentType sets the component type for the check.
func WithComponentType(componentType string) Option {
	return func(c *Check) {
		c.componentType = componentType
	}
}

// WithComponentID sets the component ID for the check.
func WithComponentID(componentID string) Option {
	return func(c *Check) {
		c.componentID = componentID
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
		paths:         []string{"/"},
		warnThreshold: 80.0,
		failThreshold: 90.0,
		componentType: "system",
		componentID:   "disk",
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
func (c *Check) Run(ctx context.Context) []checks.Result {
	var results []checks.Result

	for _, path := range c.paths {
		result := c.checkPath(path)
		results = append(results, result)
	}

	return results
}

// checkPath checks disk usage for a single path
func (c *Check) checkPath(path string) checks.Result {
	result := checks.Result{
		Status:        checks.StatusPass,
		Time:          time.Now(),
		ComponentType: c.componentType,
		ComponentID:   fmt.Sprintf("%s:%s", c.componentID, path),
	}

	diskInfo, err := c.stater.Statfs(path)
	if err != nil {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("failed to get disk stats for %s: %v", path, err)
		return result
	}

	result.ObservedValue = diskInfo.UsedPct
	result.ObservedUnit = "%"

	// Check thresholds
	if diskInfo.UsedPct >= c.failThreshold {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("disk usage critical on %s: %.1f%% used (threshold: %.1f%%)",
			path, diskInfo.UsedPct, c.failThreshold)
	} else if diskInfo.UsedPct >= c.warnThreshold {
		result.Status = checks.StatusWarn
		result.Output = fmt.Sprintf("disk usage high on %s: %.1f%% used (threshold: %.1f%%)",
			path, diskInfo.UsedPct, c.warnThreshold)
	} else {
		result.Status = checks.StatusPass
		result.Output = fmt.Sprintf("disk usage normal on %s: %.1f%% used (%.1f%% available)",
			path, diskInfo.UsedPct, diskInfo.AvailPct)
	}

	return result
}

// GetDiskInfo returns disk information for all monitored paths
func (c *Check) GetDiskInfo() ([]*DiskInfo, error) {
	var diskInfos []*DiskInfo

	for _, path := range c.paths {
		info, err := c.stater.Statfs(path)
		if err != nil {
			return nil, err
		}
		diskInfos = append(diskInfos, info)
	}

	return diskInfos, nil
}
