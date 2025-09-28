# Overview

The library provides several built-in checks that you can use to monitor the health of your application. Each check implements the `Check` interface and can be added to the healthcheck service.

## Built-in Checks

This library provides the following built-in checks:

- [HTTP Check](./http-check.md) - Checks that a specific http endpoint is reachable and returns the expected status code.
- [TCP Check](./tcp-check.md) - Checks that a TCP connection can be established to a specific host and port.
- [Disk Check](./disk-check.md) - Checks that a disk has enough free space.
- [Memory Check](./memory-check.md) - Checks that the system has enough free memory.
- [Database Check](./database-check.md) - Checks that a database is reachable.
- [Redis Check](./redis-check.md) - Checks that a Redis instance is reachable.
- [Mock Check](mock-check.md) - A mock check that returns the status passed to it. Useful for testing.

More checks may be added in the future. Pull requests are welcome!

## Adding your own check

To add your own check, the only requirement is to implement the `Check` interface, which requires a single method `Check(ctx context.Context) CheckResult`.

Here is an example of a custom check that verifies if a specific file exists on the filesystem:

```go
package main

import (
    "context"
    "os"
    "time"

    "github.com/brpaz/go-healthcheck/v2"
)

type FileExistsCheck struct {
    Name     string
    FilePath string
}

func (c *FileExistsCheck) Check(ctx context.Context) healthcheck.CheckResult {
    start := time.Now()
    _, err := os.Stat(c.FilePath)
    duration := time.Since(start).Milliseconds()

    if os.IsNotExist(err) {
        return healthcheck.CheckResult{
            Name:         c.Name,
            Status:       "fail",
            ComponentType: "file",
            Time:          time.Now().UTC(),
        }
    } else if err != nil {
        return healthcheck.CheckResult{
            Name:         c.Name,
            Status:       "fail",
            ComponentType: "file",
            Time:          time.Now().UTC(),
        }
    }

    return healthcheck.CheckResult{
        Name:         c.Name,
        Status:       "pass",
        ComponentType: "file",
        Time:          time.Now().UTC(),
    }
}
```

Note that each check can have multiple sub checks. This is useful when you want to group related checks together. For example, a database check can have sub checks for connection, and specific queries.


