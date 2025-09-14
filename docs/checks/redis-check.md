# Redis Check

The Redis Checks provides a way to monitor the health of a Redis server.

## Configuration Options

The Redis Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithClient(client *redis.Client)`: Sets the Redis client to be used for the check.
- `WithTimeout(timeout time.Duration)`: Sets the timeout for the Redis PING command (default is 5 seconds).

## Example

```go
import (
	"github.com/brpaz/go-healthcheck/v2/checks/redischeck"
)

func main() {
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	check := redischeck.NewCheck(
        redischeck.WithName("My Redis Check"),
        redischeck.WithClient(redisClient),
		redischeck.WithTimeout(5*time.Second),
	)
}
