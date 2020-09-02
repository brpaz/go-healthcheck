package checks_test

import (
	"testing"

	"github.com/brpaz/go-healthcheck"
	"github.com/brpaz/go-healthcheck/checks"
	"github.com/stretchr/testify/assert"
)

func TestDummyCheck_returnsSpecifiedStatus(t *testing.T) {

	check := checks.NewDummyCheck(healthcheck.Pass)

	result := check.Execute()
	assert.Equal(t, 1, len(result))

	assert.Equal(t, healthcheck.Pass, result["dummy"][0].Status)
}
