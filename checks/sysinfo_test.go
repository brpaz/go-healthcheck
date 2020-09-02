package checks_test

import (
	"testing"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/stretchr/testify/assert"
)

func TestSysinfo(t *testing.T) {

	checker := checks.NewSysInfoChecker()

	result := checker.Execute()

	assert.NotNil(t, "uptime", result)
	assert.NotNil(t, "cpu:utilization", result)
	assert.NotNil(t, "memory:utilization", result)
}
