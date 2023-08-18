package metrics

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkMeter(b *testing.B) {
	m := NewMeter(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mark(1)
	}
}

func BenchmarkMeterParallel(b *testing.B) {
	m := NewMeter(nil)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Mark(1)
		}
	})
}

// exercise race detector
func TestMeterConcurrency(t *testing.T) {
	m := newStandardMeter(nil)
	wg := &sync.WaitGroup{}
	reps := 100
	for i := 0; i < reps; i++ {
		wg.Add(1)
		go func(m Meter, wg *sync.WaitGroup) {
			m.Mark(1)
			wg.Done()
		}(m, wg)

		// Test reading from EWMA concurrently.
		wg.Add(1)
		go func(m Meter, wg *sync.WaitGroup) {
			m.Snapshot()
			wg.Done()
		}(m, wg)
	}
	wg.Wait()
}

func TestGetOrRegisterMeter(t *testing.T) {
	r := NewRegistry()
	NewRegisteredMeter("foo", r, nil).Mark(47)
	if m := GetOrRegisterMeter("foo", r, nil); m.Count() != 47 {
		t.Fatal(m)
	}
}

func TestMeterDecay(t *testing.T) {
	m := newStandardMeter(nil)
	m.Mark(1)
	rateMean := m.RateMean()
	time.Sleep(100 * time.Millisecond)
	if m.RateMean() >= rateMean {
		t.Error("m.RateMean() didn't decrease")
	}
}

func TestMeterNonzero(t *testing.T) {
	m := NewMeter(nil)
	m.Mark(3)
	if count := m.Count(); count != 3 {
		t.Errorf("m.Count(): 3 != %v\n", count)
	}
}

func TestMeterSnapshot(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	m := NewMeter(nil)
	m.Mark(r.Int63())
	if snapshot := m.Snapshot(); m.Count() != snapshot.Count() {
		t.Fatal(snapshot)
	}
}

func TestMeterZero(t *testing.T) {
	m := NewMeter(nil)
	if count := m.Count(); count != 0 {
		t.Errorf("m.Count(): 0 != %v\n", count)
	}
}

func TestMeterLabels(t *testing.T) {
	labels := Labels{"key1": "value1"}
	m := NewMeter(labels)
	if len(m.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(m.Labels()))
	}
	if lbls := m.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}

	// Labels passed by value.
	labels["key1"] = "value3"
	if lbls := m.Labels()["key1"]; lbls != "value1" {
		t.Error("Labels(): labels passed by reference")
	}

	// Labels in snapshot.
	ss := m.Snapshot()
	if len(ss.Labels()) != 1 {
		t.Fatalf("Labels(): %v != 1", len(m.Labels()))
	}
	if lbls := ss.Labels()["key1"]; lbls != "value1" {
		t.Errorf("Labels(): %v != value1", lbls)
	}
}

func TestMeterWithLabels(t *testing.T) {
	m := NewMeter(Labels{"foo": "bar"})
	new := m.WithLabels(Labels{"bar": "foo"})
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
