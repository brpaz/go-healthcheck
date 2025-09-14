// Package memorycheck provides system memory monitoring health checks for Linux systems.
// It monitors RAM and swap usage and alerts when thresholds are exceeded.
package memorycheck

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brpaz/go-healthcheck/checks"
)

const (
	Name = "memory-check"
)

// MemoryInfo represents system memory usage information
type MemoryInfo struct {
	TotalRAM     uint64  // Total RAM in bytes
	AvailableRAM uint64  // Available RAM in bytes
	UsedRAM      uint64  // Used RAM in bytes
	UsedRAMPct   float64 // Used RAM percentage

	TotalSwap     uint64  // Total swap in bytes
	AvailableSwap uint64  // Available swap in bytes
	UsedSwap      uint64  // Used swap in bytes
	UsedSwapPct   float64 // Used swap percentage
}

// getMemoryInfo gets memory information from /proc/meminfo
func getMemoryInfo() (*MemoryInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/meminfo: %w", err)
	}
	defer func() { _ = file.Close() }()

	memInfo := &MemoryInfo{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		valueStr := fields[1]
		value, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			continue
		}

		// Convert from KB to bytes
		value *= 1024

		switch key {
		case "MemTotal":
			memInfo.TotalRAM = value
		case "MemAvailable":
			memInfo.AvailableRAM = value
		case "SwapTotal":
			memInfo.TotalSwap = value
		case "SwapFree":
			memInfo.AvailableSwap = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read /proc/meminfo: %w", err)
	}

	// Calculate used memory and percentages
	memInfo.UsedRAM = memInfo.TotalRAM - memInfo.AvailableRAM
	if memInfo.TotalRAM > 0 {
		memInfo.UsedRAMPct = float64(memInfo.UsedRAM) / float64(memInfo.TotalRAM) * 100
	}

	memInfo.UsedSwap = memInfo.TotalSwap - memInfo.AvailableSwap
	if memInfo.TotalSwap > 0 {
		memInfo.UsedSwapPct = float64(memInfo.UsedSwap) / float64(memInfo.TotalSwap) * 100
	}

	return memInfo, nil
}

// Check represents a memory health check that monitors system memory usage.
type Check struct {
	name              string
	ramWarnThreshold  float64 // RAM usage percentage that triggers warning
	ramFailThreshold  float64 // RAM usage percentage that triggers failure
	swapWarnThreshold float64 // Swap usage percentage that triggers warning
	swapFailThreshold float64 // Swap usage percentage that triggers failure
	componentType     string
	componentID       string
	checkSwap         bool // Whether to check swap usage
}

// Option is a functional option for configuring Check.
type Option func(*Check)

// WithName sets the name of the check.
func WithName(name string) Option {
	return func(c *Check) {
		c.name = name
	}
}

// WithRAMWarnThreshold sets the RAM usage percentage that triggers a warning (default: 80%).
func WithRAMWarnThreshold(threshold float64) Option {
	return func(c *Check) {
		c.ramWarnThreshold = threshold
	}
}

// WithRAMFailThreshold sets the RAM usage percentage that triggers a failure (default: 90%).
func WithRAMFailThreshold(threshold float64) Option {
	return func(c *Check) {
		c.ramFailThreshold = threshold
	}
}

// WithSwapWarnThreshold sets the swap usage percentage that triggers a warning (default: 50%).
func WithSwapWarnThreshold(threshold float64) Option {
	return func(c *Check) {
		c.swapWarnThreshold = threshold
	}
}

// WithSwapFailThreshold sets the swap usage percentage that triggers a failure (default: 80%).
func WithSwapFailThreshold(threshold float64) Option {
	return func(c *Check) {
		c.swapFailThreshold = threshold
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

// WithSwapCheck enables or disables swap usage checking.
func WithSwapCheck(enabled bool) Option {
	return func(c *Check) {
		c.checkSwap = enabled
	}
}

// New creates a new Memory Check instance with optional configuration.
func New(opts ...Option) *Check {
	check := &Check{
		name:              Name,
		ramWarnThreshold:  80.0,
		ramFailThreshold:  90.0,
		swapWarnThreshold: 50.0,
		swapFailThreshold: 80.0,
		componentType:     "system",
		componentID:       "memory",
		checkSwap:         false,
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

// Run executes the memory health check and returns results for RAM and optionally swap.
func (c *Check) Run(ctx context.Context) []checks.Result {
	var results []checks.Result

	memInfo, err := getMemoryInfo()
	if err != nil {
		result := checks.Result{
			Status:        checks.StatusFail,
			Output:        fmt.Sprintf("failed to get memory info: %v", err),
			Time:          time.Now(),
			ComponentType: c.componentType,
			ComponentID:   c.componentID,
		}
		return []checks.Result{result}
	}

	// Check RAM usage
	ramResult := c.checkRAM(memInfo)
	results = append(results, ramResult)

	// Check swap usage if enabled and swap is available
	if c.checkSwap && memInfo.TotalSwap > 0 {
		swapResult := c.checkSwap_usage(memInfo)
		results = append(results, swapResult)
	}

	return results
}

// checkRAM checks RAM usage against thresholds
func (c *Check) checkRAM(memInfo *MemoryInfo) checks.Result {
	result := checks.Result{
		Status:        checks.StatusPass,
		Time:          time.Now(),
		ComponentType: c.componentType,
		ComponentID:   c.componentID + ":ram",
		ObservedValue: memInfo.UsedRAMPct,
		ObservedUnit:  "%",
	}

	if memInfo.UsedRAMPct >= c.ramFailThreshold {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("RAM usage critical: %.1f%% used (threshold: %.1f%%)",
			memInfo.UsedRAMPct, c.ramFailThreshold)
	} else if memInfo.UsedRAMPct >= c.ramWarnThreshold {
		result.Status = checks.StatusWarn
		result.Output = fmt.Sprintf("RAM usage high: %.1f%% used (threshold: %.1f%%)",
			memInfo.UsedRAMPct, c.ramWarnThreshold)
	} else {
		result.Status = checks.StatusPass
	}

	return result
}

// checkSwap_usage checks swap usage against thresholds
func (c *Check) checkSwap_usage(memInfo *MemoryInfo) checks.Result {
	result := checks.Result{
		Status:        checks.StatusPass,
		Time:          time.Now(),
		ComponentType: c.componentType,
		ComponentID:   c.componentID + ":swap",
		ObservedValue: memInfo.UsedSwapPct,
		ObservedUnit:  "%",
	}

	if memInfo.UsedSwapPct >= c.swapFailThreshold {
		result.Status = checks.StatusFail
		result.Output = fmt.Sprintf("swap usage critical: %.1f%% used (threshold: %.1f%%)",
			memInfo.UsedSwapPct, c.swapFailThreshold)
	} else if memInfo.UsedSwapPct >= c.swapWarnThreshold {
		result.Status = checks.StatusWarn
		result.Output = fmt.Sprintf("swap usage high: %.1f%% used (threshold: %.1f%%)",
			memInfo.UsedSwapPct, c.swapWarnThreshold)
	} else {
		result.Status = checks.StatusPass
	}

	return result
}

// GetMemoryInfo returns current system memory information
func (c *Check) GetMemoryInfo() (*MemoryInfo, error) {
	return getMemoryInfo()
}
