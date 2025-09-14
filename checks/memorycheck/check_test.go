package memorycheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/memorycheck"
)

func TestMemoryCheck_New(t *testing.T) {
	t.Parallel()

	t.Run("creates check with default values", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New()

		assert.NotNil(t, check)
		assert.Equal(t, "memory-check", check.GetName())
	})

	t.Run("creates check with custom options", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New(
			memorycheck.WithName("custom-memory-check"),
			memorycheck.WithRAMWarnThreshold(70.0),
			memorycheck.WithRAMFailThreshold(85.0),
			memorycheck.WithSwapWarnThreshold(40.0),
			memorycheck.WithSwapFailThreshold(70.0),
			memorycheck.WithComponentType("resource"),
			memorycheck.WithComponentID("system-memory"),
			memorycheck.WithSwapCheck(false),
		)

		assert.NotNil(t, check)
		assert.Equal(t, "custom-memory-check", check.GetName())
	})
}

func TestMemoryCheck_GetName(t *testing.T) {
	t.Parallel()

	check := memorycheck.New()
	assert.Equal(t, "memory-check", check.GetName())
}

func TestMemoryCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("runs successfully with real system memory", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New()
		results := check.Run(context.Background())

		// Should have at least RAM result
		assert.NotEmpty(t, results)

		// First result should be RAM
		ramResult := results[0]
		assert.Equal(t, "system", ramResult.ComponentType)
		assert.Equal(t, "memory:ram", ramResult.ComponentID)
		if observedValue, ok := ramResult.ObservedValue.(float64); ok {
			assert.True(t, observedValue >= 0)
		}
		assert.Equal(t, "%", ramResult.ObservedUnit)

		// Status should be valid
		assert.Contains(t, []checks.Status{checks.StatusPass, checks.StatusWarn, checks.StatusFail}, ramResult.Status)
	})

	t.Run("skips swap check when disabled", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New(memorycheck.WithSwapCheck(false))
		results := check.Run(context.Background())

		// Should only have RAM result
		assert.Len(t, results, 1)
		assert.Equal(t, "memory:ram", results[0].ComponentID)
	})

	t.Run("includes swap check when enabled and available", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New(memorycheck.WithSwapCheck(true))
		results := check.Run(context.Background())

		// Should have at least RAM result
		assert.NotEmpty(t, results)
		assert.Equal(t, "memory:ram", results[0].ComponentID)

		// May have swap result if swap is available on the system
		if len(results) > 1 {
			swapResult := results[1]
			assert.Equal(t, "memory:swap", swapResult.ComponentID)
			if observedValue, ok := swapResult.ObservedValue.(float64); ok {
				assert.True(t, observedValue >= 0)
			}
			assert.Equal(t, "%", swapResult.ObservedUnit)
		}
	})
}

func TestMemoryCheck_GetMemoryInfo(t *testing.T) {
	t.Parallel()

	check := memorycheck.New()

	info, err := check.GetMemoryInfo()

	require.NoError(t, err)
	require.NotNil(t, info)

	// Verify we get reasonable memory values
	assert.Greater(t, info.TotalRAM, uint64(0), "Total RAM should be greater than 0")
	assert.LessOrEqual(t, info.UsedRAM, info.TotalRAM, "Used RAM should not exceed total RAM")
	assert.True(t, info.UsedRAMPct >= 0 && info.UsedRAMPct <= 100, "RAM percentage should be between 0 and 100")

	// Swap may or may not be available
	if info.TotalSwap > 0 {
		assert.LessOrEqual(t, info.UsedSwap, info.TotalSwap, "Used swap should not exceed total swap")
		assert.True(t, info.UsedSwapPct >= 0 && info.UsedSwapPct <= 100, "Swap percentage should be between 0 and 100")
	}
}

func TestMemoryCheck_CustomThresholds(t *testing.T) {
	t.Parallel()

	t.Run("uses custom RAM thresholds", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New(
			memorycheck.WithRAMWarnThreshold(70.0),
			memorycheck.WithRAMFailThreshold(85.0),
		)

		results := check.Run(context.Background())

		assert.NotEmpty(t, results)

		// The actual status depends on real system memory usage,
		// but we can verify the result structure is correct
		ramResult := results[0]
		assert.Equal(t, "memory:ram", ramResult.ComponentID)
	})

	t.Run("uses custom swap thresholds", func(t *testing.T) {
		t.Parallel()

		check := memorycheck.New(
			memorycheck.WithSwapWarnThreshold(40.0),
			memorycheck.WithSwapFailThreshold(70.0),
		)

		results := check.Run(context.Background())

		assert.NotEmpty(t, results)

		// If swap is available, verify the structure
		for _, result := range results {
			if result.ComponentID == "memory:swap" {
				assert.Equal(t, "%", result.ObservedUnit)
			}
		}
	})
}
