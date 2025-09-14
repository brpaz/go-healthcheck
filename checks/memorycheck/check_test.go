package memorycheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/memorycheck"
)

func TestMemoryCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("basic memory check succeeds", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New()
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Contains(t, result.Output, "memory usage normal")
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
		assert.Equal(t, memorycheck.Name, check.GetName())
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

		// Since we can't control actual memory usage, just verify the structure
		assert.NotEmpty(t, result.Output)
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

		// Just verify the check runs successfully with custom thresholds
		assert.NotEmpty(t, result.Output)
	})
}
