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
		labels := EncodeGraphiteLabels(metric.Labels())
		fmt.Fprintf(w, "%s%s %d %d\n", head, labels, metric.Count(), ts)
	case metrics.Gauge:
		labels := EncodeGraphiteLabels(metric.Labels())
		fmt.Fprintf(w, "%s%s %d %d\n", head, labels, metric.Value(), ts)
	case metrics.GaugeFloat64:
		labels := EncodeGraphiteLabels(metric.Labels())
		fmt.Fprintf(w, "%s%s %f %d\n", head, labels, metric.Value(), ts)
	case metrics.Meter:
		m := metric.Snapshot()
		labels := EncodeGraphiteLabels(m.Labels())
		fmt.Fprintf(w, "%s.count%s %d %d\n", head, labels, m.Count(), ts)
		fmt.Fprintf(w, "%s.rate.1min%s %f %d\n", head, labels, m.Rate1(), ts)
		fmt.Fprintf(w, "%s.rate.5min%s %f %d\n", head, labels, m.Rate5(), ts)
		fmt.Fprintf(w, "%s.rate.15min%s %f %d\n", head, labels, m.Rate15(), ts)
		fmt.Fprintf(w, "%s.rate.mean%s %f %d\n", head, labels, m.RateMean(), ts)
	case metrics.Timer:
		t := metric.Snapshot()
		labels := EncodeGraphiteLabels(t.Labels())
		ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s.count%s %d %d\n", head, labels, t.Count(), ts)
		fmt.Fprintf(w, "%s.min%s %d %d\n", head, labels, t.Min(), ts)
		fmt.Fprintf(w, "%s.max%s %d %d\n", head, labels, t.Max(), ts)
		fmt.Fprintf(w, "%s.mean%s %f %d\n", head, labels, t.Mean(), ts)
		fmt.Fprintf(w, "%s.sum%s %d %d\n", head, labels, t.Sum(), ts)
		fmt.Fprintf(w, "%s.stddev%s %f %d\n", head, labels, t.StdDev(), ts)
		fmt.Fprintf(w, "%s.variance%s %f %d\n", head, labels, t.Variance(), ts)
		fmt.Fprintf(w, "%s.median%s %f %d\n", head, labels, ps[0], ts)
		fmt.Fprintf(w, "%s.percentile.75%s %f %d\n", head, labels, ps[1], ts)
		fmt.Fprintf(w, "%s.percentile.95%s %f %d\n", head, labels, ps[2], ts)
		fmt.Fprintf(w, "%s.percentile.99.0%s %f %d\n", head, labels, ps[3], ts)
		fmt.Fprintf(w, "%s.percentile.99.9%s %f %d\n", head, labels, ps[4], ts)
		fmt.Fprintf(w, "%s.rate.1min%s %f %d\n", head, labels, t.Rate1(), ts)
		fmt.Fprintf(w, "%s.rate.5min%s %f %d\n", head, labels, t.Rate5(), ts)
		fmt.Fprintf(w, "%s.rate.15min%s %f %d\n", head, labels, t.Rate15(), ts)
		fmt.Fprintf(w, "%s.rate.mean%s %f %d\n", head, labels, t.RateMean(), ts)
	case metrics.Histogram:
		h := metric.Snapshot()
		labels := EncodeLabels(h.Labels())
		ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s.count%s %d %v\n", head, labels, h.Count(), ts)
		fmt.Fprintf(w, "%s.min%s %d %v\n", head, labels, h.Min(), ts)
		fmt.Fprintf(w, "%s.max%s %d %v\n", head, labels, h.Max(), ts)
		fmt.Fprintf(w, "%s.mean%s %f %v\n", head, labels, h.Mean(), ts)
		fmt.Fprintf(w, "%s.sum%s %d %d\n", head, labels, h.Sum(), ts)
		fmt.Fprintf(w, "%s.stddev%s %f %v\n", head, labels, h.StdDev(), ts)
		fmt.Fprintf(w, "%s.variance%s %f %d\n", head, labels, h.Variance(), ts)
		fmt.Fprintf(w, "%s.median%s %f %v\n", head, labels, ps[0], ts)
		fmt.Fprintf(w, "%s.percentile.75%s %f %v\n", head, labels, ps[1], ts)
		fmt.Fprintf(w, "%s.percentile.95%s %f %v\n", head, labels, ps[2], ts)
		fmt.Fprintf(w, "%s.percentile.99.0%s %f %v\n", head, labels, ps[3], ts)
		fmt.Fprintf(w, "%s.percentile.99.9%s %f %v\n", head, labels, ps[4], ts)
	}
}

// EncodeGraphiteLabels encodes a series of labels into a string label format
// that's recognized by a graphite server.
func EncodeGraphiteLabels(labels metrics.Labels) string {
	if labels == nil || len(labels) < 1 {
		return ""
	}
	str := ";"
	for k, v := range labels {
		str += k + "=" + v + ";"
	}
	// Remove last semicolon character.
	return str[:len(str)-1]
}
