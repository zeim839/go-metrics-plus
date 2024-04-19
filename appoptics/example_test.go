package appoptics

import (
	"github.com/zeim839/go-metrics-plus"
	"time"
)

func ExampleAppOptics() {
	metrics.GetOrRegisterCounter("myCounter", nil)
	metrics.GetOrRegisterMeter("myMeter", nil)

	go AppOptics(metrics.DefaultRegistry, time.Second, "token",
		[]float64{0.5, 0.75, 0.95, 0.99},
		time.Millisecond, "myservice.", nil, DefaultMeasurementsURI)
}
