# Prometheusmetrics

Prometheusmetrics is the Prometheus driver for [go-metrics-plus](https://github.com/zeim839/go-metrics-plus). It collects metrics from a go-metrics-plus registry and periodically pushes them to a Prometheus registry (a type of the [Go Prometheus library](https://github.com/prometheus/client_golang)). The registry can then be exposed to a prometheus instance via an HTTP route (using the prometheus scraper) or via a push-gateway - the decision is implementation-specific.

Please note that if your sole intent is to integrate your application's metrics with Prometheus (and you're not planning on integrating with other programs), then the [Go Prometheus Library](https://pkg.go.dev/github.com/prometheus/client_golang) may be better suited to your needs. It offers the same metric types as Go-metrics-plus, but with better performance and its own native statistical calculations - whereas go-metrics-plus is forced to expose everything as a gauge.

## Usage

```go
import (
	pmetrics "github.com/zeim839/go-metrics-plus/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeim839/go-metrics-plus"
	"net/http"
	"time"
)

// Create Prometheus registry.
r := prometheus.NewRegistry()

// Create go-metrics-plus driver.
driver, err := pmetrics.New(metrics.DefaultRegistry, 1*time.Second,
	"namespace", "subsystem", r)

// Or... create driver with custom config.
driver, err := pmetrics.NewWithConfig(prmetrics.Config{
	// ...
}, r)

if err != nil {
	// ...
}

// Flush every 1 second.
go pr.Run()

// Expose metrics to prometheus scraper on /metrics route.
http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{Registry: r}))
log.Fatal(http.ListenAndServe(":8080", nil))
```

## Example

```go

import (
	"github.com/zeim839/go-metrics-plus"
	prmetrics "github.com/zeim839/go-metrics-plus/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func ExamplePrometheus() {
	// Register some metrics.
	metrics.GetOrRegisterTimer("myTimer", nil, nil).Update(time.Second)
	metrics.GetOrRegisterCounter("myCounter", nil, nil).Inc(50)
	metrics.GetOrRegisterMeter("myMeter", nil, nil).Mark(10)

	// Create prometheus driver.
	r := prometheus.NewRegistry()
	pr, err := prmetrics.New(metrics.DefaultRegistry, time.Second,
		"namespace", "subsystem", r)

	if err != nil {
		panic(err)
	}

	// Flush every 1 second.pp
	go pr.Run()

	// Expose metrics to prometheus scraper on /metrics route.
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{Registry: r}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
