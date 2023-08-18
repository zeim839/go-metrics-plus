# InfluxDB
InfluxDB is the InfluxDB driver for [go-metrics-plus](https://github.com/zeim839/go-metrics-plus), with support for both V1 and V2 InfluxDB releases. The InfluxDBv1 and InfluxDBv2 drivers can be found in the v1 and v2 directories, respectively.

## V2 - Example

```go
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
```

## V1 - Example

```go
package v1

import (
	"os"
	client "github.com/influxdata/influxdb1-client"
	"github.com/zeim839/go-metrics-plus"
	"log"
	"net/url"
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

```
