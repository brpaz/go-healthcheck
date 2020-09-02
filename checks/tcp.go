package checks

import (
	"fmt"
	"net"

	"github.com/brpaz/go-healthcheck"
)

// TCPChecker checks if a resource is reacbable by doing a TCP request to the specified host.
// The host can be defined in the format "host:port"
type TCPChecker struct {
	Component string
	Host      string
}

// NewTCPChecker Initializes a new instance of the "TCP Checker"
func NewTCPChecker(component string, host string) TCPChecker {
	return TCPChecker{
		Component: component,
		Host:      host,
	}
}

// Execute Executes the check by doing a TCP request to the specified host.
func (c *TCPChecker) Execute() map[string][]healthcheck.Check {

	checkName := fmt.Sprintf("%s:connection", c.Component)

	status := healthcheck.Pass
	output := ""

	conn, err := net.Dial("tcp", c.Host)

	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		status = healthcheck.Fail
		output = err.Error()
	}

	return map[string][]healthcheck.Check{
		checkName: {
			{
				ComponentID:   c.Component,
				ComponentType: "network-service",
				Status:        status,
				Output:        output,
			},
		},
	}
}
