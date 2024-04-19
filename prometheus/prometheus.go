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
func (p *Prometheus) getVector(name string) *pr.GaugeVec {
	vec, ok := p.vectors[name]
	if !ok {
		vec = pr.NewGaugeVec(pr.GaugeOpts{
			Namespace: p.config.Namespace,
			Subsystem: p.config.Subsystem,
			Name:      name,
			Help:      name,
		}, nil)
		p.vectors[name] = vec
		p.reg.MustRegister(vec)
	}
	return vec
}

// Sets the value of the gauge with name and labels. Prints an error if the
// gauge cannot be set or labels dont match predefined schema. Not threadsafe,
// must be called with a mutex.
func (p *Prometheus) setValue(name string, val float64) {
	vec := p.getVector(name)
	gauge, err := vec.GetMetricWith(nil)
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
			p.setValue(name+"_count", float64(m.Count()))
		case metrics.Gauge:
			m := metric.Snapshot()
			p.setValue(name+"_gauge", float64(m.Value()))
		case metrics.GaugeFloat64:
			m := metric.Snapshot()
			p.setValue(name+"_gauge", m.Value())
		case metrics.Meter:
			m := metric.Snapshot()
			p.setValue(name+"_count", float64(m.Count()))
			p.setValue(name+"_rate_1min", m.Rate1())
			p.setValue(name+"_rate_5min", m.Rate5())
			p.setValue(name+"_rate_15min", m.Rate15())
			p.setValue(name+"_rate_mean", m.RateMean())
		case metrics.Timer:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			p.setValue(name+"_count", float64(m.Count()))
			p.setValue(name+"_min", float64(m.Min()))
			p.setValue(name+"_max", float64(m.Max()))
			p.setValue(name+"_mean", m.Mean())
			p.setValue(name+"_sum", float64(m.Sum()))
			p.setValue(name+"_variance", m.Variance())
			p.setValue(name+"_stddev", m.StdDev())
			p.setValue(name+"_median", ps[0])
			p.setValue(name+"_percentile_75", ps[1])
			p.setValue(name+"_percentile_95", ps[2])
			p.setValue(name+"_percentile_99_0", ps[3])
			p.setValue(name+"_percentile_99_9", ps[4])
			p.setValue(name+"_rate_1min", m.Rate1())
			p.setValue(name+"_rate_5min", m.Rate5())
			p.setValue(name+"_rate_15min", m.Rate15())
			p.setValue(name+"_rate_mean", m.RateMean())
		case metrics.Histogram:
			m := metric.Snapshot()
			ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			p.setValue(name+"_count", float64(m.Count()))
			p.setValue(name+"_min", float64(m.Min()))
			p.setValue(name+"_max", float64(m.Max()))
			p.setValue(name+"_mean", m.Mean())
			p.setValue(name+"_sum", float64(m.Sum()))
			p.setValue(name+"_variance", m.Variance())
			p.setValue(name+"_stddev", m.StdDev())
			p.setValue(name+"_median", ps[0])
			p.setValue(name+"_percentile_75", ps[1])
			p.setValue(name+"_percentile_95", ps[2])
			p.setValue(name+"_percentile_99_0", ps[3])
			p.setValue(name+"_percentile_99_9", ps[4])
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
