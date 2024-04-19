package logging

import (
	"fmt"
	"github.com/zeim839/go-metrics-plus"
	"io"
	"time"
)

// EncodeGraphite encodes a metric into graphite format. Some interfaces
// are encoded as multi-line summaries. Healthchecks are not supported.
func EncodeGraphite(w io.Writer, name, prefix string, i interface{}) {
	if prefix != "" {
		prefix = prefix + "."
	}
	head := prefix + name
	ts := time.Now().UTC().Unix()

	switch metric := i.(type) {
	case metrics.Counter:
		fmt.Fprintf(w, "%s %d %d\n", head, metric.Count(), ts)
	case metrics.Gauge:
		fmt.Fprintf(w, "%s %d %d\n", head, metric.Value(), ts)
	case metrics.GaugeFloat64:
		fmt.Fprintf(w, "%s %f %d\n", head, metric.Value(), ts)
	case metrics.Meter:
		m := metric.Snapshot()
		fmt.Fprintf(w, "%s.count %d %d\n", head, m.Count(), ts)
		fmt.Fprintf(w, "%s.rate.1min %f %d\n", head, m.Rate1(), ts)
		fmt.Fprintf(w, "%s.rate.5min %f %d\n", head, m.Rate5(), ts)
		fmt.Fprintf(w, "%s.rate.15min %f %d\n", head, m.Rate15(), ts)
		fmt.Fprintf(w, "%s.rate.mean %f %d\n", head, m.RateMean(), ts)
	case metrics.Timer:
		t := metric.Snapshot()
		ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s.count %d %d\n", head, t.Count(), ts)
		fmt.Fprintf(w, "%s.min %d %d\n", head, t.Min(), ts)
		fmt.Fprintf(w, "%s.max %d %d\n", head, t.Max(), ts)
		fmt.Fprintf(w, "%s.mean %f %d\n", head, t.Mean(), ts)
		fmt.Fprintf(w, "%s.sum %d %d\n", head, t.Sum(), ts)
		fmt.Fprintf(w, "%s.stddev %f %d\n", head, t.StdDev(), ts)
		fmt.Fprintf(w, "%s.variance %f %d\n", head, t.Variance(), ts)
		fmt.Fprintf(w, "%s.median %f %d\n", head, ps[0], ts)
		fmt.Fprintf(w, "%s.percentile.75 %f %d\n", head, ps[1], ts)
		fmt.Fprintf(w, "%s.percentile.95 %f %d\n", head, ps[2], ts)
		fmt.Fprintf(w, "%s.percentile.99.0 %f %d\n", head, ps[3], ts)
		fmt.Fprintf(w, "%s.percentile.99.9 %f %d\n", head, ps[4], ts)
		fmt.Fprintf(w, "%s.rate.1min %f %d\n", head, t.Rate1(), ts)
		fmt.Fprintf(w, "%s.rate.5min %f %d\n", head, t.Rate5(), ts)
		fmt.Fprintf(w, "%s.rate.15min %f %d\n", head, t.Rate15(), ts)
		fmt.Fprintf(w, "%s.rate.mean %f %d\n", head, t.RateMean(), ts)
	case metrics.Histogram:
		h := metric.Snapshot()
		ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s.count %d %v\n", head, h.Count(), ts)
		fmt.Fprintf(w, "%s.min %d %v\n", head, h.Min(), ts)
		fmt.Fprintf(w, "%s.max %d %v\n", head, h.Max(), ts)
		fmt.Fprintf(w, "%s.mean %f %v\n", head, h.Mean(), ts)
		fmt.Fprintf(w, "%s.sum %d %d\n", head, h.Sum(), ts)
		fmt.Fprintf(w, "%s.stddev %f %v\n", head, h.StdDev(), ts)
		fmt.Fprintf(w, "%s.variance %f %d\n", head, h.Variance(), ts)
		fmt.Fprintf(w, "%s.median %f %v\n", head, ps[0], ts)
		fmt.Fprintf(w, "%s.percentile.75 %f %v\n", head, ps[1], ts)
		fmt.Fprintf(w, "%s.percentile.95 %f %v\n", head, ps[2], ts)
		fmt.Fprintf(w, "%s.percentile.99.0 %f %v\n", head, ps[3], ts)
		fmt.Fprintf(w, "%s.percentile.99.9 %f %v\n", head, ps[4], ts)
	}
}
