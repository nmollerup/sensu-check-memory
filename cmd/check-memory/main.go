package main

import (
	"fmt"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/shirou/gopsutil/mem"
)

// Config represents the check plugin config
type Config struct {
	sensu.PluginConfig
	Warning  float64
	Critical float64
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-memory",
			Short:    "Memory check",
			Keyspace: "",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[float64]{
			Path:      "Warning",
			Argument:  "warning",
			Shorthand: "w",
			Usage:     "Warning ",
			Value:     &plugin.Warning,
		},
		&sensu.PluginConfigOption[float64]{
			Path:      "Critical",
			Argument:  "critical",
			Shorthand: "c",
			Usage:     "Critical ",
			Value:     &plugin.Critical,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	if plugin.Critical <= 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--critical is required")
	}
	if plugin.Warning <= 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--warning is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	memVirtual, _ := mem.VirtualMemory()
	if plugin.Critical <= memVirtual.UsedPercent {
		fmt.Printf("Virtual memory usage over critical treshhold! %.2f%%\n", memVirtual.UsedPercent)
		return sensu.CheckStateCritical, nil
	}
	if plugin.Warning <= memVirtual.UsedPercent {
		fmt.Printf("Virtual memory usage over warning treshhold! %.2f%%\n", memVirtual.UsedPercent)
		return sensu.CheckStateWarning, nil
	}

	fmt.Print("Virtual Memory usage OK")
	return sensu.CheckStateOK, nil
}
