package metrics

import (
	"math/rand"
	"sync"
	"testing"
)

func BenchmarkGuage(b *testing.B) {
	g := NewGauge(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Update(int64(i))
	}
}

// exercise race detector
func TestGaugeConcurrency(t *testing.T) {
	g := NewGauge(nil)
	wg := &sync.WaitGroup{}
	reps := 100
	for i := 0; i < reps; i++ {
		wg.Add(1)
		go func(g Gauge, wg *sync.WaitGroup) {
			g.Update(rand.Int63())
			wg.Done()
		}(g, wg)
	}
	wg.Wait()
}

func TestGauge(t *testing.T) {
	g := NewGauge(nil)
	g.Update(int64(47))
	if v := g.Value(); v != 47 {
		t.Errorf("g.Value(): 47 != %v\n", v)
	}
}

func TestGaugeSnapshot(t *testing.T) {
	g := NewGauge(nil)
	g.Update(int64(47))
	snapshot := g.Snapshot()
	g.Update(int64(0))
	if v := snapshot.Value(); v != 47 {
		t.Errorf("g.Value(): 47 != %v\n", v)
	}
}

func TestGetOrRegisterGauge(t *testing.T) {
	r := NewRegistry()
	NewRegisteredGauge("foo", r, nil).Update(47)
	if g := GetOrRegisterGauge("foo", r, nil); g.Value() != 47 {
		t.Fatal(g)
	}
}

func TestFunctionalGauge(t *testing.T) {
	var counter int64
	fg := NewFunctionalGauge(func() int64 {
		counter++
		return counter
	}, nil)
	fg.Value()
	fg.Value()
	if counter != 2 {
		t.Error("counter != 2")
	}
}

func TestGetOrRegisterFunctionalGauge(t *testing.T) {
	r := NewRegistry()
	NewRegisteredFunctionalGauge("foo", r, func() int64 { return 47 }, nil)
	if g := GetOrRegisterGauge("foo", r, nil); g.Value() != 47 {
		t.Fatal(g)
	}
}

func TestGaugeLabels(t *testing.T) {
	labels := Labels{"key1": "value1"}
	g := NewGauge(labels)
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

func TestGaugeWithLabels(t *testing.T) {
	g := NewGauge(Labels{"foo": "bar"})
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
