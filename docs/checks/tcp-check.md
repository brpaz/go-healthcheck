# TCP Check

The TCP Check verifies that a TCP connection can be established to a specific host and port. This is useful for monitoring the availability of services that communicate over TCP, such as databases, message brokers, and other network services.

## Configuration

The TCP Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithHost(host string)`: Sets the hostname or IP address of the TCP endpoint to check.
- `WithPort(port int)`: Sets the port number of the TCP endpoint to check.
- `WithNetwork(network string)`: Sets the network type (e.g., "tcp", "tcp4", "tcp6"). Default is "tcp".
- `WithDiale(r(dialer *net.Dialer)`: Sets a custom net.Dialer to be used for the connection.
- `WithTimeout(timeout time.Duration)`: Sets the timeout for the TCP connection (default is 2 seconds).
- `WithComponentType(componentType string)`: Sets the component type of the check (default is "tcp").
- `WithComponentID(componentID string)`: Sets a unique identifier for the component being checked.

## Example

```go
package main

import (
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks/tcpcheck"
)

func main() {
    check := tcpcheck.New(
        tcpcheck.WithName("My TCP Check"),
        tcpcheck.WithHost("localhost"),
        tcpcheck.WithPort(8080),
        tcpcheck.WithTimeout(5 * time.Second),
    )
}
```
