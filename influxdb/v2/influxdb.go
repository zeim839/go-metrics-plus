package v2

import (
	"context"
	influx "github.com/influxdata/influxdb-client-go"
	"github.com/zeim839/go-metrics-plus"
	"time"
)

// Config provides a container with configuration parameters for the InfluxDB V2
// exporter. Closing the client is the caller's responsibility.
type Config struct {
	Client        influx.Client    // InfluxDB V2 Client.
	Org           string           // InfluxDB Org to Write to.
	Bucket        string           // InfluxDB Bucket to Write to.
	Registry      metrics.Registry // Registry to be exported.
	FlushInterval time.Duration    // Flush interval.
	DurationUnit  time.Duration    // Time conversion unit for durations.
	Prefix        string           // Prefix to be prepended to metric names.
}

// InfluxDBV2 is a blocking exporter function which reports metrics in r to an
// influxdb V2 client c, flushing them every d duration and prepending metric
// names with prefix. Closing the client is the caller's responsibility.
func InfluxDBV2(r metrics.Registry, d time.Duration, prefix, bucket,
	org string, c influx.Client) {
	WithConfig(Config{
		Client:        c,
		Org:           org,
		Bucket:        bucket,
		Registry:      r,
		FlushInterval: d,
		DurationUnit:  time.Nanosecond,
		Prefix:        prefix,
	})
}

// WithConfig is a blocking exporter function just like InfluxDBV2, but it takes
// a Config instead. Closing the client is the caller's responsibility.
func WithConfig(c Config) {
	for range time.Tick(c.FlushInterval) {
		influxdb(&c)
	}
}

// Once performs a single submission to InfluxDBV2. Closing the client is the
// caller's responsibility.
func Once(c Config) {
	influxdb(&c)
}

func influxdb(c *Config) {
	prefix := ""
	if c.Prefix != "" {
		prefix = c.Prefix + "."
	}

	api := c.Client.WriteAPIBlocking(c.Org, c.Bucket)
	now := time.Now().UTC()
	c.Registry.Each(func(name string, i interface{}) {
		name = prefix + name
		switch metric := i.(type) {
		case metrics.Counter:
			m := metric.Snapshot()
			p := influx.NewPoint(name, m.Labels(),
				map[string]interface{}{"count": m.Count()}, now)
			api.WritePoint(context.Background(), p)
		case metrics.Gauge:
			m := metric.Snapshot()
			p := influx.NewPoint(name, m.Labels(),
				map[string]interface{}{"gauge": m.Value()}, now)
			api.WritePoint(context.Background(), p)
		case metrics.GaugeFloat64:
			m := metric.Snapshot()
			p := influx.NewPoint(name, m.Labels(),
				map[string]interface{}{"gauge": m.Value()}, now)
			api.WritePoint(context.Background(), p)
		case metrics.Meter:
			m := metric.Snapshot()
			p := influx.NewPoint(name, m.Labels(), map[string]interface{}{
				"count":      m.Count(),
				"rate.1min":  m.Rate1(),
				"rate.5min":  m.Rate5(),
				"rate.15min": m.Rate15(),
				"rate.mean":  m.RateMean(),
			}, now)
			api.WritePoint(context.Background(), p)
		case metrics.Timer:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			p := influx.NewPoint(name, m.Labels(), map[string]interface{}{
				"count":           m.Count(),
				"min":             m.Min(),
				"max":             m.Max(),
				"mean":            m.Mean(),
				"sum":             m.Sum(),
				"variance":        m.Variance(),
				"stddev":          m.StdDev(),
				"median":          ps[0],
				"percentile.75":   ps[1],
				"percentile.95":   ps[2],
				"percentile.99.0": ps[3],
				"percentile.99.9": ps[4],
				"rate.1min":       m.Rate1(),
				"rate.5min":       m.Rate5(),
				"rate.15min":      m.Rate15(),
				"rate.mean":       m.RateMean(),
			}, now)
			api.WritePoint(context.Background(), p)
		case metrics.Histogram:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			p := influx.NewPoint(name, m.Labels(), map[string]interface{}{
				"count":           m.Count(),
				"min":             m.Min(),
				"max":             m.Max(),
				"mean":            m.Mean(),
				"sum":             m.Sum(),
				"variance":        m.Variance(),
				"stddev":          m.StdDev(),
				"median":          ps[0],
				"percentile.75":   ps[1],
				"percentile.95":   ps[2],
				"percentile.99.0": ps[3],
				"percentile.99.9": ps[4],
			}, now)
			api.WritePoint(context.Background(), p)
		}
	})
}
