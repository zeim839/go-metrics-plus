# go-metrics-plus

Go Metrics Library with support for Graphite, InfluxDB, Prometheus, StatsD, and AppOptics. This is a lively fork of RCrowley's [go-metrics](https://github.com/rcrowley/go-metrics) including updated backend drivers, support for labels/tags, and various optimizations.

## Install

Download as a Go module dependency:
```bash
go get github.com/zeim839/go-metrics-plus
```

Import into your project:
```go
import "github.com/zeim839/go-metrics-plus"
```

## Usage

```go
c := metrics.NewCounter(nil)
metrics.Register("foo", c)
c.Inc(47)

g := metrics.NewGauge(metrics.Labels{"key":"value"})
metrics.Register("bar", g)
g.Update(47)

r := NewRegistry()
g := metrics.NewRegisteredFunctionalGauge("cache-evictions", r, func() int64 {
	return cache.getEvictionsCount()
}, metrics.Labels{"key":"value"})

s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
h := metrics.NewHistogram(s, nil)
metrics.Register("baz", h)
h.Update(47)

m := metrics.NewMeter(nil)
metrics.Register("quux", m)
m.Mark(47)

t := metrics.NewTimer(nil)
metrics.Register("bang", t)
t.Time(func() {})
t.Update(47)
```

Register() is not threadsafe. For threadsafe metric registration use GetOrRegister:

```go
t := metrics.GetOrRegisterTimer("account.create.latency", nil, nil)
t.Time(func() {})
t.Update(47)
```

Periodically log every metric in human-readable form to standard error:
```go
import (
	"github.com/zeim839/go-metrics-plus"
	"github.com/zeim839/go-metrics-plus/logging"
	"time"
)

// Log the DefaultRegistry to stdout every second.
go logging.Logger(logging.Encode, metrics.DefaultRegistry, time.Second, "some.prefix"
```

## Publishing Metrics

* AppOptics: [Documentation](appoptics/README.md).
* Graphite: [Documentation](graphite/README.md).
* InfluxDB: [Documentation](influxdb/README.md).
* Stdout/syslog: [Documentation](logging/README.md).
* Prometheus: [Documentation](prometheus/README.md).
* StatsD: [Documentation](statsd/README.md).

## Contributing

Thank you for considering to contribute to go-metrics-plus. We accept contributions from anyone on the internet.

If you'd like to propose a new feature or report an issue, please do so in the GitHub issues section. If you'd like to contribute documentation or code, then please fork, fix, commit, and send a pull request. A maintainer will then review and merge your changes to the codebase. Pull requests must be opened on the 'master' branch.

Please make sure your contributions adhere to our coding guidelines:
* Code must adhere to the official Go [formatting](https://go.dev/doc/effective_go#formatting) guidelines (i.e uses [gofmt](https://pkg.go.dev/cmd/gofmt)).
* Code must be documented adhering to the official Go [commentary](https://go.dev/doc/effective_go#commentary) guidelines.
* Pull requests need to be based on and opened against the `master` branch.

## License

This project is an amalgamation of various open-source projects. You may assume that each directory in the source tree represents a separate package distribution, with its corresponding license documented in a LICENSE.md file contained within the directory.

The root directory is governed by go-metrics's original project license, reproduced in [LICENSE.md](LICENSE.md).
