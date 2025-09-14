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
		assert.Equal(t, "system", result.ComponentType)
		assert.Equal(t, "memory check placeholder", result.Output)
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
