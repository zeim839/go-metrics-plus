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
		fmt.Fprintf(w, "%s%s %d %v\n", head, EncodeLabels(metric.Labels()),
			metric.Count(), ts)
	case metrics.Gauge:
		fmt.Fprintf(w, "%s%s %d %v\n", head, EncodeLabels(metric.Labels()),
			metric.Value(), ts)
	case metrics.GaugeFloat64:
		fmt.Fprintf(w, "%s%s %f %v\n", head, EncodeLabels(metric.Labels()),
			metric.Value(), ts)
	case metrics.Meter:
		m := metric.Snapshot()
		labels := EncodeLabels(m.Labels())
		fmt.Fprintf(w, "%s_count%s %d %v\n", head, labels, m.Count(), ts)
		fmt.Fprintf(w, "%s_rate_1min%s %f %v\n", head, labels, m.Rate1(), ts)
		fmt.Fprintf(w, "%s_rate_5min%s %f %v\n", head, labels, m.Rate5(), ts)
		fmt.Fprintf(w, "%s_rate_15min%s %f %v\n", head, labels, m.Rate15(), ts)
		fmt.Fprintf(w, "%s_rate_mean%s %f %v\n", head, labels, m.RateMean(), ts)
	case metrics.Timer:
		t := metric.Snapshot()
		labels := EncodeLabels(t.Labels())
		ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s_count%s %d %v\n", head, labels, t.Count(), ts)
		fmt.Fprintf(w, "%s_min%s %d %v\n", head, labels, t.Min(), ts)
		fmt.Fprintf(w, "%s_max%s %d %v\n", head, labels, t.Max(), ts)
		fmt.Fprintf(w, "%s_mean%s %f %v\n", head, labels, t.Mean(), ts)
		fmt.Fprintf(w, "%s_sum%s %d %v\n", head, labels, t.Sum(), ts)
		fmt.Fprintf(w, "%s_stddev%s %f %v\n", head, labels, t.StdDev(), ts)
		fmt.Fprintf(w, "%s_variance%s %f %v\n", head, labels, t.Variance(), ts)
		fmt.Fprintf(w, "%s_median%s %f %v\n", head, labels, ps[0], ts)
		fmt.Fprintf(w, "%s_percentile_75%s %f %v\n", head, labels, ps[1], ts)
		fmt.Fprintf(w, "%s_percentile_95%s %f %v\n", head, labels, ps[2], ts)
		fmt.Fprintf(w, "%s_percentile_99_0%s %f %v\n", head, labels, ps[3], ts)
		fmt.Fprintf(w, "%s_percentile_99_9%s %f %v\n", head, labels, ps[4], ts)
		fmt.Fprintf(w, "%s_rate_1min%s %f %v\n", head, labels, t.Rate1(), ts)
		fmt.Fprintf(w, "%s_rate_5min%s %f %v\n", head, labels, t.Rate5(), ts)
		fmt.Fprintf(w, "%s_rate_15min%s %f %v\n", head, labels, t.Rate15(), ts)
		fmt.Fprintf(w, "%s_rate_mean%s %f %v\n", head, labels, t.RateMean(), ts)
	case metrics.Histogram:
		h := metric.Snapshot()
		labels := EncodeLabels(h.Labels())
		ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
		fmt.Fprintf(w, "%s_count%s %d %v\n", head, labels, h.Count(), ts)
		fmt.Fprintf(w, "%s_min%s %d %v\n", head, labels, h.Min(), ts)
		fmt.Fprintf(w, "%s_max%s %d %v\n", head, labels, h.Max(), ts)
		fmt.Fprintf(w, "%s_mean%s %f %v\n", head, labels, h.Mean(), ts)
		fmt.Fprintf(w, "%s_sum%s %d %v\n", head, labels, h.Sum(), ts)
		fmt.Fprintf(w, "%s_stddev%s %f %v\n", head, labels, h.StdDev(), ts)
		fmt.Fprintf(w, "%s_variance%s %f %v\n", head, labels, h.Variance(), ts)
		fmt.Fprintf(w, "%s_median%s %f %v\n", head, labels, ps[0], ts)
		fmt.Fprintf(w, "%s_percentile_75%s %f %v\n", head, labels, ps[1], ts)
		fmt.Fprintf(w, "%s_percentile_95%s %f %v\n", head, labels, ps[2], ts)
		fmt.Fprintf(w, "%s_percentile_99_0%s %f %v\n", head, labels, ps[3], ts)
		fmt.Fprintf(w, "%s_percentile_99_9%s %f %v\n", head, labels, ps[4], ts)
	}
}

// EncodeLabels encodes labels into JSON format. Returns "" if the
// slice is empty.
func EncodeLabels(labels metrics.Labels) string {
	if labels == nil || len(labels) < 1 {
		return ""
	}
	str := "{"
	for k, v := range labels {
		str += k + ":\"" + v + "\","
	}
	// Remove last comma character and add closing brace.
	return str[:len(str)-1] + "}"
}
