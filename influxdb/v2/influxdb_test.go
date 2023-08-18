package v2

import (
	influx "github.com/influxdata/influxdb-client-go"
	"github.com/zeim839/go-metrics-plus"
	"time"
)

func ExampleInfluxDBV2() {
	// Register some metrics.
	m := metrics.GetOrRegisterMeter("myMeter", nil, metrics.Labels{"foo": "bar"})
	t := metrics.GetOrRegisterTimer("myTimer", nil, nil)
	m.Mark(100)
	t.Update(30 * time.Second)

	// Set up Influx client.
	token := "my_token"
	client := influx.NewClient("http://localhost:8086", token)

	// Start flushing ever 1 second.
	go InfluxDBV2(metrics.DefaultRegistry, time.Second, "prefix",
		"myBucket", "myOrg", client)
}
