package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/brpaz/go-healthcheck"
	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/mockcheck"
)

var (
	testServiceID   = "test-service"
	testDescription = "A test service"
	testVersion     = "v1.0.0"
	testReleaseID   = "sha-123456"
	testCheckName   = "mockcheck-check"
)

func newHealthChechkerWithResult(t *testing.T, status checks.Status) *healthcheck.HealthCheck {
	t.Helper()
	return healthcheck.NewHealthCheck(
		healthcheck.WithServiceID(testServiceID),
		healthcheck.WithReleaseID(testReleaseID),
		healthcheck.WithVersion(testVersion),
		healthcheck.WithDescription(testDescription),
		healthcheck.WithCheck(mockcheck.New(
			mockcheck.WithName(testCheckName),
			mockcheck.WithStatus(status),
		)),
	)
}

func TestHandler(t *testing.T) {
	t.Parallel()

	t.Run("Successful Health Check", func(t *testing.T) {
		t.Parallel()

		hc := newHealthChechkerWithResult(t, checks.StatusPass)

		req, _ := http.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()

		healthHandler := healthcheck.HealthHandler(hc)
		healthHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/health+json", rr.Header().Get("Content-Type"))

		var response healthcheck.HealthHttpResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, checks.StatusPass, response.Status)
		assert.Len(t, response.Checks, 1)

		checkResult, exists := response.Checks[testCheckName]
		assert.True(t, exists)
		assert.NotEmpty(t, checkResult)
		assert.Empty(t, response.Output)
		assert.Equal(t, checks.StatusPass, checkResult[0].Status)
		assert.Equal(t, testServiceID, response.ServiceID)
		assert.Equal(t, testDescription, response.Description)
		assert.Equal(t, testVersion, response.Version)
		assert.Equal(t, testReleaseID, response.ReleaseID)
		assert.Equal(t, response.Status, checkResult[0].Status)
	})

	t.Run("Failed Health Check", func(t *testing.T) {
		t.Parallel()

		svc := newHealthChechkerWithResult(t, checks.StatusFail)

		req, _ := http.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()

		healthHandler := healthcheck.HealthHandler(svc)
		healthHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
		assert.Equal(t, "application/health+json", rr.Header().Get("Content-Type"))

		var response healthcheck.HealthHttpResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, checks.StatusFail, response.Status)
		assert.Len(t, response.Checks, 1)

		checkResult, exists := response.Checks[testCheckName]
		assert.True(t, exists)
		assert.NotEmpty(t, response.Checks)
		assert.NotEmpty(t, checkResult)
		assert.Equal(t, checks.StatusFail, checkResult[0].Status)
	})
}
