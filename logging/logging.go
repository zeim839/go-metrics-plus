package logging

import (
	"github.com/zeim839/go-metrics-plus"
	"os"
	"time"
)

// Logger is a block exporter function which flushes metrics in r to stdout
// using the given Encoder, sinking them every d duration and prepending
// metric names with prefix.
func Logger(f Encoder, r metrics.Registry, d time.Duration, prefix string) {
	for range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			f(os.Stdout, name, prefix, i)
		})
	}
}
