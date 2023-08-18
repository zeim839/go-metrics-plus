package metrics

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func BenchmarkTimer(b *testing.B) {
	tm := NewTimer(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.Update(1)
	}
}

func TestGetOrRegisterTimer(t *testing.T) {
	r := NewRegistry()
	NewRegisteredTimer("foo", r, nil).Update(47)
	if tm := GetOrRegisterTimer("foo", r, nil); tm.Count() != 1 {
		t.Fatal(tm)
	}
}

func TestTimerExtremes(t *testing.T) {
	tm := NewTimer(nil)
	tm.Update(math.MaxInt64)
	tm.Update(0)
	if stdDev := tm.StdDev(); stdDev != 4.611686018427388e+18 {
		t.Errorf("tm.StdDev(): 4.611686018427388e+18 != %v\n", stdDev)
	}
}

func TestTimerFunc(t *testing.T) {
	tm := NewTimer(nil)
	tm.Time(func() { time.Sleep(50e6) })
	if max := tm.Max(); 45e6 > max || max > 55e6 {
		t.Errorf("tm.Max(): 45e6 > %v || %v > 55e6\n", max, max)
	}
}

func TestTimerZero(t *testing.T) {
	tm := NewTimer(nil)
	if count := tm.Count(); count != 0 {
		t.Errorf("tm.Count(): 0 != %v\n", count)
	}
	if min := tm.Min(); min != 0 {
		t.Errorf("tm.Min(): 0 != %v\n", min)
	}
	if max := tm.Max(); max != 0 {
		t.Errorf("tm.Max(): 0 != %v\n", max)
	}
	if mean := tm.Mean(); mean != 0.0 {
		t.Errorf("tm.Mean(): 0.0 != %v\n", mean)
	}
	if stdDev := tm.StdDev(); stdDev != 0.0 {
		t.Errorf("tm.StdDev(): 0.0 != %v\n", stdDev)
	}
	ps := tm.Percentiles([]float64{0.5, 0.75, 0.99})
	if ps[0] != 0.0 {
		t.Errorf("median: 0.0 != %v\n", ps[0])
	}
	if ps[1] != 0.0 {
		t.Errorf("75th percentile: 0.0 != %v\n", ps[1])
	}
	if ps[2] != 0.0 {
		t.Errorf("99th percentile: 0.0 != %v\n", ps[2])
	}
	if rate1 := tm.Rate1(); rate1 != 0.0 {
		t.Errorf("tm.Rate1(): 0.0 != %v\n", rate1)
	}
	if rate5 := tm.Rate5(); rate5 != 0.0 {
		t.Errorf("tm.Rate5(): 0.0 != %v\n", rate5)
	}
	if rate15 := tm.Rate15(); rate15 != 0.0 {
		t.Errorf("tm.Rate15(): 0.0 != %v\n", rate15)
	}
	if rateMean := tm.RateMean(); rateMean != 0.0 {
		t.Errorf("tm.RateMean(): 0.0 != %v\n", rateMean)
	}
}

func TestTimerLabels(t *testing.T) {
	labels := Labels{"key1": "value1"}
	c := NewTimer(labels)
	if len(c.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(c.Labels()))
	}
	if lbls := c.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}

	// Labels passed by value.
	labels["key1"] = "value2"
	if lbls := c.Labels()["key1"]; lbls != "value1" {
		t.Error("Labels(): labels passed by reference")
	}

	// Labels in snapshot.
	ss := c.Snapshot()
	if len(ss.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(c.Labels()))
	}
	if lbls := ss.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}
}

func TestTimerWithLabels(t *testing.T) {
	c := NewTimer(Labels{"foo": "bar"})
	new := c.WithLabels(Labels{"bar": "foo"})
	if len(new.Labels()) != 2 {
		t.Fatalf("WithLabels() len: %v != 2", len(new.Labels()))
	}
	if lbls := new.Labels()["foo"]; lbls != "bar" {
		t.Errorf("WithLabels(): %v != bar", lbls)
	}
	if lbls := new.Labels()["bar"]; lbls != "foo" {
		t.Errorf("WithLabels(): %v != foo", lbls)
	}
}

func ExampleGetOrRegisterTimer() {
	m := "account.create.latency"
	t := GetOrRegisterTimer(m, nil, nil)
	t.Update(47)
	fmt.Println(t.Max()) // Output: 47
}
