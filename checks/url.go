package checks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brpaz/go-healthcheck"
)

// URLChecker checks if the specified URL is reachable and returns 200 OK. A timeout can be defined.
type URLChecker struct {
	URL       string
	Timeout   time.Duration
	Component string
}

// NewURLChecker Returns a new instance of the URLChecker
func NewURLChecker(component string, url string, timeout time.Duration) *URLChecker {
	return &URLChecker{Component: component, URL: url, Timeout: timeout}
}

// GetCheckName Returns the check name based on the specified component
func (c *URLChecker) GetCheckName() string {
	return c.Component + ":http"
}

// Execute Executes a check by making a HEAD request to the given URL.
// If the call returns 200 OK, the check will pass, otherwise will be marked as "Fail"
func (c *URLChecker) Execute() map[string][]healthcheck.Check {

	checkName := c.GetCheckName()

	status := healthcheck.Pass
	output := ""

	client := http.Client{
		Timeout: c.Timeout,
	}

	resp, err := client.Head(c.URL)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		status = healthcheck.Fail
		output = "HTTP Request error" + err.Error()
	}

	if resp.StatusCode != http.StatusOK {
		status = healthcheck.Fail
		output = fmt.Sprintf("Unexcepted Status Code. Endpoint returned with status %d", resp.StatusCode)
	}

	return map[string][]healthcheck.Check{
		checkName: {
			{
				ComponentID:   c.Component,
				ComponentType: componentTypeExternalSystem,
				Status:        status,
				Output:        output,
			},
		},
	}
}
