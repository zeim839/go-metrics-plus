package metrics

import (
	"math"
	"sync"
	"time"
)

// EWMA calculates an exponentially-weighted per-second moving average.
type EWMA interface {
	Rate() float64
	Snapshot() EWMA
	Update(int64)
}

// NewEWMA constructs a new EWMA with the given alpha and period.
func NewEWMA(alpha float64, period time.Duration) EWMA {
	if UseNilMetrics {
		return NilEWMA{}
	}
	return &StandardEWMA{
		alpha:     alpha,
		period:    period,
		timestamp: time.Now(),
	}
}

// NewEWMA1 constructs a new EWMA for a one-minute moving average.
func NewEWMA1() EWMA {
	return NewEWMA(1-math.Exp(-5.0/60.0/1), 5*time.Second)
}

// NewEWMA5 constructs a new EWMA for a five-minute moving average.
func NewEWMA5() EWMA {
	return NewEWMA(1-math.Exp(-5.0/60.0/5), 5*time.Second)
}

// NewEWMA15 constructs a new EWMA for a fifteen-minute moving average.
func NewEWMA15() EWMA {
	return NewEWMA(1-math.Exp(-5.0/60.0/15), 5*time.Second)
}

// EWMASnapshot is a read-only copy of another EWMA.
type EWMASnapshot float64

// Rate returns the rate of events per second at the time the snapshot was
// taken.
func (a EWMASnapshot) Rate() float64 { return float64(a) }

// Snapshot returns the snapshot.
func (a EWMASnapshot) Snapshot() EWMA { return a }

// Update is a no-op.
func (EWMASnapshot) Update(int64) {}

// NilEWMA is a no-op EWMA.
type NilEWMA struct{}

// Rate is a no-op.
func (NilEWMA) Rate() float64 { return 0.0 }

// Snapshot is a no-op.
func (NilEWMA) Snapshot() EWMA { return NilEWMA{} }

// Tick is a no-op.
func (NilEWMA) Tick() {}

// Update is a no-op.
func (NilEWMA) Update(n int64) {}

// StandardEWMA is the standard implementation of an EWMA.
type StandardEWMA struct {
	alpha     float64
	period    time.Duration
	ewma      float64
	uncounted int64
	timestamp time.Time
	init      bool
	mutex     sync.Mutex
}

func (s *StandardEWMA) updateRate() {
	periods := time.Since(s.timestamp) / s.period
	rate := float64(s.uncounted) / float64(s.period)

	s.ewma = s.alpha*(rate) + (1-s.alpha)*s.ewma
	s.timestamp = s.timestamp.Add(s.period)
	s.uncounted = 0
	periods -= 1

	if !s.init {
		s.ewma = rate
		s.init = true
	}

	s.ewma = math.Pow(1-s.alpha, float64(periods)) * s.ewma
	s.timestamp = s.timestamp.Add(time.Duration(periods) * s.period)
}

// Rate returns the moving average rate of events per second.
func (s *StandardEWMA) Rate() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if time.Since(s.timestamp)/s.period < 1 {
		return s.ewma * float64(time.Second)
	}
	s.updateRate()
	return s.ewma * float64(time.Second)
}

// Snapshot returns a read-only copy of the EWMA.
func (s *StandardEWMA) Snapshot() EWMA {
	return EWMASnapshot(s.Rate())
}

// Update registers n events that occured within the last ± 0.5 sec.
func (s *StandardEWMA) Update(n int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if time.Since(s.timestamp)/s.period < 1 {
		s.uncounted += n
		return
	}
	s.updateRate()
}

// Used to elapse time in unit tests.
func (s *StandardEWMA) addToTimestamp(d time.Duration) {
	s.timestamp = s.timestamp.Add(d)
}
