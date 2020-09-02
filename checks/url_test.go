package checks_test

import (
	"testing"
	"time"

	"github.com/brpaz/go-healthcheck"
	"github.com/brpaz/go-healthcheck/checks"
	"github.com/stretchr/testify/assert"
)

func TestURLCheck_Success(t *testing.T) {

	checker := checks.NewURLChecker("serviceA", "http://example.com", 5*time.Second)

	checks := checker.Execute()

	result := checks[checker.GetCheckName()][0]

	assert.Equal(t, 1, len(checks))
	assert.Equal(t, healthcheck.Pass, result.Status)
}

func TestURLheck_Fail(t *testing.T) {

	checker := checks.NewURLChecker("serviceA", "https://httpbin.org/status/500", 5*time.Second)

	checks := checker.Execute()

	result := checks[checker.GetCheckName()][0]

	assert.Equal(t, 1, len(checks))
	assert.Equal(t, healthcheck.Fail, result.Status)
	assert.Equal(t, "Unexcepted Status Code. Endpoint returned with status 500", result.Output)
}
