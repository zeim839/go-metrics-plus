package statsd

import (
	"bufio"
	"github.com/zeim839/go-metrics-plus"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func ExampleStatsd() {
	go Statsd(metrics.DefaultRegistry, time.Second, "prefix", ":8125", "tcp")
}

func ExampleWithConfig() {
	go WithConfig(Config{
		Addr:          ":8125",
		Protocol:      "tcp",
		Registry:      metrics.DefaultRegistry,
		FlushInterval: time.Second,
		Prefix:        "some.prefix",
	})
}

func newTestServer(t *testing.T, ctx *atomic.Bool) (map[string]string,
	net.Listener, Config, *sync.WaitGroup) {

	res := make(map[string]string)
	ln, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		t.Fatal("could not start dummy server:", err)
	}

	var wg sync.WaitGroup
	go func() {
		for ctx.Load() {
			conn, err := ln.Accept()
			if err != nil {
				t.Errorf("dummy server error: %s", err)
				return
			}
			r := bufio.NewReader(conn)
			line, err := r.ReadString('\n')
			for err == nil {
				parts := strings.Split(line, ":")
				res[parts[0]] = res[parts[0]] + line
				line, err = r.ReadString('\n')
			}
			wg.Done()
			conn.Close()
		}
	}()

	c := Config{
		Addr:          ":9999",
		Protocol:      "tcp",
		Registry:      metrics.DefaultRegistry,
		FlushInterval: 10 * time.Millisecond,
		DurationUnit:  time.Millisecond,
		Prefix:        "p",
	}

	return res, ln, c, &wg
}

func TestWrites(t *testing.T) {
	var ctx atomic.Bool
	ctx.Store(true)
	res, ln, c, wg := newTestServer(t, &ctx)
	defer ln.Close()

	metrics.GetOrRegisterCounter("foo", nil, nil).Inc(2)
	metrics.GetOrRegisterMeter("bar", nil, nil).Mark(1)

	ctx.Store(false)
	wg.Add(1)
	err := Once(c)
	if err != nil {
		t.Errorf("Once(): %s", err)
		return
	}
	wg.Wait()

	expect := "p.bar.count:1|c\n"
	if str := res["p.bar.count"]; expect != str {
		t.Errorf("%v != %v", expect, str)
	}

	expect = "p.bar.rate.15min:1.000000|g\n"
	if str := res["p.bar.rate.15min"]; expect != str {
		t.Errorf("%v != %v", expect, str)
	}

	expect = "p.bar.rate.5min:1.000000|g\n"
	if str := res["p.bar.rate.5min"]; expect != str {
		t.Errorf("%v != %v", expect, str)
	}

	expect = "p.bar.rate.1min:1.000000|g\n"
	if str := res["p.bar.rate.1min"]; expect != str {
		t.Errorf("%v != %v", expect, str)
	}

	expect = "p.bar.rate.1min:1.000000|g\n"
	if str := res["p.bar.rate.1min"]; expect != str {
		t.Errorf("%v != %v", expect, str)
	}

	expect = "p.foo:2|c\n"
	if str := res["p.foo"]; expect != str {
		t.Errorf("%v != %v", expect, str)
	}
}
