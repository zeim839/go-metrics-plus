package graphite

import (
	"bufio"
	"github.com/zeim839/go-metrics-plus"
	"github.com/zeim839/go-metrics-plus/logging"
	"net"
	"time"
)

// Config provides a container with configuration parameters for
// the Graphite exporter.
type Config struct {
	Addr          *net.TCPAddr     // Network address to connect to.
	Registry      metrics.Registry // Registry to be exported.
	FlushInterval time.Duration    // Flush interval.
	DurationUnit  time.Duration    // Time conversion unit for durations.
	Prefix        string           // Prefix to be prepended to metric names.
}

// Graphite is a blocking exporter function which reports metrics in r
// to a graphite server located at addr, flushing them every d duration
// and prepending metric names with prefix.
func Graphite(r metrics.Registry, d time.Duration, prefix string, addr *net.TCPAddr) {
	WithConfig(Config{
		Addr:          addr,
		Registry:      r,
		FlushInterval: d,
		DurationUnit:  time.Nanosecond,
		Prefix:        prefix,
	})
}

// WithConfig is a blocking exporter function just like Graphite,
// but it takes a GraphiteConfig instead. Returns a non-nil error
// on failed connections.
func WithConfig(c Config) error {
	conn, err := net.DialTCP("tcp", nil, c.Addr)
	if nil != err {
		return err
	}
	defer conn.Close()
	w := bufio.NewWriter(conn)
	//lint:ignore SA1015 TODO
	for range time.Tick(c.FlushInterval) {
		graphite(w, &c)
	}
	return nil
}

// Once performs a single submission to Graphite, returning a
// non-nil error on failed connections. This can be used in a loop
// similar to GraphiteWithConfig for custom error handling.
func Once(c Config) error {
	conn, err := net.DialTCP("tcp", nil, c.Addr)
	if nil != err {
		return err
	}
	defer conn.Close()
	w := bufio.NewWriter(conn)
	graphite(w, &c)
	return nil
}

func graphite(w *bufio.Writer, c *Config) {
	c.Registry.Each(func(name string, i interface{}) {
		logging.EncodeGraphite(w, name, c.Prefix, i)
		w.Flush()
	})
}
