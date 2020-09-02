package checks

import "github.com/brpaz/go-healthcheck"

// DummyCheck Check that always return the status that was passed as an argument. Useful for testing.
type DummyCheck struct {
	status healthcheck.Status
}

// NewDummyCheck Creates an instance of the DummyCheck
func NewDummyCheck(status healthcheck.Status) DummyCheck {
	return DummyCheck{
		status: status,
	}
}

// Execute Executes the Checks. The dummy check, returns the status based on the argument passed during the initialization.
func (c DummyCheck) Execute() map[string][]healthcheck.Check {
	return map[string][]healthcheck.Check{
		"dummy": {
			{
				Status: c.status,
			},
		},
	}
}
