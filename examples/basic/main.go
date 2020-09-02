package main

import (
	"encoding/json"
	"os"

	"github.com/brpaz/go-healthcheck"
	"github.com/brpaz/go-healthcheck/checks"
)

func main() {
	health := healthcheck.New("myservice", "Some Test service", "1", "1").WithNote("some note")
	health.AddCheckProvider(checks.NewDummyCheck(healthcheck.Pass))
	health.AddCheckProvider(checks.NewSysInfoChecker())

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(health.Get())
}
