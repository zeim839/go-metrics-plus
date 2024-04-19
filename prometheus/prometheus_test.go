package prometheusmetrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeim839/go-metrics-plus"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func ExamplePrometheus() {
	// Register some metrics.
	metrics.GetOrRegisterTimer("myTimer", nil).Update(time.Second)
	metrics.GetOrRegisterCounter("myCounter", nil).Inc(50)
	metrics.GetOrRegisterMeter("myMeter", nil).Mark(10)

	// Create prometheus driver.
	r := prometheus.NewRegistry()
	pr, err := New(metrics.DefaultRegistry, time.Second,
		"namespace", "subsystem", r)

	if err != nil {
		panic(err)
	}

	// Flush every 1 second.pp
	go pr.Run()

	// Expose metrics to prometheus scraper on /metrics route.
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{Registry: r}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func BenchmarkPrometheus(b *testing.B) {
	s := metrics.NewUniformSample(100)
	metricRegistry := metrics.NewRegistry()
	metrics.GetOrRegisterMeter("myMeter", metricRegistry).Mark(420)
	metrics.GetOrRegisterHistogram("myHist", metricRegistry, s).Update(33)

	r := prometheus.NewRegistry()
	pr, _ := New(metricRegistry, time.Nanosecond, "ns", "ss", r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pr.Once()
	}
}

// Exercise race detector.
func TestPrometheusConcurrency(t *testing.T) {
	s := metrics.NewUniformSample(100)
	metricRegistry := metrics.NewRegistry()
	metrics.GetOrRegisterMeter("myMeter", metricRegistry).Mark(420)
	metrics.GetOrRegisterHistogram("myHist", metricRegistry, s).Update(33)

	r := prometheus.NewRegistry()
	pr, err := New(metricRegistry, time.Nanosecond, "ns", "ss", r)
	if err != nil {
		t.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(pr *Prometheus, wg *sync.WaitGroup) {
			pr.Once()
			wg.Done()
		}(pr, wg)
		wg.Wait()
	}
}

func TestPrometheusCreate(t *testing.T) {
	reg := metrics.NewRegistry()

	// Prometheus registry should not be nil.
	_, err := New(reg, time.Second, "", "", nil)
	if err == nil {
		t.Error("New(): created driver with nil prometheus registry")
	}

	// Should not error if non-nil.
	_, err = New(reg, time.Second, "", "", prometheus.NewRegistry())
	if err != nil {
		t.Errorf("New(): failed with error %s", err)
	}
}

func TestPrometheusOnce(t *testing.T) {
	metrics.GetOrRegisterCounter("counter", nil).Inc(45)
	metrics.GetOrRegisterGauge("gauge", nil).Update(45)
	metrics.GetOrRegisterMeter("meter", nil).Mark(45)

	r := prometheus.NewRegistry()
	pr, err := New(metrics.DefaultRegistry, time.Second, "", "", r)
	if err != nil {
		t.Errorf("Once(): failed with error %s", err)
	}

	pr.Once()
	metrics, _ := r.Gather()
	if len(metrics) != 7 {
		t.Errorf("Once(): expected 36 metrics to be registered but found %d", len(metrics))
	}

	// Counter
	expected := "name:\"counter_count\" help:\"counter_count\" type:GAUGE " +
		"metric:<gauge:<value:45 > > "
	if expected != fmt.Sprint(metrics[0]) {
		t.Errorf("Once(): %s != %s", expected, metrics[0])
	}

	// GaugeFloat64
	expected = "name:\"gauge_gauge\" help:\"gauge_gauge\" type:GAUGE " +
		"metric:<gauge:<value:45 > > "
	if expected != fmt.Sprint(metrics[1]) {
		t.Errorf("Once(): %s != %s", expected, metrics[1])
	}

	// Meter
	expected = "name:\"meter_count\" help:\"meter_count\" type:GAUGE " +
		"metric:<gauge:<value:45 > > "
	if expected != fmt.Sprint(metrics[2]) {
		t.Errorf("Once(): %s != %s", expected, metrics[2])
	}

	expected = "name:\"meter_rate_15min\" help:\"meter_rate_15min\" type:GAUGE " +
		"metric:<gauge:<value:45 > > "
	if expected != fmt.Sprint(metrics[3]) {
		t.Errorf("Once(): %s != %s", expected, metrics[3])
	}

	expected = "name:\"meter_rate_1min\" help:\"meter_rate_1min\" type:GAUGE " +
		"metric:<gauge:<value:45 > > "
	if expected != fmt.Sprint(metrics[4]) {
		t.Errorf("Once(): %s != %s", expected, metrics[4])
	}

	expected = "name:\"meter_rate_5min\" help:\"meter_rate_5min\" type:GAUGE " +
		"metric:<gauge:<value:45 > > "
	if expected != fmt.Sprint(metrics[5]) {
		t.Errorf("Once(): %s != %s", expected, metrics[5])
	}
	// Mean rate is too volatile to calculate.
}
