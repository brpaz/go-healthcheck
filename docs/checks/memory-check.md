# Memory Check

The Memory Check verifies that the system has enough free memory. This is useful for monitoring the memory usage of your application and ensuring that it does not run out of memory.

## Configuration

The Memory Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithWarnThreshold(threshold float64)`: Sets the RAM usage percentage threshold to trigger a warning status. Default is 80.0 (80%).
- `WithFailThreshold(threshold float64)`: Sets the RAM usage percentage threshold to trigger


## Example

```go
package main

import (
    "github.com/brpaz/go-healthcheck/v2"
    "github.com/brpaz/go-healthcheck/v2/checks/memorycheck"
)

func main() {
    check := memorycheck.New(
        memorycheck.WithName("memory:utilization"),
        memorycheck.WithWarnThreshold(70.0),
        memorycheck.WithFailThreshold(85.0),
    )
}
```
