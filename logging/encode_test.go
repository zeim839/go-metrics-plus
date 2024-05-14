package logging

import (
	"bytes"
	"github.com/zeim839/go-metrics-plus"
	"strings"
	"testing"
	"time"
)

func BenchmarkEncode(b *testing.B) {
	// Timer is worst-case scenario (most verbose).
	timer := metrics.NewTimer()
	timer.Update(time.Second)
	b.ResetTimer()
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		Encode(buf, "foo", "bar", timer)
	}
}

func TestEncodeCounter(t *testing.T) {
	counter := metrics.NewCounter()
	counter.Inc(500)
	buf := new(bytes.Buffer)
	Encode(buf, "foo", "bar", counter)
	expect := "bar_foo 500"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without namespace.
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "", counter)
	expect = "foo 500"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without labels.
	counter = metrics.NewCounter()
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "bar", counter)
	expect = "bar_foo 0"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}
}

func TestEncodeGauge(t *testing.T) {
	gauge := metrics.NewGauge()
	gauge.Update(10)
	buf := new(bytes.Buffer)
	Encode(buf, "foo", "bar", gauge)
	expect := "bar_foo 10"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without namespace.
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "", gauge)
	expect = "foo 10"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without labels.
	buf = new(bytes.Buffer)
	gauge = metrics.NewGauge()
	Encode(buf, "foo", "bar", gauge)
	expect = "bar_foo 0"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}
}

func TestEncodeGaugeFloat64(t *testing.T) {
	gauge := metrics.NewGaugeFloat64()
	gauge.Update(10)
	buf := new(bytes.Buffer)
	Encode(buf, "foo", "bar", gauge)
	expect := "bar_foo 10.000000"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without namespace.
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "", gauge)
	expect = "foo 10.000000"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}

	// Without labels.
	gauge = metrics.NewGaugeFloat64()
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "bar", gauge)
	expect = "bar_foo 0.000000"
	if str := buf.String(); str[:len(str)-12] != expect {
		t.Errorf("Encode(): %s != %s", str[:len(str)-12], expect)
	}
}

func TestEncodeHealthcheck(t *testing.T) {
	check := metrics.NewHealthcheck(func(metrics.Healthcheck) {})
	buf := new(bytes.Buffer)
	Encode(buf, "foo", "bar", check)
	if str := buf.String(); str != "" {
		t.Errorf("Encode(): Healthcheck returned non-empty string: %s", str)
	}
}

func TestEncodeHistogram(t *testing.T) {
	hist := metrics.NewHistogram(metrics.NewUniformSample(100))
	hist.Update(100.0)
	buf := new(bytes.Buffer)
	Encode(buf, "foo", "bar", hist)
	lines := strings.Split(buf.String(), "\n")
	if len(lines) != 13 {
		t.Fatal("Encode(): Did not produce 13 lines for histogram")
	}
	expect := "bar_foo_count 1"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_min 100"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_max 100"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_mean 100.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_sum 100"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_stddev 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_variance 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_median 100.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_75 100.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_95 100.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_99_0 100.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_99_9 100.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}

	// Without namespace
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "", hist)
	lines = strings.Split(buf.String(), "\n")
	if len(lines) != 13 {
		t.Fatal("Encode(): Did not produce 13 lines for histogram")
	}
	expect = "foo_count 1"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_min 100"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_max 100"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_mean 100.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_sum 100"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_stddev 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_variance 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_median 100.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_75 100.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_95 100.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_99_0 100.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_99_9 100.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}

	// Without labels.
	hist = metrics.NewHistogram(metrics.NewUniformSample(100))
	hist.Update(100)
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "bar", hist)
	lines = strings.Split(buf.String(), "\n")
	if len(lines) != 13 {
		t.Fatal("Encode(): Did not produce 13 lines for histogram")
	}
	expect = "bar_foo_count 1"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_min 100"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_max 100"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_mean 100.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_sum 100"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_stddev 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_variance 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_median 100.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_75 100.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_95 100.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_99_0 100.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_99_9 100.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
}

func TestEncodeMeter(t *testing.T) {
	meter := metrics.NewMeter()
	meter.Mark(20)
	buf := new(bytes.Buffer)
	Encode(buf, "foo", "bar", meter)
	lines := strings.Split(buf.String(), "\n")
	if len(lines) != 6 {
		t.Fatal("Encode(): Did not produce six lines for meter")
	}
	expect := "bar_foo_count 20"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_1min 0.000000"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_5min 0.000000"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_15min 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_mean"
	if line := lines[4][:len(expect)]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}

	// Without namespace.
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "", meter)
	lines = strings.Split(buf.String(), "\n")
	if len(lines) != 6 {
		t.Fatal("Encode(): Did not produce six lines for meter")
	}
	expect = "foo_count 20"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_1min 0.000000"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_5min 0.000000"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_15min 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_mean"
	if line := lines[4][:len(expect)]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}

	// Without labels.
	meter = metrics.NewMeter()
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "bar", meter)
	lines = strings.Split(buf.String(), "\n")
	if len(lines) != 6 {
		t.Fatal("Encode(): Did not produce six lines for meter")
	}
	expect = "bar_foo_count 0"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_1min 0.000000"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_5min 0.000000"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_15min 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_mean"
	if line := lines[4][:len(expect)]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
}

func TestEncodeTimer(t *testing.T) {
	// Do not timer.Update() without some time.Sleep, results are erratic.
	timer := metrics.NewTimer()
	buf := new(bytes.Buffer)
	Encode(buf, "foo", "bar", timer)
	lines := strings.Split(buf.String(), "\n")
	if len(lines) != 17 {
		t.Fatal("Encode(): Did not produce 17 lines for timer")
	}
	expect := "bar_foo_count 0"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_min 0"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_max 0"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_mean 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_sum 0"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_stddev 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_variance 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_median 0.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_75 0.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_95 0.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_99_0 0.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_percentile_99_9 0.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_1min 0.000000"
	if line := lines[12][:len(lines[12])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_5min 0.000000"
	if line := lines[13][:len(lines[13])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_15min 0.000000"
	if line := lines[14][:len(lines[14])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "bar_foo_rate_mean 0.000000"
	if line := lines[15][:len(lines[15])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}

	// Without namespace.
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "", timer)
	lines = strings.Split(buf.String(), "\n")
	if len(lines) != 17 {
		t.Error("Encode(): Did not produce 17 lines for timer")
	}
	expect = "foo_count 0"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_min 0"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_max 0"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_mean 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_sum 0"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_stddev 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_variance 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_median 0.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_75 0.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_95 0.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_99_0 0.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_99_9 0.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_1min 0.000000"
	if line := lines[12][:len(lines[12])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_5min 0.000000"
	if line := lines[13][:len(lines[13])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_15min 0.000000"
	if line := lines[14][:len(lines[14])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_mean 0.000000"
	if line := lines[15][:len(lines[15])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}

	// Without labels.
	timer = metrics.NewTimer()
	buf = new(bytes.Buffer)
	Encode(buf, "foo", "", timer)
	lines = strings.Split(buf.String(), "\n")
	if len(lines) != 17 {
		t.Fatal("Encode(): Did not produce 17 lines for timer")
	}
	expect = "foo_count 0"
	if line := lines[0][:len(lines[0])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_min 0"
	if line := lines[1][:len(lines[1])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_max 0"
	if line := lines[2][:len(lines[2])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_mean 0.000000"
	if line := lines[3][:len(lines[3])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_sum 0"
	if line := lines[4][:len(lines[4])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_stddev 0.000000"
	if line := lines[5][:len(lines[5])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_variance 0.000000"
	if line := lines[6][:len(lines[6])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_median 0.000000"
	if line := lines[7][:len(lines[7])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_75 0.000000"
	if line := lines[8][:len(lines[8])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_95 0.000000"
	if line := lines[9][:len(lines[9])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_99_0 0.000000"
	if line := lines[10][:len(lines[10])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_percentile_99_9 0.000000"
	if line := lines[11][:len(lines[11])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_1min 0.000000"
	if line := lines[12][:len(lines[12])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_5min 0.000000"
	if line := lines[13][:len(lines[13])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_15min 0.000000"
	if line := lines[14][:len(lines[14])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
	expect = "foo_rate_mean 0.000000"
	if line := lines[15][:len(lines[15])-11]; line != expect {
		t.Errorf("Encode(): %s != %s", line, expect)
	}
}

func TestEncodeUnknown(t *testing.T) {
	srt := struct {
		a string
		b int16
	}{"asd", 123}

	buf := new(bytes.Buffer)
	Encode(buf, "foo", "bar", srt)
	if str := buf.String(); str != "" {
		t.Errorf("Encode(): Unknown struct returned non-empty string: %s", str)
	}
}
