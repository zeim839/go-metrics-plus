package graphite

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

func ExampleGraphite() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go Graphite(metrics.DefaultRegistry, 1*time.Second, "some.prefix", addr)
}

func ExampleWithConfig() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go WithConfig(Config{
		Addr:          addr,
		Registry:      metrics.DefaultRegistry,
		FlushInterval: 1 * time.Second,
		DurationUnit:  time.Millisecond,
	})
}

func newTestServer(t *testing.T, ctx *atomic.Bool) (map[string]string,
	net.Listener, Config, *sync.WaitGroup) {

	res := make(map[string]string)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
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
				parts := strings.Split(line, ";")
				res[parts[0]] = res[parts[0]] + line
				line, err = r.ReadString('\n')
			}
			wg.Done()
			conn.Close()
		}
	}()

	c := Config{
		Addr:          ln.Addr().(*net.TCPAddr),
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

	metrics.GetOrRegisterCounter("foo", nil, metrics.Labels{"key1": "value1"}).Inc(2)
	metrics.GetOrRegisterMeter("bar", nil, metrics.Labels{"key2": "value2"}).Mark(1)

	ctx.Store(false)
	wg.Add(1)
	err := Once(c)
	if err != nil {
		t.Errorf("Once(): %s", err)
		return
	}
	wg.Wait()

	expect := "p.bar.count;key2=value2 1"
	if str := res["p.bar.count"]; expect != str[:len(str)-12] {
		t.Errorf("%s != %s", expect, str[:len(str)-12])
	}

	expect = "p.bar.rate.15min;key2=value2 1.000000"
	if str := res["p.bar.rate.15min"]; expect != str[:len(str)-12] {
		t.Errorf("%s != %s", expect, str[:len(str)-12])
	}

	expect = "p.bar.rate.1min;key2=value2 1.000000"
	if str := res["p.bar.rate.1min"]; expect != str[:len(str)-12] {
		t.Errorf("%s != %s", expect, str[:len(str)-12])
	}

	expect = "p.bar.rate.5min;key2=value2 1.000000"
	if str := res["p.bar.rate.5min"]; expect != str[:len(str)-12] {
		t.Errorf("%s != %s", expect, str[:len(str)-12])
	}

	// p.bar.rate.mean changes every nanosecond. It is too erratic,
	// so it is ignored in this test.

	expect = "p.foo;key1=value1 2"
	if str := res["p.foo"]; expect != str[:len(str)-12] {
		t.Errorf("%s != %s", expect, str[:len(str)-12])
	}
}
