package logging

import (
	"fmt"
	"github.com/zeim839/go-metrics-plus"
	"io"
)

// EncodeStatsd encodes a metric into statsd line protocol. Some interfaces
// are encoded as multi-line summaries. Healthchecks are not supported. Labels
// are not (natively) supported by Statsd. It is assumed that the sampling rate
// is the same as the flush rate configured on the Statsd server.
func EncodeStatsd(w io.Writer, name, prefix string, i interface{}) {
	if prefix != "" {
		prefix = prefix + "."
	}
	head := prefix + name
	switch metric := i.(type) {
	case metrics.Counter:
		fmt.Fprintf(w, "%s:%d|c\n", head, metric.Count())
		metric.Clear()
	case metrics.Gauge:
		fmt.Fprintf(w, "%s:%d|g\n", head, metric.Value())
	case metrics.GaugeFloat64:
		fmt.Fprintf(w, "%s:%f|g\n", head, metric.Value())
	case metrics.Meter:
		m := metric.Snapshot()
		fmt.Fprintf(w, "%s.count:%d|c\n", head, m.Count())
		fmt.Fprintf(w, "%s.rate.1min:%f|g\n", head, m.Rate1())
		fmt.Fprintf(w, "%s.rate.5min:%f|g\n", head, m.Rate5())
		fmt.Fprintf(w, "%s.rate.15min:%f|g\n", head, m.Rate15())
		fmt.Fprintf(w, "%s.rate.mean:%f|g\n", head, m.RateMean())
	case metrics.Timer:
		m := metric.Snapshot()
		ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s.count:%d|c\n", head, m.Count())
		fmt.Fprintf(w, "%s.max:%d|g\n", head, m.Max())
		fmt.Fprintf(w, "%s.mean:%f|g\n", head, m.Mean())
		fmt.Fprintf(w, "%s.min:%d|g\n", head, m.Min())
		fmt.Fprintf(w, "%s.percentile.mean:%f|g\n", head, ps[0])
		fmt.Fprintf(w, "%s.percentile.75:%f|g\n", head, ps[1])
		fmt.Fprintf(w, "%s.percentile.95:%f|g\n", head, ps[2])
		fmt.Fprintf(w, "%s.percentile.99.0:%f|g\n", head, ps[3])
		fmt.Fprintf(w, "%s.percentile.99.9:%f|g\n", head, ps[4])
		fmt.Fprintf(w, "%s.rate.1min:%f|g\n", head, m.Rate1())
		fmt.Fprintf(w, "%s.rate.5min:%f|g\n", head, m.Rate5())
		fmt.Fprintf(w, "%s.rate.15min:%f|g\n", head, m.Rate15())
		fmt.Fprintf(w, "%s.rate.mean:%f|g\n", head, m.RateMean())
		fmt.Fprintf(w, "%s.stddev:%f|g\n", head, m.StdDev())
		fmt.Fprintf(w, "%s.sum:%d|c", head, m.Sum())
		fmt.Fprintf(w, "%s.variance:%f|c", head, m.Variance())
	case metrics.Histogram:
		m := metric.Snapshot()
		ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s.count:%d|c\n", head, m.Count())
		fmt.Fprintf(w, "%s.max:%d|g\n", head, m.Max())
		fmt.Fprintf(w, "%s.mean:%f|g\n", head, m.Mean())
		fmt.Fprintf(w, "%s.min:%d|g\n", head, m.Min())
		fmt.Fprintf(w, "%s.percentile.mean:%f|g\n", head, ps[0])
		fmt.Fprintf(w, "%s.percentile.75:%f|g\n", head, ps[1])
		fmt.Fprintf(w, "%s.percentile.95:%f|g\n", head, ps[2])
		fmt.Fprintf(w, "%s.percentile.99.0:%f|g\n", head, ps[3])
		fmt.Fprintf(w, "%s.percentile.99.9:%f|g\n", head, ps[4])
		fmt.Fprintf(w, "%s.stddev:%f|g\n", head, m.StdDev())
		fmt.Fprintf(w, "%s.sum:%d|c", head, m.Sum())
		fmt.Fprintf(w, "%s.variance:%f|c", head, m.Variance())
	}
}
