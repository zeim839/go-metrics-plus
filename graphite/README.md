# Graphite

Graphite is the Graphite driver for [go-metrics-plus](https://github.com/zeim839/go-metrics-plus). It collects metrics from a registry and periodically posts them to a Graphite instance. The driver has been ported to go-metrics-plus (along with several bug fixes and optimizations) from cyberdelia's [go-metrics-graphite](https://github.com/cyberdelia/go-metrics-graphite) project.

## Usage

```go
import "github.com/zeim839/go-metrics-plus/graphite"

// Sinks metrics every 1 second.
go graphite.Graphite(metrics.DefaultRegistry, 1*time.Second, "some.prefix", addr)
```

## Example

```go
import (
	"github.com/zeim839/go-metrics-plus"
	"github.com/zeim839/go-metrics-plus/graphite"
	"net"
	"time"
)

func ExampleGraphite() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go graphite.Graphite(metrics.DefaultRegistry, 1*time.Second, "some.prefix", addr)
}

func ExampleWithConfig() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go graphite.WithConfig(Config{
		Addr:          addr,
		Registry:      metrics.DefaultRegistry,
		FlushInterval: 1 * time.Second,
		DurationUnit:  time.Millisecond,
	})
}
```
