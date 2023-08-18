package v1

import (
	client "github.com/influxdata/influxdb1-client"
	"github.com/zeim839/go-metrics-plus"
	"log"
	"net/url"
	"os"
	"time"
)

func ExampleInfluxDBV1() {
	m := metrics.GetOrRegisterMeter("myMeter", nil, metrics.Labels{"foo": "bar"})
	t := metrics.GetOrRegisterTimer("myTimer", nil, nil)
	m.Mark(100)
	t.Update(30 * time.Second)

	// Set up InfluxDB client.
	host, err := url.Parse("http://localhost:8086")
	if err != nil {
		log.Fatal("Bad Host URL")
	}

	conn, err := client.NewClient(client.Config{
		URL:      *host,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	})
	if err != nil {
		log.Fatal("Client could not connect to InfluxDB")
	}

	// Start flushing every 1 second.
	go InfluxDBV1(metrics.DefaultRegistry, time.Second, "prefix", "dummy", conn)
}
