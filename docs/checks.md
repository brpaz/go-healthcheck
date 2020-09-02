
# Checks

This package includes some built-in checks for common use cases.

## Database Check

The database check runs a simple `SELECT 1` query and returns "Failed" status if there was an error. It´s compatible with any struct that implements the native "database/sql" interface.

```go
package main

import (
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks"
)

func main() {
    health := healthcheck.New("myservice", "Some Test service", "1.0.0", "1.0.0-SPANSHOT")
    
    dbConn, _ := sql.Open("mysql", "dsn")

    // the first argument is the component name. It will be used as identifier for this check in the response
    health.AddCheckProvider(checks.NewDBChecker("main-db", dbConn))

    result := health.Get()
} 
```

## TCP Check

This check verifies that the application can connect to an external service via TCP. It is useful to check connectivity with external services that the application depends on.


```go
package main

import (
	"github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks"
)

func main() {
    health := healthcheck.New("myservice", "Some Test service", "1.0.0", "1.0.0-SPANSHOT")
    
    // the first argument is the component name. It will be used as identifier for this check in the response
    health.AddCheckProvider(checks.NewTCPChecker("my-service", "google.com:80"))

    result := health.Get()
} 
```

## URL Check

This check verifies if the application can connect to the specified URL. It works by doing an "HEAD" request to the specified URL and checks for the response status code. The check will pass if receives 200 OK.

```go
package main

import (
	"github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks"
)

func main() {
    health := healthcheck.New("myservice", "Some Test service", "1.0.0", "1.0.0-SPANSHOT")
    
    // the first argument is the component name. It will be used as identifier for this check in the response
    health.AddCheckProvider(checks.NewURLChecker("serviceA", "http://example.com", 5*time.Second))

    result := health.Get()
} 
```

## SysInfo Check

This check returns some basic metrics for the system like Uptime, Memory usage and Load.

```go
package main

import (
	"github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks"
)

func main() {
    health := healthcheck.New("myservice", "Some Test service", "1.0.0", "1.0.0-SPANSHOT")
    
    // the first argument is the component name. It will be used as identifier for this check in the response
    health.AddCheckProvider(checks.NewSysInfoChecker()))

    result := health.Get()
} 
```
