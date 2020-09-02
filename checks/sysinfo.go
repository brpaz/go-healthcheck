package checks

import (
	"os"
	"syscall"
	"time"

	"github.com/brpaz/go-healthcheck"
)

const cpuUsageWarnThereshold = 85
const freeRamWarnThereshold = 100

type SysInfoChecker struct {
}

// NewSysInfoChecker Creates a new instance of the SysInfoChecker
func NewSysInfoChecker() *SysInfoChecker {
	return &SysInfoChecker{}
}

func (c *SysInfoChecker) uptime(now string, si *syscall.Sysinfo_t) []healthcheck.Check {
	return []healthcheck.Check{
		{
			ComponentID:   "uptime",
			ComponentType: componentTypeSystem,
			ObservedValue: si.Uptime,
			ObservedUnit:  "s",
			Status:        healthcheck.Pass,
			Time:          now,
		},
	}
}

func (c *SysInfoChecker) hostname(now string, si *syscall.Sysinfo_t) []healthcheck.Check {
	hostname, _ := os.Hostname()
	return []healthcheck.Check{
		{
			ComponentID:   "hostname",
			ComponentType: componentTypeSystem,
			ObservedValue: hostname,
			Status:        healthcheck.Pass,
			Time:          now,
		},
	}
}

func (c *SysInfoChecker) cpuStatus(load uint64) healthcheck.Status {
	if load > cpuUsageWarnThereshold {
		return healthcheck.Warn
	}

	return healthcheck.Pass
}

func (c *SysInfoChecker) cpu(now string, si *syscall.Sysinfo_t) []healthcheck.Check {

	load1 := si.Loads[0] / 65536.0
	load5 := si.Loads[1] / 65536.0
	load15 := si.Loads[2] / 65536.0

	return []healthcheck.Check{
		{
			ComponentID:   "1 minute",
			ComponentType: componentTypeSystem,
			ObservedValue: load1,
			ObservedUnit:  "%",
			Status:        c.cpuStatus(load1),
			Time:          now,
		},
		{
			ComponentID:   "5 minutes",
			ComponentType: componentTypeSystem,
			ObservedValue: load5,
			ObservedUnit:  "%",
			Status:        c.cpuStatus(load5),
			Time:          now,
		},
		{
			ComponentID:   "15 minutes",
			ComponentType: componentTypeSystem,
			ObservedValue: load15,
			ObservedUnit:  "%",
			Status:        c.cpuStatus(load5),
			Time:          now,
		},
	}
}

func (c *SysInfoChecker) memory(now string, si *syscall.Sysinfo_t) []healthcheck.Check {
	return []healthcheck.Check{
		{
			ComponentID:   "Total Memory",
			ComponentType: componentTypeSystem,
			ObservedValue: si.Totalram / 1024 / 1024,
			ObservedUnit:  "MB",
			Status:        healthcheck.Pass,
			Time:          now,
		},
		{
			ComponentID:   "Free Memory",
			ComponentType: componentTypeSystem,
			ObservedValue: si.Freeram / 1024 / 1024,
			ObservedUnit:  "MB",
			Status:        healthcheck.Pass,
			Time:          now,
		},
	}
}

func (c *SysInfoChecker) Execute() map[string][]healthcheck.Check {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	si := &syscall.Sysinfo_t{}

	_ = syscall.Sysinfo(si)

	return map[string][]healthcheck.Check{
		"uptime":             c.uptime(now, si),
		"hostname":           c.hostname(now, si),
		"cpu:utilization":    c.cpu(now, si),
		"memory:utilization": c.memory(now, si),
	}
}
