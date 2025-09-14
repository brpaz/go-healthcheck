# Mock Check

The Mock Check is a simple check that returns the status passed to it. This is useful for testing and simulating different health states in your application.

## Configuration

The Mock Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithStatus(status string)`: Sets the status to be returned by the check. Valid values are "pass", "warn", and "fail". Default is "pass".


## Example

```go
package main

import (
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks/mockcheck"
)

func main() {
    check := mockcheck.New(
        mockcheck.WithName("mock:example"),
        mockcheck.WithStatus("pass"),
    )
}
```
