package main

import (
	"fmt"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/shirou/gopsutil/mem"
)

type Config struct {
	sensu.PluginConfig
	Warning  float64
	Critical float64
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-memory",
			Short:    "Memory usage check by percent",
			Keyspace: "",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[float64]{
			Path:      "warning",
			Argument:  "warning",
			Shorthand: "w",
			Usage:     "Warning threshold in percent for memory usage (default 80)",
			Value:     &plugin.Warning,
		},
		&sensu.PluginConfigOption[float64]{
			Path:      "critical",
			Argument:  "critical",
			Shorthand: "c",
			Usage:     "Critical threshold in percent for memory usage (default 90)",
			Value:     &plugin.Critical,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	if plugin.Warning == 0 {
		plugin.Warning = 80.0
	}
	if plugin.Critical == 0 {
		plugin.Critical = 90.0
	}
	if plugin.Warning >= plugin.Critical {
		return sensu.CheckStateWarning, fmt.Errorf("--warning (%.2f) must be less than --critical (%.2f)", plugin.Warning, plugin.Critical)
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to get memory stats: %w", err)
	}

	totalMB := float64(v.Total) / 1024 / 1024
	usedMB := float64(v.Used) / 1024 / 1024
	availMB := float64(v.Available) / 1024 / 1024
	buffersMB := float64(v.Buffers) / 1024 / 1024
	cachedMB := float64(v.Cached) / 1024 / 1024

	perfdata := fmt.Sprintf("used_percent=%.2f%%;%.2f;%.2f;0;100 total_mb=%.0f used_mb=%.0f available_mb=%.0f buffers_mb=%.0f cached_mb=%.0f",
		v.UsedPercent, plugin.Warning, plugin.Critical,
		totalMB, usedMB, availMB, buffersMB, cachedMB)

	detail := fmt.Sprintf("used %.2f%%, total %.0fMB, available %.0fMB | %s", v.UsedPercent, totalMB, availMB, perfdata)

	switch {
	case v.UsedPercent >= plugin.Critical:
		fmt.Printf("CheckMemory CRITICAL: %s\n", detail)
		return sensu.CheckStateCritical, nil
	case v.UsedPercent >= plugin.Warning:
		fmt.Printf("CheckMemory WARNING: %s\n", detail)
		return sensu.CheckStateWarning, nil
	default:
		fmt.Printf("CheckMemory OK: %s\n", detail)
		return sensu.CheckStateOK, nil
	}
}
