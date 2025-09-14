package healthcheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	healthcheck "github.com/brpaz/go-healthcheck/v2"
	"github.com/brpaz/go-healthcheck/v2/checks"
	"github.com/brpaz/go-healthcheck/v2/checks/mockcheck"
)

var defaultOpts = []healthcheck.Option{
	healthcheck.WithServiceID("test-service"),
	healthcheck.WithDescription("A test service"),
	healthcheck.WithVersion("v1.0.0"),
	healthcheck.WithReleaseID("release-1"),
}

func newHealthTest(opts ...healthcheck.Option) *healthcheck.HealthCheck {
	options := append(defaultOpts, opts...)
	return healthcheck.New(options...)
}

func TestHealthCheck_New(t *testing.T) {
	t.Parallel()

	t.Run("With Default Options", func(t *testing.T) {
		t.Parallel()

		h := healthcheck.New()
		assert.NotNil(t, h)
		assert.Empty(t, h.ServiceID)
		assert.Empty(t, h.Description)
		assert.Empty(t, h.Version)
		assert.Empty(t, h.ReleaseID)
		assert.Empty(t, h.Checks)
	})

	t.Run("With Release Information", func(t *testing.T) {
		t.Parallel()

		h := healthcheck.New(
			healthcheck.WithServiceID("test-service"),
			healthcheck.WithDescription("A test service"),
			healthcheck.WithVersion("v1.0.0"),
			healthcheck.WithReleaseID("release-1"),
		)
		assert.NotNil(t, h)
		assert.Equal(t, "test-service", h.ServiceID)
		assert.Equal(t, "A test service", h.Description)
		assert.Equal(t, "v1.0.0", h.Version)
		assert.Equal(t, "release-1", h.ReleaseID)
		assert.Empty(t, h.Checks)
	})

	t.Run("With Checks", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.New(
			mockcheck.WithName("mockcheck-check"),
			mockcheck.WithStatus(checks.StatusPass),
		)
		healthcheck := newHealthTest(
			healthcheck.WithCheck(check),
		)

		assert.NotNil(t, healthcheck)
		assert.Len(t, healthcheck.Checks, 1)
		assert.Equal(t, "mockcheck-check", healthcheck.Checks[0].GetName())
	})
}

func TestHealthCheck_AddCheck(t *testing.T) {
	t.Parallel()

	healthcheck := newHealthTest()
	healthcheck.AddCheck(mockcheck.New(
		mockcheck.WithName("mockcheck-check"),
		mockcheck.WithStatus(checks.StatusPass),
	))

	assert.Len(t, healthcheck.Checks, 1)
	assert.Equal(t, "mockcheck-check", healthcheck.Checks[0].GetName())
}

func TestHealthCheck_Execute(t *testing.T) {
	t.Parallel()

	t.Run("Passes When No Checks Are Configured", func(t *testing.T) {
		t.Parallel()

		healthcheck := newHealthTest()
		response := healthcheck.Execute(context.Background())

		assert.NotNil(t, response)
		assert.Equal(t, checks.StatusPass, response.Status)
		assert.Empty(t, response.Checks)
	})

	t.Run("Passes With Single Passing Check", func(t *testing.T) {
		t.Parallel()

		check := mockcheck.New(
			mockcheck.WithName("mockcheck-check"),
			mockcheck.WithStatus(checks.StatusPass),
		)
		healthcheck := newHealthTest(
			healthcheck.WithCheck(check),
		)

		response := healthcheck.Execute(context.Background())

		assert.NotNil(t, response)

		checkResult, exists := response.Checks["mockcheck-check"]
		assert.True(t, exists)
		assert.NotEmpty(t, checkResult)
		assert.Equal(t, checks.StatusPass, checkResult[0].Status)
	})

	t.Run("Fails With Single Failing Check", func(t *testing.T) {
		t.Parallel()

		failedcheck := mockcheck.New(
			mockcheck.WithName("fail-check"),
			mockcheck.WithStatus(checks.StatusFail),
		)

		successcheck := mockcheck.New(
			mockcheck.WithName("success-check"),
			mockcheck.WithStatus(checks.StatusPass),
		)

		healthcheck := newHealthTest(
			healthcheck.WithCheck(failedcheck),
			healthcheck.WithCheck(successcheck),
		)

		response := healthcheck.Execute(context.Background())

		assert.NotNil(t, response)
		assert.Equal(t, checks.StatusFail, response.Status)
		assert.Len(t, response.Checks, 2)

		failedCheckResult, exists := response.Checks["fail-check"]
		assert.True(t, exists)

		successCheckResult, exists := response.Checks["success-check"]
		assert.True(t, exists)

		assert.Equal(t, checks.StatusFail, failedCheckResult[0].Status)
		assert.Equal(t, checks.StatusPass, successCheckResult[0].Status)
		assert.Equal(t, failedCheckResult[0].Output, "mock check failed")
	})

	t.Run("Warns With One Warning Check", func(t *testing.T) {
		t.Parallel()

		successCheck := mockcheck.New(
			mockcheck.WithName("success-check"),
			mockcheck.WithStatus(checks.StatusPass),
		)
		warningcheck := mockcheck.New(
			mockcheck.WithName("warning-check"),
			mockcheck.WithStatus(checks.StatusWarn),
		)

		healthcheck := newHealthTest(
			healthcheck.WithCheck(warningcheck),
			healthcheck.WithCheck(successCheck),
		)

		response := healthcheck.Execute(context.Background())

		assert.NotNil(t, response)
		assert.Equal(t, checks.StatusWarn, response.Status)

		warningCheckResult, exists := response.Checks["warning-check"]
		assert.True(t, exists)
		assert.Equal(t, checks.StatusWarn, warningCheckResult[0].Status)
	})
}
