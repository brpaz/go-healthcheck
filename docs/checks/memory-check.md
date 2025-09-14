# Memory Check

The Memory Check verifies that the system has enough free memory. This is useful for monitoring the memory usage of your application and ensuring that it does not run out of memory.

## Configuration

The Memory Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithComponentType(componentType string)`: Sets the component type of the check (default is "memory").
- `WithComponentID(componentID string)`: Sets a unique identifier for the component being checked If not set, it defaults to "memory_check".
- `WithRAMWarnThreshold(threshold float64)`: Sets the RAM usage percentage threshold to trigger a warning status. Default is 80.0 (80%).
- `WithRAMFailThreshold(threshold float64)`: Sets the RAM usage percentage threshold to trigger a failure status. Default is 90.0 (90%).
- `WithSwapWarnThreshold(threshold float64)`: Sets the Swap usage percentage threshold to trigger a warning status. Default is 50.0 (50%).
- `WithSwapFailThreshold(threshold float64)`: Sets the Swap usage percentage threshold to trigger a failure status. Default is 75.0 (75%).
- `WithCheckSwap(checkSwap bool)`: Enables or disables checking swap memory. Default is false.


## Example

```go
package main

import (
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks/memorycheck"
)

func main() {
    check := memorycheck.New(
        memorycheck.WithName("My Memory Check"),
        memorycheck.WithRAMWarnThreshold(70.0),
        memorycheck.WithRAMFailThreshold(85.0),
        memorycheck.WithCheckSwap(true),
        memorycheck.WithSwapWarnThreshold(40.0),
        memorycheck.WithSwapFailThreshold(60.0),
    )
}
```
