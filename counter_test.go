package metrics

import "testing"

func BenchmarkCounter(b *testing.B) {
	c := NewCounter(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Inc(1)
	}
}

func TestCounterClear(t *testing.T) {
	c := NewCounter(nil)
	c.Inc(1)
	c.Clear()
	if count := c.Count(); count != 0 {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestCounterDec1(t *testing.T) {
	c := NewCounter(nil)
	c.Dec(1)
	if count := c.Count(); count != -1 {
		t.Errorf("c.Count(): -1 != %v\n", count)
	}
}

func TestCounterDec2(t *testing.T) {
	c := NewCounter(nil)
	c.Dec(2)
	if count := c.Count(); count != -2 {
		t.Errorf("c.Count(): -2 != %v\n", count)
	}
}

func TestCounterInc1(t *testing.T) {
	c := NewCounter(nil)
	c.Inc(1)
	if count := c.Count(); count != 1 {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestCounterInc2(t *testing.T) {
	c := NewCounter(nil)
	c.Inc(2)
	if count := c.Count(); count != 2 {
		t.Errorf("c.Count(): 2 != %v\n", count)
	}
}

func TestCounterSnapshot(t *testing.T) {
	c := NewCounter(nil)
	c.Inc(1)
	snapshot := c.Snapshot()
	c.Inc(1)
	if count := snapshot.Count(); count != 1 {
		t.Errorf("c.Count(): 1 != %v\n", count)
	}
}

func TestCounterZero(t *testing.T) {
	c := NewCounter(nil)
	if count := c.Count(); count != 0 {
		t.Errorf("c.Count(): 0 != %v\n", count)
	}
}

func TestGetOrRegisterCounter(t *testing.T) {
	r := NewRegistry()
	NewRegisteredCounter("foo", r, nil).Inc(47)
	if c := GetOrRegisterCounter("foo", r, nil); c.Count() != 47 {
		t.Fatal(c)
	}
}

func TestCounterLabels(t *testing.T) {
	labels := Labels{"key1": "value1"}
	c := NewCounter(labels)
	if len(c.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(c.Labels()))
	}
	if lbls := c.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}

	// Labels passed by value.
	labels["key1"] = "valye2"
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

func TestCounterWithLabels(t *testing.T) {
	c := NewCounter(Labels{"foo": "bar"})
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
