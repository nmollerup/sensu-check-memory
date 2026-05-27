package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

type Config struct {
	sensu.PluginConfig
	Scheme string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "metrics-memory-vmstat",
			Short:    "Kernel virtual memory statistics from /proc/vmstat in Graphite plaintext format (Linux only)",
			Keyspace: "",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "scheme",
			Argument:  "scheme",
			Shorthand: "s",
			Usage:     "Metric naming scheme prefix (default: <hostname>.vmstat)",
			Value:     &plugin.Scheme,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	if runtime.GOOS != "linux" {
		return sensu.CheckStateCritical, fmt.Errorf("metrics-memory-vmstat is only supported on Linux (current OS: %s)", runtime.GOOS)
	}
	if plugin.Scheme == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		plugin.Scheme = hostname + ".vmstat"
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	f, err := os.Open("/proc/vmstat")
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to open /proc/vmstat: %w", err)
	}
	defer f.Close()

	ts := time.Now().Unix()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) != 2 {
			continue
		}
		val, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}
		fmt.Printf("%s.%s %d %d\n", plugin.Scheme, parts[0], val, ts)
	}
	if err := scanner.Err(); err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("error reading /proc/vmstat: %w", err)
	}

	return sensu.CheckStateOK, nil
}
