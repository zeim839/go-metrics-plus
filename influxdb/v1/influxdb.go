package v1

import (
	client "github.com/influxdata/influxdb1-client"
	"github.com/zeim839/go-metrics-plus"
	"log"
	"time"
)

// Config provides a container with configuration parameters for
// the InfluxDB V1 exporter.
type Config struct {
	Client        *client.Client   // InfluxDB V1 Client.
	Database      string           // The InfluxDB Database to use.
	Registry      metrics.Registry // Registry to be exported.
	FlushInterval time.Duration    // Flush interval.
	DurationUnit  time.Duration    // Time conversion unit for durations.
	Prefix        string           // Prefix to be prepended to metric names.
}

// InfluxDBV1 is a blocking exporter function which reports metrics in r to an
// influxdb v1 client c, flushing them every d duration and prepending metric
// names with prefix.
func InfluxDBV1(r metrics.Registry, d time.Duration, prefix, db string, c *client.Client) {
	WithConfig(Config{
		Client:        c,
		Database:      db,
		Registry:      r,
		FlushInterval: d,
		DurationUnit:  time.Nanosecond,
		Prefix:        prefix,
	})
}

// WithConfig is a blocking exporter function just like InfluxDBV1,
// but it takes a Config instead.
func WithConfig(c Config) {
	//lint:ignore SA1015 TODO
	for range time.Tick(c.FlushInterval) {
		if err := influxdb(&c); nil != err {
			log.Fatal(err)
		}
	}
}

// Once performs a single submission to InfluxDB, returning a
// non-nil error on failed connections. This can be used in a loop
// similar to WithConfig for custom error handling.
func Once(c Config) error {
	return influxdb(&c)
}

func influxdb(c *Config) error {
	prefix := ""
	if c.Prefix != "" {
		prefix = c.Prefix + "."
	}

	// Data points to append.
	pts := []client.Point{}

	now := time.Now().UTC()
	c.Registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			m := metric.Snapshot()
			pts = append(pts, client.Point{
				Measurement: prefix + name,
				Time:        now,
				Fields: map[string]interface{}{
					"count": m.Count(),
				},
			})
		case metrics.Gauge:
			m := metric.Snapshot()
			pts = append(pts, client.Point{
				Measurement: prefix + name,
				Time:        now,
				Fields: map[string]interface{}{
					"gauge": m.Value(),
				},
			})
		case metrics.GaugeFloat64:
			m := metric.Snapshot()
			pts = append(pts, client.Point{
				Measurement: prefix + name,
				Time:        now,
				Fields: map[string]interface{}{
					"gauge": m.Value(),
				},
			})
		case metrics.Meter:
			m := metric.Snapshot()
			pts = append(pts, client.Point{
				Measurement: prefix + name,
				Time:        now,
				Fields: map[string]interface{}{
					"count":      m.Count(),
					"rate.1min":  m.Rate1(),
					"rate.5min":  m.Rate5(),
					"rate.15min": m.Rate15(),
					"rate.mean":  m.RateMean(),
				},
			})
		case metrics.Timer:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			pts = append(pts, client.Point{
				Measurement: prefix + name,
				Time:        now,
				Fields: map[string]interface{}{
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
				},
			})
		case metrics.Histogram:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			pts = append(pts, client.Point{
				Measurement: prefix + name,
				Time:        now,
				Fields: map[string]interface{}{
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
				},
			})
		}
	})

	bps := client.BatchPoints{
		Points:   pts,
		Database: c.Database,
	}

	_, err := c.Client.Write(bps)
	return err
}
