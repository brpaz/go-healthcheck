# Disk Check

The Disk Check verifies that a disk has enough free space. This is useful for monitoring the availability of disk space on your system.

## Configuration

The Disk Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithPath(path string)`: Sets a disk path to be checked.
- `WithFileSystemStater` : Sets a custom FileSystemStater to be used for retrieving disk usage information.
- `WithWarnThreshold(threshold float64)`: Sets the disk usage percentage threshold to trigger a warning status. Default is 80.0 (80%).
- `WithFailThreshold(threshold float64)`: Sets the disk usage percentage threshold to trigger a failure status. Default is 90.0 (90%).


## Example

```go
package main

import (
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks/diskcheck"
)

func main() {
    check := diskcheck.New(
        diskcheck.WithName("disk:root"),
        diskcheck.WithPath("/"),
        diskcheck.WithWarnThreshold(75.0),
        diskcheck.WithFailThreshold(90.0),
    )
}
```
