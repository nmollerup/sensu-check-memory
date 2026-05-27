[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nmollerup/sensu-check-memory)
[![goreleaser](https://github.com/nmollerup/sensu-check-memory/actions/workflows/release.yml/badge.svg)](https://github.com/nmollerup/sensu-check-memory/actions/workflows/release.yml) [![Go Test](https://github.com/nmollerup/sensu-check-memory/actions/workflows/test.yml/badge.svg)](https://github.com/nmollerup/sensu-check-memory/actions/workflows/test.yml)

# sensu-check-memory

Cross-platform memory and swap checks and metrics for [Sensu Go](https://sensu.io).

## Plugins

### Checks

- **check-memory** — Alert on memory usage percentage with perfdata output
- **check-swap** — Alert on swap usage percentage with perfdata output

### Metrics

- **metrics-memory** — Memory and swap statistics in Graphite plaintext format
- **metrics-memory-vmstat** — Kernel virtual memory statistics from `/proc/vmstat` in Graphite plaintext format (Linux only)

## Usage

All plugins support `--help` for full option details.

### check-memory

Checks virtual memory usage and exits with WARNING or CRITICAL when the used percentage exceeds the configured threshold. Includes perfdata for graphing.

```
-w, --warning   Warning threshold in percent (default: 80)
-c, --critical  Critical threshold in percent (default: 90)
```

**Example output:**

```
CheckMemory OK: used 42.31%, total 15987MB, available 9214MB | used_percent=42.31%;80.00;90.00;0;100 total_mb=15987 used_mb=6773 available_mb=9214 buffers_mb=512 cached_mb=2048
```

**Example:**

```bash
check-memory --warning 80 --critical 90
```

### check-swap

Checks swap usage and exits with WARNING or CRITICAL when the used percentage exceeds the configured threshold. Includes perfdata for graphing.

```
-w, --warning   Warning threshold in percent (default: 70)
-c, --critical  Critical threshold in percent (default: 80)
```

**Example output:**

```
CheckSwap OK: used 12.50%, total 2048MB, free 1792MB | used_percent=12.50%;70.00;80.00;0;100 total_mb=2048 used_mb=256 free_mb=1792
```

**Example:**

```bash
check-swap --warning 70 --critical 80
```

### metrics-memory

Outputs virtual memory and swap statistics in Graphite plaintext format. Metrics are prefixed with `<hostname>.memory` by default.

```
-s, --scheme  Metric naming prefix (default: <hostname>.memory)
```

**Metrics emitted:**

| Metric | Description |
|---|---|
| `total_bytes` | Total physical memory |
| `used_bytes` | Used memory |
| `free_bytes` | Free memory |
| `available_bytes` | Available memory (including reclaimable cache) |
| `buffers_bytes` | Kernel buffers |
| `cached_bytes` | Page cache |
| `used_percent` | Percentage of memory in use |
| `available_percent` | Percentage of memory available |
| `swap_total_bytes` | Total swap space |
| `swap_used_bytes` | Used swap space |
| `swap_free_bytes` | Free swap space |
| `swap_used_percent` | Percentage of swap in use |

**Example output:**

```
myhost.memory.total_bytes 16776138752 1716800000
myhost.memory.used_bytes 7123456789 1716800000
myhost.memory.used_percent 42.4800 1716800000
...
```

**Example:**

```bash
metrics-memory --scheme myhost.memory
```

### metrics-memory-vmstat

Reads all key/value pairs from `/proc/vmstat` and outputs them in Graphite plaintext format. Metrics are prefixed with `<hostname>.vmstat` by default.

**Linux only.** The plugin exits CRITICAL on non-Linux systems.

```
-s, --scheme  Metric naming prefix (default: <hostname>.vmstat)
```

**Example output:**

```
myhost.vmstat.nr_free_pages 123456 1716800000
myhost.vmstat.pgpgin 987654 1716800000
myhost.vmstat.pgpgout 543210 1716800000
...
```

**Example:**

```bash
metrics-memory-vmstat --scheme myhost.vmstat
```

## Metric format

All metric plugins output Graphite plaintext format:

```
<scheme>.<metric_name> <value> <unix_timestamp>
```

## Installation

Download the latest release from the [releases page](https://github.com/nmollerup/sensu-check-memory/releases) or register as a Sensu Bonsai asset.

## Configuration

Checks return standard Sensu exit codes:

| Code | Status |
|---|---|
| 0 | OK |
| 1 | WARNING |
| 2 | CRITICAL |

Metric plugins always return 0 (OK) on success and output metrics to stdout for Sensu metric collection.
