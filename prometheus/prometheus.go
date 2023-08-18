package prometheusmetrics

import (
	"fmt"
	pr "github.com/prometheus/client_golang/prometheus"
	"github.com/zeim839/go-metrics-plus"
	"log"
	"sync"
	"time"
)

// Config provides a container with configuration parameters for Prometheus
// exposer. Each metric's name will be prepended by namespace and subsystem,
// like so: namespace_subsystem_myMetric.
type Config struct {
	Namespace     string           // The Prometheus namespace.
	Subsystem     string           // The Prometheus subsystem.
	Registry      metrics.Registry // Registry to be exported.
	FlushInterval time.Duration    // Flush interval.
	DurationUnit  time.Duration    // Time conversion unit for durations.
}

// The Prometheus exposer's state. Can be created with New() or NewWithConfig().
type Prometheus struct {
	config  Config
	reg     *pr.Registry
	vectors map[string]*pr.GaugeVec
	mutex   sync.Mutex
}

// New creates a new prometheus exposer instance that will expose metrics registry
// 'r' every 'd' nanoseconds, using namespace 'ns', subsystem 'ss', and will
// write results to prometheus registry 'p'. Returns an error if 'p' is nil.
func New(r metrics.Registry, d time.Duration, ns, ss string,
	p *pr.Registry) (*Prometheus, error) {
	return NewWithConfig(Config{
		Namespace:     ns,
		Subsystem:     ss,
		Registry:      r,
		FlushInterval: d,
		DurationUnit:  time.Nanosecond,
	}, p)
}

// NewWithConfig creates a new prometheus exposer instance using config 'c' and
// prometheus registry 'p'. Returns an error if 'p' is nil.
func NewWithConfig(c Config, p *pr.Registry) (*Prometheus, error) {
	if p == nil {
		return nil, fmt.Errorf("Prometheus registry cannot be nil")
	}
	prom := &Prometheus{
		vectors: map[string]*pr.GaugeVec{},
		config:  c,
		reg:     p,
	}
	return prom, nil
}

// Retrieves or creates a gauge vector for the given name and label set. Not
// threadsafe, must be called with a mutex.
func (p *Prometheus) getVector(name string, labels metrics.Labels) *pr.GaugeVec {
	vec, ok := p.vectors[name]
	if !ok {
		// Get label keys.
		keys := []string{}
		for k := range labels {
			keys = append(keys, k)
		}
		vec = pr.NewGaugeVec(pr.GaugeOpts{
			Namespace: p.config.Namespace,
			Subsystem: p.config.Subsystem,
			Name:      name,
			Help:      name,
		}, keys)

		p.vectors[name] = vec
		p.reg.MustRegister(vec)
	}
	return vec
}

// Sets the value of the gauge with name and labels. Prints an error if the
// gauge cannot be set or labels dont match predefined schema. Not threadsafe,
// must be called with a mutex.
func (p *Prometheus) setValue(name string, val float64, labels metrics.Labels) {
	vec := p.getVector(name, labels)
	gauge, err := vec.GetMetricWith(pr.Labels(labels))
	if err != nil {
		log.Printf("Error: (metrics) ignoring %s due to error: %s", name, err)
		return
	}
	gauge.Set(val)
}

// Once performs a single submission of metrics to the configured prometheus
// registry.
func (p *Prometheus) Once() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.config.Registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			m := metric.Snapshot()
			p.setValue(name+"_count", float64(m.Count()), m.Labels())
		case metrics.Gauge:
			m := metric.Snapshot()
			p.setValue(name+"_gauge", float64(m.Value()), m.Labels())
		case metrics.GaugeFloat64:
			m := metric.Snapshot()
			p.setValue(name+"_gauge", m.Value(), m.Labels())
		case metrics.Meter:
			m := metric.Snapshot()
			labels := m.Labels()
			p.setValue(name+"_count", float64(m.Count()), labels)
			p.setValue(name+"_rate_1min", m.Rate1(), labels)
			p.setValue(name+"_rate_5min", m.Rate5(), labels)
			p.setValue(name+"_rate_15min", m.Rate15(), labels)
			p.setValue(name+"_rate_mean", m.RateMean(), labels)
		case metrics.Timer:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			labels := m.Labels()
			p.setValue(name+"_count", float64(m.Count()), labels)
			p.setValue(name+"_min", float64(m.Min()), labels)
			p.setValue(name+"_max", float64(m.Max()), labels)
			p.setValue(name+"_mean", m.Mean(), labels)
			p.setValue(name+"_sum", float64(m.Sum()), labels)
			p.setValue(name+"_variance", m.Variance(), labels)
			p.setValue(name+"_stddev", m.StdDev(), labels)
			p.setValue(name+"_median", ps[0], labels)
			p.setValue(name+"_percentile_75", ps[1], labels)
			p.setValue(name+"_percentile_95", ps[2], labels)
			p.setValue(name+"_percentile_99_0", ps[3], labels)
			p.setValue(name+"_percentile_99_9", ps[4], labels)
			p.setValue(name+"_rate_1min", m.Rate1(), labels)
			p.setValue(name+"_rate_5min", m.Rate5(), labels)
			p.setValue(name+"_rate_15min", m.Rate15(), labels)
			p.setValue(name+"_rate_mean", m.RateMean(), labels)
		case metrics.Histogram:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			labels := m.Labels()
			p.setValue(name+"_count", float64(m.Count()), labels)
			p.setValue(name+"_min", float64(m.Min()), labels)
			p.setValue(name+"_max", float64(m.Max()), labels)
			p.setValue(name+"_mean", m.Mean(), labels)
			p.setValue(name+"_sum", float64(m.Sum()), labels)
			p.setValue(name+"_variance", m.Variance(), labels)
			p.setValue(name+"_stddev", m.StdDev(), labels)
			p.setValue(name+"_median", ps[0], labels)
			p.setValue(name+"_percentile_75", ps[1], labels)
			p.setValue(name+"_percentile_95", ps[2], labels)
			p.setValue(name+"_percentile_99_0", ps[3], labels)
			p.setValue(name+"_percentile_99_9", ps[4], labels)
		}
	})
}

// Run performs a submission of all metrics data into the configured prometheus
// registry once every FlushInterval.
func (p *Prometheus) Run() {
	for range time.Tick(p.config.FlushInterval) {
		p.Once()
	}
}
