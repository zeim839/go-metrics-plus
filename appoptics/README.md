# AppOptics

AppOptics is the AppOptics driver for [go-metrics-plus](https://github.com/zeim839/go-metrics-plus). It collects metrics from a registry and periodically posts them to AppOptics. The driver has been ported to go-metrics-plus from the AppOptics' [go-metrics-appoptics](https://github.com/appoptics/go-metrics-appoptics) library.

## Usage

```go
import (
	"github.com/zeim839/go-metrics-plus"
	"github.com/zeim839/go-metrics-plus/appoptics"
	"time"
)

func main() {
	metrics.GetOrRegisterCounter("myCounter", nil, metrics.Labels{"foo":"bar"})
	metrics.GetOrRegisterMeter("myMeter", nil)
	// ...

	go appoptics.AppOptics(metrics.DefaultRegistry, time.Second, "access-token",
		metrics.Labels{"hostname": "localhost"}, []float64{0.5, 0.75, 0.95, 0.99},
		time.Millisecond, "myservice.", nil, appoptics.DefaultMeasurementsURI)
}
```
