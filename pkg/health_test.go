package healthcheck_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	hc "github.com/brpaz/go-healthcheck/pkg"
	"github.com/brpaz/go-healthcheck/pkg/checks"
	"github.com/brpaz/go-healthcheck/pkg/checks/mockcheck"
)

var defaultOpts = []hc.Option{
	hc.WithServiceID("test-service"),
	hc.WithDescription("A test service"),
	hc.WithVersion("v1.0.0"),
	hc.WithReleaseID("release-1"),
}

func createTestInstance(opts ...hc.Option) *hc.HealthCheck {
	options := append(defaultOpts, opts...)
	return hc.New(options...)
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("With default options", func(t *testing.T) {
		t.Parallel()

		h := hc.New()
		assert.NotNil(t, h)
		assert.Empty(t, h.ServiceID)
		assert.Empty(t, h.Description)
		assert.Empty(t, h.Version)
		assert.Empty(t, h.ReleaseID)
		assert.Empty(t, h.Checks)
	})

	t.Run("With full options", func(t *testing.T) {
		t.Parallel()

		h := hc.New(
			hc.WithServiceID("test-service"),
			hc.WithDescription("A test service"),
			hc.WithVersion("v1.0.0"),
			hc.WithReleaseID("release-1"),
		)
		assert.NotNil(t, h)
		assert.Equal(t, "test-service", h.ServiceID)
		assert.Equal(t, "A test service", h.Description)
		assert.Equal(t, "v1.0.0", h.Version)
		assert.Equal(t, "release-1", h.ReleaseID)
		assert.Empty(t, h.Checks)
	})

	t.Run("With checks", func(t *testing.T) {
		t.Parallel()
		check := mockcheck.New(
			mockcheck.WithName("mockcheck-check"),
			mockcheck.WithStatus(checks.StatusPass),
		)
		hc := createTestInstance(
			hc.WithCheck(check),
		)

		assert.NotNil(t, hc)
		assert.Len(t, hc.Checks, 1)
		assert.Equal(t, "mockcheck-check", hc.Checks[0].GetName())
	})
}

func TestExecute(t *testing.T) {
	t.Parallel()

	t.Run("HealthCheck passes with no checks configured", func(t *testing.T) {
		t.Parallel()

		hc := createTestInstance()
		response := hc.Execute(context.Background())

		assert.NotNil(t, response)
		assert.Equal(t, hc.ServiceID, response.ServiceID)
		assert.Equal(t, hc.Description, response.Description)
		assert.Equal(t, hc.Version, response.Version)
		assert.Equal(t, hc.ReleaseID, response.ReleaseID)
		assert.Empty(t, response.Output)
		assert.Equal(t, checks.StatusPass, response.Status)
		assert.Empty(t, response.Checks)
	})

	t.Run("HealthCheck passes with one passing check", func(t *testing.T) {
		t.Parallel()

		check := mockcheck.New(
			mockcheck.WithName("mockcheck-check"),
			mockcheck.WithStatus(checks.StatusPass),
		)
		hc := createTestInstance(
			hc.WithCheck(check),
		)

		response := hc.Execute(context.Background())

		assert.NotNil(t, response)

		checkResult, exists := response.Checks["mockcheck-check"]
		assert.True(t, exists)
		assert.NotEmpty(t, checkResult)
		assert.Equal(t, checks.StatusPass, checkResult[0].Status)
	})

	t.Run("HealthCheck fails with one failing check", func(t *testing.T) {
		t.Parallel()

		failedcheck := mockcheck.New(
			mockcheck.WithName("fail-check"),
			mockcheck.WithStatus(checks.StatusFail),
		)

		successcheck := mockcheck.New(
			mockcheck.WithName("success-check"),
			mockcheck.WithStatus(checks.StatusPass),
		)

		hc := createTestInstance(
			hc.WithCheck(failedcheck),
			hc.WithCheck(successcheck),
		)

		response := hc.Execute(context.Background())

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

	t.Run("HealthCheck warns with one warning check", func(t *testing.T) {
		t.Parallel()

		successCheck := mockcheck.New(
			mockcheck.WithName("success-check"),
			mockcheck.WithStatus(checks.StatusPass),
		)
		warningcheck := mockcheck.New(
			mockcheck.WithName("warning-check"),
			mockcheck.WithStatus(checks.StatusWarn),
		)

		hc := createTestInstance(
			hc.WithCheck(warningcheck),
			hc.WithCheck(successCheck),
		)

		response := hc.Execute(context.Background())

		assert.NotNil(t, response)
		assert.Equal(t, checks.StatusWarn, response.Status)

		warningCheckResult, exists := response.Checks["warning-check"]
		assert.True(t, exists)
		assert.Equal(t, checks.StatusWarn, warningCheckResult[0].Status)
	})
}

func TestHandler(t *testing.T) {
	t.Parallel()

	t.Run("Handler returns 200 OK for passing checks", func(t *testing.T) {
		t.Parallel()

		svc := createTestInstance(
			hc.WithCheck(mockcheck.New(
				mockcheck.WithName("mockcheck-check"),
				mockcheck.WithStatus(checks.StatusPass),
			)),
		)

		req, _ := http.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()

		handler := hc.Handler(svc)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/health+json", rr.Header().Get("Content-Type"))

		var response hc.Response
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, checks.StatusPass, response.Status)
		assert.Equal(t, svc.ServiceID, response.ServiceID)
		assert.Equal(t, svc.Description, response.Description)
		assert.Equal(t, svc.Version, response.Version)
		assert.Equal(t, svc.ReleaseID, response.ReleaseID)
		assert.Empty(t, response.Output)
		assert.Len(t, response.Checks, 1)

		checkResult, exists := response.Checks["mockcheck-check"]
		assert.True(t, exists)
		assert.NotEmpty(t, checkResult)
		assert.Equal(t, checks.StatusPass, checkResult[0].Status)
		assert.NotEmpty(t, checkResult[0].Time)
	})

	t.Run("Handler returns 503 with failing checks", func(t *testing.T) {
		t.Parallel()

		svc := createTestInstance(
			hc.WithCheck(mockcheck.New(
				mockcheck.WithName("fail-check"),
				mockcheck.WithStatus(checks.StatusFail),
			)),
		)

		req, _ := http.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()

		handler := hc.Handler(svc)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
		assert.Equal(t, "application/health+json", rr.Header().Get("Content-Type"))

		var response hc.Response
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, checks.StatusFail, response.Status)
		assert.Equal(t, svc.ServiceID, response.ServiceID)
		assert.Equal(t, svc.Description, response.Description)
		assert.Equal(t, svc.Version, response.Version)
		assert.Equal(t, svc.ReleaseID, response.ReleaseID)
		assert.Empty(t, response.Output)
		assert.Len(t, response.Checks, 1)

		checkResult, exists := response.Checks["fail-check"]
		assert.True(t, exists)
		assert.NotEmpty(t, checkResult)
		assert.Equal(t, checks.StatusFail, checkResult[0].Status)
	})
}
