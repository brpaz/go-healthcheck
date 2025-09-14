# HTTP Check

The HTTP Check verifies that a specific HTTP(S) endpoint is reachable and returns the expected status code. This is useful for monitoring the health of web services and APIs.

## Configuration

The HTTP Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithURL(url string)`: Sets the URL of the HTTP endpoint to check.
- `WithExpectedStatus(status []int)`: Sets a list of expected status codes. If the response status code is not in this list, the check will fail. By default any status code in the range 200-399 is considered healthy.
- `WithTimeout(timeout time.Duration)`: Sets the timeout for the HTTP request (default is 5 seconds).
- `WithHTTPClient(client *http.Client)`: Sets a custom HTTP client to be used for the request.
- `WithComponentType(componentType string)`: Sets the component type of the check (default is "http").
- `WithComponentID(componentID string)`: Sets a unique identifier for the component being checked.

## Example

```go
package main

import (
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks/httpcheck"
    "net/http"
)

func main() {
    check := httpcheck.New(
        httpcheck.WithName("Google HTTP Check"),
        httpcheck.WithURL("https://www.google.com"),
        httpcheck.WithExpectedStatus([]int{200,201}),
        httpcheck.WithTimeout(5 * time.Second),
    )
}
```
