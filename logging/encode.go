package logging

import (
	"fmt"
	"github.com/zeim839/go-metrics-plus"
	"io"
	"time"
)

// An Encoder encodes an interface and writes into the writer w.
type Encoder func(w io.Writer, name string, prefix string, i interface{})

// Encode encodes a metric into prometheus expositional format. Some interfaces
// are encoded as multi-line summaries. Healthchecks are not supported.
func Encode(w io.Writer, name, prefix string, i interface{}) {
	if prefix != "" {
		prefix = prefix + "_"
	}
	head := prefix + name
	ts := time.Now().UTC().Unix()

	switch metric := i.(type) {
	case metrics.Counter:
		fmt.Fprintf(w, "%s %d %v\n", head, metric.Count(), ts)
	case metrics.Gauge:
		fmt.Fprintf(w, "%s %d %v\n", head, metric.Value(), ts)
	case metrics.GaugeFloat64:
		fmt.Fprintf(w, "%s %f %v\n", head, metric.Value(), ts)
	case metrics.Meter:
		m := metric.Snapshot()
		fmt.Fprintf(w, "%s_count %d %v\n", head, m.Count(), ts)
		fmt.Fprintf(w, "%s_rate_1min %f %v\n", head, m.Rate1(), ts)
		fmt.Fprintf(w, "%s_rate_5min %f %v\n", head, m.Rate5(), ts)
		fmt.Fprintf(w, "%s_rate_15min %f %v\n", head, m.Rate15(), ts)
		fmt.Fprintf(w, "%s_rate_mean %f %v\n", head, m.RateMean(), ts)
	case metrics.Timer:
		t := metric.Snapshot()
		ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s_count %d %v\n", head, t.Count(), ts)
		fmt.Fprintf(w, "%s_min %d %v\n", head, t.Min(), ts)
		fmt.Fprintf(w, "%s_max %d %v\n", head, t.Max(), ts)
		fmt.Fprintf(w, "%s_mean %f %v\n", head, t.Mean(), ts)
		fmt.Fprintf(w, "%s_sum %d %v\n", head, t.Sum(), ts)
		fmt.Fprintf(w, "%s_stddev %f %v\n", head, t.StdDev(), ts)
		fmt.Fprintf(w, "%s_variance %f %v\n", head, t.Variance(), ts)
		fmt.Fprintf(w, "%s_median %f %v\n", head, ps[0], ts)
		fmt.Fprintf(w, "%s_percentile_75 %f %v\n", head, ps[1], ts)
		fmt.Fprintf(w, "%s_percentile_95 %f %v\n", head, ps[2], ts)
		fmt.Fprintf(w, "%s_percentile_99_0 %f %v\n", head, ps[3], ts)
		fmt.Fprintf(w, "%s_percentile_99_9 %f %v\n", head, ps[4], ts)
		fmt.Fprintf(w, "%s_rate_1min %f %v\n", head, t.Rate1(), ts)
		fmt.Fprintf(w, "%s_rate_5min %f %v\n", head, t.Rate5(), ts)
		fmt.Fprintf(w, "%s_rate_15min %f %v\n", head, t.Rate15(), ts)
		fmt.Fprintf(w, "%s_rate_mean %f %v\n", head, t.RateMean(), ts)
	case metrics.Histogram:
		h := metric.Snapshot()
		ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s_count %d %v\n", head, h.Count(), ts)
		fmt.Fprintf(w, "%s_min %d %v\n", head, h.Min(), ts)
		fmt.Fprintf(w, "%s_max %d %v\n", head, h.Max(), ts)
		fmt.Fprintf(w, "%s_mean %f %v\n", head, h.Mean(), ts)
		fmt.Fprintf(w, "%s_sum %d %v\n", head, h.Sum(), ts)
		fmt.Fprintf(w, "%s_stddev %f %v\n", head, h.StdDev(), ts)
		fmt.Fprintf(w, "%s_variance %f %v\n", head, h.Variance(), ts)
		fmt.Fprintf(w, "%s_median %f %v\n", head, ps[0], ts)
		fmt.Fprintf(w, "%s_percentile_75 %f %v\n", head, ps[1], ts)
		fmt.Fprintf(w, "%s_percentile_95 %f %v\n", head, ps[2], ts)
		fmt.Fprintf(w, "%s_percentile_99_0 %f %v\n", head, ps[3], ts)
		fmt.Fprintf(w, "%s_percentile_99_9 %f %v\n", head, ps[4], ts)
	}
}
