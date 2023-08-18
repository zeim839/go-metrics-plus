package appoptics

import (
	"github.com/zeim839/go-metrics-plus"
	"time"
)

func ExampleAppOptics() {
	metrics.GetOrRegisterCounter("myCounter", nil, metrics.Labels{"foo": "bar"})
	metrics.GetOrRegisterMeter("myMeter", nil, nil)

	go AppOptics(metrics.DefaultRegistry, time.Second, "token",
		metrics.Labels{"hostname": "localhost"}, []float64{0.5, 0.75, 0.95, 0.99},
		time.Millisecond, "myservice.", nil, DefaultMeasurementsURI)
}
