package metrics

import "sync/atomic"

// Counter holds an int64 value that can be incremented and decremented.
type Counter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
	Snapshot() Counter
}

// GetOrRegisterCounter returns an existing Counter or constructs and registers
// a new StandardCounter.
func GetOrRegisterCounter(name string, r Registry) Counter {
	if r == nil {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewCounter).(Counter)
}

// NewRegisteredCounter constructs and registers a new StandardCounter.
func NewRegisteredCounter(name string, r Registry) Counter {
	if nil == r {
		r = DefaultRegistry
	}
	c := NewCounter()
	r.Register(name, c)
	return c
}

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter {
	if UseNilMetrics {
		return NilCounter{}
	}
	return &StandardCounter{}
}

// CounterSnapshot is a read-only copy of another Counter.
type CounterSnapshot struct {
	count int64
}

// Clear panics.
func (CounterSnapshot) Clear() {
	panic("Clear called on a CounterSnapshot")
}

// Count returns the count at the time the snapshot was taken.
func (c CounterSnapshot) Count() int64 { return c.count }

// Dec panics.
func (CounterSnapshot) Dec(int64) {
	panic("Dec called on a CounterSnapshot")
}

// Inc panics.
func (CounterSnapshot) Inc(int64) {
	panic("Inc called on a CounterSnapshot")
}

// Snapshot returns the snapshot.
func (c CounterSnapshot) Snapshot() Counter { return c }

// NilCounter is a no-op Counter.
type NilCounter struct{}

// Clear is a no-op.
func (NilCounter) Clear() {}

// Count is a no-op.
func (NilCounter) Count() int64 { return 0 }

// Dec is a no-op.
func (NilCounter) Dec(i int64) {}

// Inc is a no-op.
func (NilCounter) Inc(i int64) {}

// Snapshot is a no-op.
func (NilCounter) Snapshot() Counter { return NilCounter{} }

// StandardCounter is the standard implementation of a Counter and uses the
// sync/atomic package to manage a single int64 value.
type StandardCounter struct {
	count atomic.Int64
}

// Clear sets the counter to zero.
func (c *StandardCounter) Clear() {
	c.count.Store(0)
}

// Count returns the current count.
func (c *StandardCounter) Count() int64 {
	return c.count.Load()
}

// Dec decrements the counter by the given amount.
func (c *StandardCounter) Dec(i int64) {
	c.count.Add(-i)
}

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	c.count.Add(i)
}

// Snapshot returns a read-only copy of the counter.
func (c *StandardCounter) Snapshot() Counter {
	return CounterSnapshot{
		count: c.Count(),
	}
}
