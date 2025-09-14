# Redis Check

The Redis Check verifies that a Redis instance is reachable and can respond to a PING command. This is useful for monitoring the availability of Redis instances in your application.

## Configuration Options

The Redis Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithClient(client *redis.Client)`: Sets the Redis client to be used for the check.
- `WithTimeout(timeout time.Duration)`: Sets the timeout for the Redis PING command (default is 5 seconds).
- `WithComponentType(componentType string)`: Sets the component type of the check (default is "redis").
- `WithComponentID(componentID string)`: Sets a unique identifier for the component being checked.

## Example

```go
import (
	"github.com/brpaz/go-healthcheck/checks/redischeck"
)

func main() {
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	check := redischeck.New(
        redischeck.WithName("My Redis Check"),
        redischeck.WithClient(redisClient),
		redischeck.WithTimeout(5*time.Second),
	)
}
