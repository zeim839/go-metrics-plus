# StatsD

StatsD is the StatsD driver for [go-metrics-plus](https://github.com/zeim839/go-metrics-plus). It collects metrics from a registry and periodically posts them to a Statsd instance. Unlike other drivers, it does not support labels.

## Usage

```go
import "github.com/zeim839/go-metrics-plus/statsd"

go Statsd(metrics.DefaultRegistry, time.Second, "prefix", ":8125", "tcp")
```
