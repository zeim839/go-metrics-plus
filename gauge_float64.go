package metrics

import (
	"math"
	"sync/atomic"
)

// GaugeFloat64 holds a float64 value that can be set arbitrarily.
type GaugeFloat64 interface {
	Snapshot() GaugeFloat64
	Update(float64)
	Value() float64
	Labels() Labels
	WithLabels(Labels) GaugeFloat64
}

// GetOrRegisterGaugeFloat64 returns an existing GaugeFloat64 or constructs and registers a
// new StandardGaugeFloat64.
func GetOrRegisterGaugeFloat64(name string, r Registry, labels Labels) GaugeFloat64 {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGaugeFloat64(labels)).(GaugeFloat64)
}

// NewGaugeFloat64 constructs a new StandardGaugeFloat64.
func NewGaugeFloat64(labels Labels) GaugeFloat64 {
	if UseNilMetrics {
		return NilGaugeFloat64{}
	}
	return &StandardGaugeFloat64{labels: deepCopyLabels(labels)}
}

// NewRegisteredGaugeFloat64 constructs and registers a new StandardGaugeFloat64.
func NewRegisteredGaugeFloat64(name string, r Registry, labels Labels) GaugeFloat64 {
	c := NewGaugeFloat64(labels)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewFunctionalGaugeFloat64 constructs a new FunctionalGauge.
func NewFunctionalGaugeFloat64(f func() float64, labels Labels) GaugeFloat64 {
	if UseNilMetrics {
		return NilGaugeFloat64{}
	}
	return &FunctionalGaugeFloat64{value: f, labels: deepCopyLabels(labels)}
}

// NewRegisteredFunctionalGaugeFloat64 constructs and registers a new StandardGauge.
func NewRegisteredFunctionalGaugeFloat64(name string, r Registry,
	f func() float64, labels Labels) GaugeFloat64 {
	c := NewFunctionalGaugeFloat64(f, labels)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// GaugeFloat64Snapshot is a read-only copy of another GaugeFloat64.
type GaugeFloat64Snapshot struct {
	value  float64
	labels Labels
}

// Snapshot returns the snapshot.
func (g GaugeFloat64Snapshot) Snapshot() GaugeFloat64 { return g }

// Update panics.
func (GaugeFloat64Snapshot) Update(float64) {
	panic("Update called on a GaugeFloat64Snapshot")
}

// Value returns the value at the time the snapshot was taken.
func (g GaugeFloat64Snapshot) Value() float64 { return g.value }

// Labels returns a deep copy of the Snapshot's labels.
func (g GaugeFloat64Snapshot) Labels() Labels { return deepCopyLabels(g.labels) }

// WithLabels returns a copy of the snapshot with the specified labels appended
// to the current list of labels.
func (g GaugeFloat64Snapshot) WithLabels(labels Labels) GaugeFloat64 {
	newLabels := g.labels
	for k, v := range labels {
		newLabels[k] = v
	}
	return GaugeFloat64Snapshot{
		value:  g.Value(),
		labels: newLabels,
	}
}

// NilGaugeFloat64 is a no-op Gauge.
type NilGaugeFloat64 struct{}

// Snapshot is a no-op.
func (NilGaugeFloat64) Snapshot() GaugeFloat64 { return NilGaugeFloat64{} }

// Update is a no-op.
func (NilGaugeFloat64) Update(v float64) {}

// Value is a no-op.
func (NilGaugeFloat64) Value() float64 { return 0.0 }

// Labels is a no-op.
func (NilGaugeFloat64) Labels() Labels { return Labels{} }

// WithLabels is a no-op.
func (NilGaugeFloat64) WithLabels(Labels) GaugeFloat64 {
	return NilGaugeFloat64{}
}

// StandardGaugeFloat64 is the standard implementation of a GaugeFloat64 and uses
// atomic uint64 (holds float bytes) to manage a single float64 value.
type StandardGaugeFloat64 struct {
	value  atomic.Uint64
	labels Labels
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGaugeFloat64) Snapshot() GaugeFloat64 {
	return GaugeFloat64Snapshot{
		value:  g.Value(),
		labels: g.Labels(),
	}
}

// Update updates the gauge's value.
func (g *StandardGaugeFloat64) Update(v float64) {
	g.value.Store(math.Float64bits(v))
}

// Value returns the gauge's current value.
func (g *StandardGaugeFloat64) Value() float64 {
	return math.Float64frombits(g.value.Load())
}

// Labels returns a deep copy of the gauge's labels.
func (g *StandardGaugeFloat64) Labels() Labels {
	return deepCopyLabels(g.labels)
}

// WithLabels returns a snapshot of the Gauge with the given labels appended to
// the current list of labels.
func (g *StandardGaugeFloat64) WithLabels(labels Labels) GaugeFloat64 {
	return g.Snapshot().WithLabels(labels)
}

// FunctionalGaugeFloat64 returns value from given function
type FunctionalGaugeFloat64 struct {
	value  func() float64
	labels Labels
}

// Value returns the gauge's current value.
func (g FunctionalGaugeFloat64) Value() float64 {
	return g.value()
}

// Snapshot returns the snapshot.
func (g FunctionalGaugeFloat64) Snapshot() GaugeFloat64 {
	return GaugeFloat64Snapshot{
		value:  g.Value(),
		labels: g.Labels(),
	}
}

// Update panics.
func (FunctionalGaugeFloat64) Update(float64) {
	panic("Update called on a FunctionalGaugeFloat64")
}

// Labels returns a deep copy of the gauge's labels.
func (g FunctionalGaugeFloat64) Labels() Labels {
	return deepCopyLabels(g.labels)
}

// WithLabels returns a snapshot of the Gauge with the given labels appended to
// the current list of labels.
func (g FunctionalGaugeFloat64) WithLabels(labels Labels) GaugeFloat64 {
	return g.Snapshot().WithLabels(labels)
}
