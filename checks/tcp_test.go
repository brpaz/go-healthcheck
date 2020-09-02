package checks_test

import (
	"testing"

	"github.com/brpaz/go-healthcheck"
	"github.com/brpaz/go-healthcheck/checks"
	"github.com/stretchr/testify/assert"
)

func TestTCPChecker_Success(t *testing.T) {

	checker := checks.NewTCPChecker("my-service", "google.com:80")

	result := checker.Execute()

	checkResult := result["my-service:connection"][0]

	assert.Equal(t, healthcheck.Pass, checkResult.Status)
}

func TestTCPChecker_Fail(t *testing.T) {

	checker := checks.NewTCPChecker("my-service", "unkown:80")

	result := checker.Execute()

	checkResult := result["my-service:connection"][0]

	assert.Equal(t, healthcheck.Fail, checkResult.Status)
	assert.NotEmpty(t, checkResult.Output)
}
