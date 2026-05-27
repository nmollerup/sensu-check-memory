package main

import (
	"fmt"
	"os"
	"time"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/shirou/gopsutil/mem"
)

type Config struct {
	sensu.PluginConfig
	Scheme string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "metrics-memory",
			Short:    "Memory metrics in Graphite plaintext format",
			Keyspace: "",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "scheme",
			Argument:  "scheme",
			Shorthand: "s",
			Usage:     "Metric naming scheme prefix (default: <hostname>.memory)",
			Value:     &plugin.Scheme,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	if plugin.Scheme == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		plugin.Scheme = hostname + ".memory"
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to get memory stats: %w", err)
	}
	s, err := mem.SwapMemory()
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to get swap stats: %w", err)
	}

	ts := time.Now().Unix()
	sc := plugin.Scheme

	output := func(name string, value interface{}) {
		switch v := value.(type) {
		case uint64:
			fmt.Printf("%s.%s %d %d\n", sc, name, v, ts)
		case float64:
			fmt.Printf("%s.%s %.4f %d\n", sc, name, v, ts)
		}
	}

	output("total_bytes", v.Total)
	output("used_bytes", v.Used)
	output("free_bytes", v.Free)
	output("available_bytes", v.Available)
	output("buffers_bytes", v.Buffers)
	output("cached_bytes", v.Cached)
	output("used_percent", v.UsedPercent)
	output("available_percent", 100.0-v.UsedPercent)
	output("swap_total_bytes", s.Total)
	output("swap_used_bytes", s.Used)
	output("swap_free_bytes", s.Free)
	output("swap_used_percent", s.UsedPercent)

	return sensu.CheckStateOK, nil
}
