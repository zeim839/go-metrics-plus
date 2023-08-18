package statsd

import (
	"bufio"
	"fmt"
	"github.com/zeim839/go-metrics-plus"
	"github.com/zeim839/go-metrics-plus/logging"
	"net"
	"time"
)

// Config provides a container with configuration parameters for
// the Statsd exporter.
type Config struct {
	Addr          string           // Network address to connect to.
	Protocol      string           // Statsd server's network protocol.
	Registry      metrics.Registry // Registry to be exported.
	FlushInterval time.Duration    // Flush Interval.
	DurationUnit  time.Duration    // Time conversion unit for durations.
	Prefix        string           // Prefix to be prepended to metric names.
	Timeout       time.Duration    // How long to wait for a connection to establish.
}

// Statsd is an exporter function which reports metrics in r to a Statsd server
// located at addr and Protocol, flushing them every d duration and prepending
// metric names with prefix. It is assumed that d is equivalent to the flush
// interval implemented by the Statsd server.
func Statsd(r metrics.Registry, d time.Duration, prefix, addr, protocol string) {
	WithConfig(Config{
		Addr:          addr,
		Protocol:      protocol,
		Registry:      r,
		FlushInterval: d,
		DurationUnit:  time.Nanosecond,
		Prefix:        prefix,
		Timeout:       250 * time.Millisecond,
	})
}

// Naive/basic validation to prevent client from hanging on bad addresses.
func checkConfig(c Config) error {
	switch c.Protocol {
	case "tcp", "tcp4", "tcp6":
		_, err := net.ResolveTCPAddr(c.Protocol, c.Addr)
		return err
	case "udp", "udp4", "udp6":
		_, err := net.ResolveUDPAddr(c.Protocol, c.Addr)
		return err
	default:
		return fmt.Errorf("unsupported protocol %s", c.Protocol)
	}
}

// WithConfig is an exporter function just like Statsd, but it takes a Config
// instead. It is assumed that FlushInterval is equivalent to the flush interval
// implemented by the Statsd server. Returns non-nil error on failed connections.
func WithConfig(c Config) error {
	if err := checkConfig(c); err != nil {
		return err
	}
	conn, err := net.DialTimeout(c.Protocol, c.Addr, c.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	w := bufio.NewWriter(conn)
	//lint:ignore SA1015 TODO
	for range time.Tick(c.FlushInterval) {
		statsd(w, &c)
	}
	return nil
}

// Once performs a single submission to Statsd, returning a
// non-nil error on failed connections. This can be used in a loop
// similar to WithConfig for custom error handling.
func Once(c Config) error {
	if err := checkConfig(c); err != nil {
		return err
	}
	conn, err := net.DialTimeout(c.Protocol, c.Addr, c.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	w := bufio.NewWriter(conn)
	statsd(w, &c)
	return nil
}

func statsd(w *bufio.Writer, c *Config) {
	c.Registry.Each(func(name string, i interface{}) {
		logging.EncodeStatsd(w, name, c.Prefix, i)
		w.Flush()
	})
}
