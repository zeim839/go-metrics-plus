package metrics

import "testing"

func BenchmarkHistogram(b *testing.B) {
	h := NewHistogram(NewUniformSample(100), nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Update(int64(i))
	}
}

func TestGetOrRegisterHistogram(t *testing.T) {
	r := NewRegistry()
	s := NewUniformSample(100)
	NewRegisteredHistogram("foo", r, s, nil).Update(47)
	if h := GetOrRegisterHistogram("foo", r, s, nil); h.Count() != 1 {
		t.Fatal(h)
	}
}

func TestHistogram10000(t *testing.T) {
	h := NewHistogram(NewUniformSample(100000), nil)
	for i := 1; i <= 10000; i++ {
		h.Update(int64(i))
	}
	testHistogram10000(t, h)
}

func TestHistogramEmpty(t *testing.T) {
	h := NewHistogram(NewUniformSample(100), nil)
	if count := h.Count(); count != 0 {
		t.Errorf("h.Count(): 0 != %v\n", count)
	}
	if min := h.Min(); min != 0 {
		t.Errorf("h.Min(): 0 != %v\n", min)
	}
	if max := h.Max(); max != 0 {
		t.Errorf("h.Max(): 0 != %v\n", max)
	}
	if mean := h.Mean(); mean != 0.0 {
		t.Errorf("h.Mean(): 0.0 != %v\n", mean)
	}
	if stdDev := h.StdDev(); stdDev != 0.0 {
		t.Errorf("h.StdDev(): 0.0 != %v\n", stdDev)
	}
	ps := h.Percentiles([]float64{0.5, 0.75, 0.99})
	if ps[0] != 0.0 {
		t.Errorf("median: 0.0 != %v\n", ps[0])
	}
	if ps[1] != 0.0 {
		t.Errorf("75th percentile: 0.0 != %v\n", ps[1])
	}
	if ps[2] != 0.0 {
		t.Errorf("99th percentile: 0.0 != %v\n", ps[2])
	}
}

func TestHistogramSnapshot(t *testing.T) {
	h := NewHistogram(NewUniformSample(100000), nil)
	for i := 1; i <= 10000; i++ {
		h.Update(int64(i))
	}
	snapshot := h.Snapshot()
	h.Update(0)
	testHistogram10000(t, snapshot)
}

func testHistogram10000(t *testing.T, h Histogram) {
	if count := h.Count(); count != 10000 {
		t.Errorf("h.Count(): 10000 != %v\n", count)
	}
	if min := h.Min(); min != 1 {
		t.Errorf("h.Min(): 1 != %v\n", min)
	}
	if max := h.Max(); max != 10000 {
		t.Errorf("h.Max(): 10000 != %v\n", max)
	}
	if mean := h.Mean(); mean != 5000.5 {
		t.Errorf("h.Mean(): 5000.5 != %v\n", mean)
	}
	if stdDev := h.StdDev(); stdDev != 2886.751331514372 {
		t.Errorf("h.StdDev(): 2886.751331514372 != %v\n", stdDev)
	}
	ps := h.Percentiles([]float64{0.5, 0.75, 0.99})
	if ps[0] != 5000.5 {
		t.Errorf("median: 5000.5 != %v\n", ps[0])
	}
	if ps[1] != 7500.75 {
		t.Errorf("75th percentile: 7500.75 != %v\n", ps[1])
	}
	if ps[2] != 9900.99 {
		t.Errorf("99th percentile: 9900.99 != %v\n", ps[2])
	}
}

func TestHistogramLabels(t *testing.T) {
	labels := Labels{"key1": "value1"}
	h := NewHistogram(NewUniformSample(100), labels)
	if len(h.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(h.Labels()))
	}
	if lbls := h.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}

	// Labels passed by value.
	labels["key1"] = "value2"
	if lbls := h.Labels()["key1"]; lbls != "value1" {
		t.Error("Labels(): labels passed by reference")
	}

	// Labels in snapshot.
	ss := h.Snapshot()
	if len(ss.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(h.Labels()))
	}
	if lbls := ss.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}
}

func TestHistogramWithLabels(t *testing.T) {
	h := NewHistogram(NewUniformSample(100), Labels{"foo": "bar"})
	new := h.WithLabels(Labels{"bar": "foo"})
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
