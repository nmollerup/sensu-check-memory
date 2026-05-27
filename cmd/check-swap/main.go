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
			Name:     "check-swap",
			Short:    "Swap usage check by percent",
			Keyspace: "",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[float64]{
			Path:      "warning",
			Argument:  "warning",
			Shorthand: "w",
			Usage:     "Warning threshold in percent for swap usage (default 70)",
			Value:     &plugin.Warning,
		},
		&sensu.PluginConfigOption[float64]{
			Path:      "critical",
			Argument:  "critical",
			Shorthand: "c",
			Usage:     "Critical threshold in percent for swap usage (default 80)",
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
		plugin.Warning = 70.0
	}
	if plugin.Critical == 0 {
		plugin.Critical = 80.0
	}
	if plugin.Warning >= plugin.Critical {
		return sensu.CheckStateWarning, fmt.Errorf("--warning (%.2f) must be less than --critical (%.2f)", plugin.Warning, plugin.Critical)
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	s, err := mem.SwapMemory()
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to get swap stats: %w", err)
	}

	totalMB := float64(s.Total) / 1024 / 1024
	usedMB := float64(s.Used) / 1024 / 1024
	freeMB := float64(s.Free) / 1024 / 1024

	perfdata := fmt.Sprintf("used_percent=%.2f%%;%.2f;%.2f;0;100 total_mb=%.0f used_mb=%.0f free_mb=%.0f",
		s.UsedPercent, plugin.Warning, plugin.Critical,
		totalMB, usedMB, freeMB)

	detail := fmt.Sprintf("used %.2f%%, total %.0fMB, free %.0fMB | %s", s.UsedPercent, totalMB, freeMB, perfdata)

	switch {
	case s.UsedPercent >= plugin.Critical:
		fmt.Printf("CheckSwap CRITICAL: %s\n", detail)
		return sensu.CheckStateCritical, nil
	case s.UsedPercent >= plugin.Warning:
		fmt.Printf("CheckSwap WARNING: %s\n", detail)
		return sensu.CheckStateWarning, nil
	default:
		fmt.Printf("CheckSwap OK: %s\n", detail)
		return sensu.CheckStateOK, nil
	}
}
