
<a name="0.2.1"></a>
## [0.2.1](https://github.com/nmollerup/sensu-check-memory/compare/0.2.0...0.2.1) (2026-05-27)

### Build

* Fix goreleaser archive configuration for correct asset packaging


<a name="0.2.0"></a>
## [0.2.0](https://github.com/nmollerup/sensu-check-memory/compare/0.1.1...0.2.0) (2026-05-27)

### Features

* New `metrics-memory` command that outputs memory statistics in graphite format
* New `metrics-memory-vmstat` command that outputs vmstat statistics in graphite format (Linux only)

### Enhancements

* Include perfdata in output for check-memory and check-swap
* Setup defaults for thresholds in check-memory and check-swap

### Build

* Updated goreleaser config to version 2
* Added build configs for metrics-memory and metrics-memory-vmstat binaries

### Bug Fixes

* Fix defer return value not handled (errcheck lint)

### Dependencies

* Bump google.golang.org/grpc from 1.56.3 to 1.79.3
* Bump goreleaser/goreleaser-action from 6 to 7
* Bump github.com/sirupsen/logrus from 1.9.0 to 1.9.1
* Bump actions/checkout from 5 to 6
* Bump golangci/golangci-lint-action from 8 to 9
* Bump actions/setup-go from 5 to 6
* Bump actions/checkout from 4 to 5


<a name="0.1.1"></a>
## [0.1.1](https://github.com/nmollerup/sensu-check-memory/compare/0.1.0...0.1.1) (2025-06-04)

### Build

* Fix goreleaser deprecation warning
* Update go toolchain version

### Dependencies

* Bump golang.org/x/net from 0.23.0 to 0.38.0
* Bump github.com/golang-jwt/jwt/v4 from 4.4.2 to 4.5.2
* Bump goreleaser/goreleaser-action from 5 to 6
* Bump golangci/golangci-lint-action from 4 to 8
* Bump actions/setup-go from 5 to 6
* Bump actions/checkout from 4 to 5


<a name="0.1.0"></a>
## 0.1.0 (2024-02-21)

