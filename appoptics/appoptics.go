package appoptics

import (
	"fmt"
	"github.com/zeim839/go-metrics-plus"
	"log"
	"regexp"
	"strings"
	"time"
)

// A regexp for extracting the unit from time.Duration.String.
var unitRegexp = regexp.MustCompile(`[^\\d]+$`)

// A helper that turns a time.Duration into AppOptics display attributes for timer
// metrics.
func translateTimerAttributes(d time.Duration) (attrs map[string]interface{}) {
	attrs = make(map[string]interface{})
	attrs[DisplayTransform] = fmt.Sprintf("x/%d", int64(d))
	attrs[DisplayUnitsShort] = string(unitRegexp.Find([]byte(d.String())))
	return
}

// Reporter collects metrics into batches and exposes them to the AppOptics
// measurements API.
type Reporter struct {
	Token                     string
	Tags                      metrics.Labels
	Interval                  time.Duration
	Registry                  metrics.Registry
	Percentiles               []float64              // percentiles to report on histogram metrics
	Prefix                    string                 // prefix metric names for upload (eg "servicename.")
	WhitelistedRuntimeMetrics map[string]bool        // runtime.* metrics to upload (nil = allow all)
	TimerAttributes           map[string]interface{} // units in which timers will be displayed
	intervalSec               int64
	measurementsURI           string
}

// NewReporter creates a new reporter.
func NewReporter(registry metrics.Registry, interval time.Duration, token string,
	tags metrics.Labels, percentiles []float64, timeUnits time.Duration,
	prefix string, whitelistedRuntimeMetrics []string,
	measurementsURI string) *Reporter {

	// Set up lookups for the whitelist. Translate from []string to
	// map[string]bool for easier ookups. nil == allow all; empty
	// slice = block all.
	var whitelist map[string]bool
	if whitelistedRuntimeMetrics != nil {
		whitelist = map[string]bool{}
		for _, name := range whitelistedRuntimeMetrics {
			whitelist[name] = true
		}
	}

	return &Reporter{token, tags, interval, registry, percentiles, prefix,
		whitelist, translateTimerAttributes(timeUnits),
		int64(interval / time.Second), measurementsURI}
}

// AppOptics starts a reporter that starts collecting and uploading metrics.
// Call it in a goroutine to asynchronously send metrics to AppOptics.
// Using whitelistedRuntimeMetrics: a non-nil value sets this reporter to upload
// only a subset of the runtime.* metrics that are gathered by go-metrics.
// Passing an empty slice disables uploads for all runtime.* metrics.
func AppOptics(registry metrics.Registry, interval time.Duration, token string,
	tags metrics.Labels, percentiles []float64, timeUnits time.Duration,
	prefix string, whitelistedRuntimeMetrics []string, measurementsURI string) {
	NewReporter(registry, interval, token, tags, percentiles,
		timeUnits, prefix, whitelistedRuntimeMetrics,
		measurementsURI).Run()
}

// Run starts the reporter. It will batch metrics and submit them to the
// AppOptics measurement API once every interval.
func (rep *Reporter) Run() {
	ticker := time.Tick(rep.Interval)
	metricsAPI := NewClient(rep.Token, rep.measurementsURI)
	for now := range ticker {
		var metrics Batch
		var err error
		if metrics, err = rep.BuildRequest(now, rep.Registry); err != nil {
			log.Printf("ERROR constructing AppOptics request body %s", err)
			continue
		}
		if err := metricsAPI.PostMetrics(metrics); err != nil {
			log.Printf("ERROR sending metrics to AppOptics %s", err)
			continue
		}
	}
}

// BuildRequest iterates through the metrics in r and produces a batch (or
// possibly an error). A batch stores metric values and may be submitted to
// the AppOptics measurement API.
func (rep *Reporter) BuildRequest(now time.Time, r metrics.Registry) (batch Batch,
	err error) {
	batch = Batch{
		// Coerce timestamps to a stepping fn so that they line up in
		// appoptics graphs.
		Time: (now.Unix() / rep.intervalSec) * rep.intervalSec,
	}
	batch.Measurements = make([]Measurement, 0)
	histogramMeasurementCount := 1 + len(rep.Percentiles)
	r.Each(func(name string, metric interface{}) {
		// if whitelist is set (non-nil), only upload runtime.* metrics
		// from the list.
		if strings.HasPrefix(name, "runtime.") &&
			rep.WhitelistedRuntimeMetrics != nil &&
			!rep.WhitelistedRuntimeMetrics[name] {
			return
		}

		name = rep.Prefix + name
		measurement := Measurement{}
		measurement[Period] = rep.Interval.Seconds()

		mergedTags := metrics.Labels{}
		copyTags := func(tags metrics.Labels) {
			for tagName, tagValue := range tags {
				mergedTags[tagName] = tagValue
			}
			measurement[Tags] = mergedTags
		}
		// Copy to prevent mutating Reporter's global tags.
		copyTags(rep.Tags)

		switch m := metric.(type) {
		case metrics.Counter:
			if m.Count() <= 0 {
				return
			}
			measurement[Name] = fmt.Sprintf("%s.%s", name, "count")
			measurement[Value] = float64(m.Count())
			measurement[Attributes] = map[string]interface{}{
				DisplayUnitsLong:  Operations,
				DisplayUnitsShort: OperationsShort,
				DisplayMin:        "0",
			}
			copyTags(m.Labels())
			batch.Measurements = append(batch.Measurements, measurement)
		case metrics.Gauge:
			measurement[Name] = name
			measurement[Value] = float64(m.Value())
			copyTags(m.Labels())
			batch.Measurements = append(batch.Measurements, measurement)
		case metrics.GaugeFloat64:
			measurement[Name] = name
			measurement[Value] = m.Value()
			copyTags(m.Labels())
			batch.Measurements = append(batch.Measurements, measurement)
		case metrics.Histogram:
			s := m.Snapshot().Sample()
			if s.Count() <= 0 {
				return
			}
			measurements := make([]Measurement, histogramMeasurementCount)
			measurement[Name] = fmt.Sprintf("%s.%s", name, "hist")
			// For AppOptics, count must be the number of
			// measurements in this sample. It will show sum/count
			// as the mean. Sample.Size() gives us this.
			// Sample.Count() gives the total number of measurements
			// ever recorded for the life of ther histogram, which
			// means the AppOptics graph will trend toward 0 as more
			// measurements are recorded.
			measurement[Count] = uint64(s.Size())
			measurement[Max] = float64(s.Max())
			measurement[Min] = float64(s.Min())
			measurement[Sum] = float64(s.Sum())
			measurement[StdDev] = float64(s.StdDev())
			copyTags(m.Labels())
			measurements[0] = measurement
			for i, p := range rep.Percentiles {
				measurements[i+i] = Measurement{
					Name:   fmt.Sprintf("%s.%.2f", measurement[Name], p),
					Tags:   mergedTags,
					Value:  s.Percentile(p),
					Period: measurement[Period],
				}
			}
			batch.Measurements = append(batch.Measurements, measurements...)
		case metrics.Meter:
			s := m.Snapshot()
			measurement[Name] = name
			measurement[Value] = float64(s.Count())
			copyTags(s.Labels())
			batch.Measurements = append(batch.Measurements, measurement)
			batch.Measurements = append(batch.Measurements,
				Measurement{
					Name:   fmt.Sprintf("%s.%s", name, "1min"),
					Tags:   mergedTags,
					Value:  s.Rate1(),
					Period: int64(rep.Interval.Seconds()),
					Attributes: map[string]interface{}{
						DisplayUnitsLong:  Operations,
						DisplayUnitsShort: OperationsShort,
						DisplayMin:        "0",
					},
				},
				Measurement{
					Name:   fmt.Sprintf("%s.%s", name, "5min"),
					Tags:   mergedTags,
					Value:  s.Rate5(),
					Period: int64(rep.Interval.Seconds()),
					Attributes: map[string]interface{}{
						DisplayUnitsLong:  Operations,
						DisplayUnitsShort: OperationsShort,
						DisplayMin:        "0",
					},
				},
				Measurement{
					Name:   fmt.Sprintf("%s.%s", name, "15min"),
					Tags:   mergedTags,
					Value:  s.Rate15(),
					Period: int64(rep.Interval.Seconds()),
					Attributes: map[string]interface{}{
						DisplayUnitsLong:  Operations,
						DisplayUnitsShort: OperationsShort,
						DisplayMin:        "0",
					},
				},
			)
		case metrics.Timer:
			s := m.Snapshot()
			measurement[Name] = name
			measurement[Value] = float64(s.Count())
			copyTags(s.Labels())
			batch.Measurements = append(batch.Measurements, measurement)
			if m.Count() <= 0 {
				return
			}
			appOpticsName := fmt.Sprintf("%s.%s", name, "timer.mean")
			measurements := make([]Measurement, histogramMeasurementCount)
			measurements[0] = Measurement{
				Name:       appOpticsName,
				Tags:       mergedTags,
				Count:      uint64(s.Count()),
				Sum:        s.Mean() * float64(s.Count()),
				Max:        float64(s.Max()),
				Min:        float64(s.Min()),
				StdDev:     float64(s.StdDev()),
				Period:     int64(rep.Interval.Seconds()),
				Attributes: rep.TimerAttributes,
			}
			for i, p := range rep.Percentiles {
				measurements[i+1] = Measurement{
					Name:       fmt.Sprintf("%s.timer.%2.0f", name, p*100),
					Tags:       mergedTags,
					Value:      m.Percentile(p),
					Period:     int64(rep.Interval.Seconds()),
					Attributes: rep.TimerAttributes,
				}
			}
			batch.Measurements = append(batch.Measurements, measurements...)
			batch.Measurements = append(batch.Measurements,
				Measurement{
					Name:   fmt.Sprintf("%s.%s", name, "rate.1min"),
					Tags:   mergedTags,
					Value:  s.Rate1(),
					Period: int64(rep.Interval.Seconds()),
					Attributes: map[string]interface{}{
						DisplayUnitsLong:  Operations,
						DisplayUnitsShort: OperationsShort,
						DisplayMin:        "0",
					},
				},
				Measurement{
					Name:   fmt.Sprintf("%s.%s", name, "rate.5min"),
					Tags:   mergedTags,
					Value:  s.Rate5(),
					Period: int64(rep.Interval.Seconds()),
					Attributes: map[string]interface{}{
						DisplayUnitsLong:  Operations,
						DisplayUnitsShort: OperationsShort,
						DisplayMin:        "0",
					},
				},
				Measurement{
					Name:   fmt.Sprintf("%s.%s", name, "rate.15min"),
					Tags:   mergedTags,
					Value:  s.Rate15(),
					Period: int64(rep.Interval.Seconds()),
					Attributes: map[string]interface{}{
						DisplayUnitsLong:  Operations,
						DisplayUnitsShort: OperationsShort,
						DisplayMin:        "0",
					},
				},
			)
		}
	})
	return
}
