package mockcheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/mockcheck"
)

func TestMockCheck_New(t *testing.T) {
	t.Parallel()
	t.Run("With default options", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.New()
		assert.Equal(t, "mock", check.GetName())
	})

	t.Run("With custom options", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.New(
			mockcheck.WithName("custom"),
			mockcheck.WithStatus(checks.StatusFail),
		)
		assert.Equal(t, "custom", check.GetName())
	})
}

func TestMockCheck_Execute(t *testing.T) {
	t.Parallel()
	t.Run("pass", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.New(
			mockcheck.WithName("test"),
			mockcheck.WithStatus(checks.StatusPass),
		)
		results := check.Run(context.Background())
		assert.Len(t, results, 1)
		assert.Equal(t, checks.StatusPass, results[0].Status)
		assert.Equal(t, "test", check.GetName())
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.New(
			mockcheck.WithName("fail"),
			mockcheck.WithStatus(checks.StatusFail),
		)
		results := check.Run(context.Background())
		assert.Len(t, results, 1)
		assert.Equal(t, checks.StatusFail, results[0].Status)
		assert.Equal(t, "mock check failed", results[0].Output)
		assert.Equal(t, "fail", check.GetName())
	})
}
