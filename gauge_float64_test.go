package metrics

import "testing"

func BenchmarkGuageFloat64(b *testing.B) {
	g := NewGaugeFloat64(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Update(float64(i))
	}
}

func BenchmarkGuageFloat64Parallel(b *testing.B) {
	g := NewGaugeFloat64(nil)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			g.Update(float64(1))
		}
	})
}

func TestGaugeFloat64(t *testing.T) {
	g := NewGaugeFloat64(nil)
	g.Update(float64(47.0))
	if v := g.Value(); float64(47.0) != v {
		t.Errorf("g.Value(): 47.0 != %v\n", v)
	}
}

func TestGaugeFloat64Snapshot(t *testing.T) {
	g := NewGaugeFloat64(nil)
	g.Update(float64(47.0))
	snapshot := g.Snapshot()
	g.Update(float64(0))
	if v := snapshot.Value(); float64(47.0) != v {
		t.Errorf("g.Value(): 47.0 != %v\n", v)
	}
}

func TestGetOrRegisterGaugeFloat64(t *testing.T) {
	r := NewRegistry()
	NewRegisteredGaugeFloat64("foo", r, nil).Update(float64(47.0))
	t.Logf("registry: %v", r)
	if g := GetOrRegisterGaugeFloat64("foo", r, nil); float64(47.0) != g.Value() {
		t.Fatal(g)
	}
}

func TestFunctionalGaugeFloat64(t *testing.T) {
	var counter float64
	fg := NewFunctionalGaugeFloat64(func() float64 {
		counter++
		return counter
	}, nil)
	fg.Value()
	fg.Value()
	if counter != 2 {
		t.Error("counter != 2")
	}
}

func TestGetOrRegisterFunctionalGaugeFloat64(t *testing.T) {
	r := NewRegistry()
	NewRegisteredFunctionalGaugeFloat64("foo", r, func() float64 { return 47 }, nil)
	if g := GetOrRegisterGaugeFloat64("foo", r, nil); g.Value() != 47 {
		t.Fatal(g)
	}
}

func TestGaugeFloat64Labels(t *testing.T) {
	labels := Labels{"key1": "value1"}
	g := NewGaugeFloat64(labels)
	if len(g.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(g.Labels()))
	}
	if lbls := g.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}

	// Labels passed by value.
	labels["key1"] = "value2"
	if lbls := g.Labels()["key1"]; lbls != "value1" {
		t.Error("Labels(): labels passed by reference")
	}

	// Labels in snapshot.
	ss := g.Snapshot()
	if len(ss.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(g.Labels()))
	}
	if lbls := ss.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}
}

func TestGaugeFloat64WithLabels(t *testing.T) {
	g := NewGaugeFloat64(Labels{"foo": "bar"})
	new := g.WithLabels(Labels{"bar": "foo"})
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
