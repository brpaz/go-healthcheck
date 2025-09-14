package mockcheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/brpaz/go-healthcheck/v2/checks"
	"github.com/brpaz/go-healthcheck/v2/checks/mockcheck"
)

func TestMockCheck_New(t *testing.T) {
	t.Parallel()
	t.Run("With default options", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.NewCheck()
		assert.Equal(t, "mock", check.GetName())
	})

	t.Run("With custom options", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.NewCheck(
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
		check := mockcheck.NewCheck(
			mockcheck.WithName("test"),
			mockcheck.WithStatus(checks.StatusPass),
		)
		result := check.Run(context.Background())
		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test", check.GetName())
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.NewCheck(
			mockcheck.WithName("fail"),
			mockcheck.WithStatus(checks.StatusFail),
		)
		result := check.Run(context.Background())
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Equal(t, "check failed", result.Output)
		assert.Equal(t, "fail", check.GetName())
	})
}
