package healthcheck_test

import (
	"testing"

	"github.com/brpaz/go-healthcheck"
	"github.com/brpaz/go-healthcheck/checks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	health := healthcheck.New("my-service", "service description", "1.0.0", "1.0.0-release").
		WithLink("example", "https://example.com").
		WithOutput("test").
		WithNote("test note")

	assert.Equal(t, "my-service", health.ServiceID)
	assert.Equal(t, "1.0.0-release", health.ReleaseID)
	assert.Equal(t, "1.0.0", health.Version)
	assert.Equal(t, "test", health.Output)
	assert.Equal(t, []string{"test note"}, health.Notes)
	assert.Equal(t, "https://example.com", health.Links["example"])
}

func TestAddCheckProvider(t *testing.T) {
	health := healthcheck.New("my-service", "service description", "1.0.0", "1.0.0-release")
	health.AddCheckProvider(checks.NewDummyCheck(healthcheck.Pass))

	result := health.Get()

	assert.Equal(t, 1, len(result.Checks))
	assert.Equal(t, healthcheck.Pass, result.Checks["dummy"][0].Status)
	assert.Equal(t, healthcheck.Pass, result.Status)
}

func TestFailedCheck_ReturnsGlobalFailStatus(t *testing.T) {
	health := healthcheck.New("my-service", "service description", "1.0.0", "1.0.0-release")
	health.AddCheckProvider(checks.NewDummyCheck(healthcheck.Fail))

	result := health.Get()

	assert.Equal(t, 1, len(result.Checks))
	assert.Equal(t, healthcheck.Fail, result.Checks["dummy"][0].Status)
	assert.Equal(t, healthcheck.Fail, result.Status)
}

func TestWarnCheck_ReturnsGlobalWarnStatus(t *testing.T) {
	health := healthcheck.New("my-service", "service description", "1.0.0", "1.0.0-release")
	health.AddCheckProvider(checks.NewDummyCheck(healthcheck.Warn))
	health.AddCheckProvider(checks.NewSysInfoChecker())

	result := health.Get()

	assert.Equal(t, healthcheck.Warn, result.Status)
}
