package memorycheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/brpaz/go-healthcheck/v2/checks"
	"github.com/brpaz/go-healthcheck/v2/checks/memorycheck"
)

// MockMemoryReader implements MemoryReader for testing
type MockMemoryReader struct {
	UsedPct float64
}

func (m *MockMemoryReader) ReadMemoryStats() (*memorycheck.MemoryStats, error) {
	return &memorycheck.MemoryStats{
		Total:     8000000000, // 8GB
		Available: 2000000000, // 2GB
		Used:      6000000000, // 6GB
		UsedPct:   m.UsedPct,
	}, nil
}

func TestMemoryCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("basic memory check succeeds", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New()
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "%", result.ObservedUnit)
		assert.GreaterOrEqual(t, result.ObservedValue, 0.0)
		assert.LessOrEqual(t, result.ObservedValue, 100.0)
	})

	t.Run("custom name check", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New(memorycheck.WithName("custom-memory-check"))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "custom-memory-check", check.GetName())
	})
}

func TestMemoryCheck_GetName(t *testing.T) {
	t.Parallel()

	t.Run("returns default name", func(t *testing.T) {
		check := memorycheck.New()
		assert.Equal(t, "memory", check.GetName())
	})

	t.Run("returns custom name when set", func(t *testing.T) {
		check := memorycheck.New(memorycheck.WithName("my-memory-check"))
		assert.Equal(t, "my-memory-check", check.GetName())
	})
}

func TestMemoryCheck_Thresholds(t *testing.T) {
	t.Parallel()

	t.Run("passes when usage is below warn threshold", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New(
			memorycheck.WithWarnThreshold(80.0),
			memorycheck.WithFailThreshold(90.0),
		)
		result := check.Run(context.Background())

		assert.Equal(t, "%", result.ObservedUnit)
		assert.GreaterOrEqual(t, result.ObservedValue, 0.0)
		assert.LessOrEqual(t, result.ObservedValue, 100.0)
	})

	t.Run("supports custom thresholds", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New(
			memorycheck.WithWarnThreshold(70.0),
			memorycheck.WithFailThreshold(85.0),
		)
		result := check.Run(context.Background())
		assert.Equal(t, checks.StatusPass, result.Status)
	})
}

func TestMemoryCheck_ThresholdLogic(t *testing.T) {
	t.Parallel()

	t.Run("fail status when usage exceeds fail threshold", func(t *testing.T) {
		t.Parallel()

		mockReader := &MockMemoryReader{UsedPct: 95.0} // 95% usage
		check := memorycheck.New(
			memorycheck.WithWarnThreshold(80.0), // 80%
			memorycheck.WithFailThreshold(90.0), // 90%
			memorycheck.WithMemoryReader(mockReader),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "memory usage critical")
		assert.Contains(t, result.Output, "95.0%")
		assert.Contains(t, result.Output, "90.0%")
	})

	t.Run("warn status when usage exceeds warn threshold but not fail", func(t *testing.T) {
		t.Parallel()

		mockReader := &MockMemoryReader{UsedPct: 85.0} // 85% usage
		check := memorycheck.New(
			memorycheck.WithWarnThreshold(80.0), // 80%
			memorycheck.WithFailThreshold(90.0), // 90%
			memorycheck.WithMemoryReader(mockReader),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusWarn, result.Status)
		assert.Contains(t, result.Output, "memory usage high")
		assert.Contains(t, result.Output, "85.0%")
		assert.Contains(t, result.Output, "80.0%")
	})

	t.Run("pass status when usage is below warn threshold", func(t *testing.T) {
		t.Parallel()

		mockReader := &MockMemoryReader{UsedPct: 70.0} // 70% usage
		check := memorycheck.New(
			memorycheck.WithWarnThreshold(80.0), // 80%
			memorycheck.WithFailThreshold(90.0), // 90%
			memorycheck.WithMemoryReader(mockReader),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
	})
}
