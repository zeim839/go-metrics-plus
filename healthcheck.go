package metrics

// Healthcheck holds an error value describing an arbitrary up/down status.
type Healthcheck interface {
	Check()
	Error() error
	Healthy()
	Unhealthy(error)
	Labels() Labels
	WithLabels(Labels) Healthcheck
}

// NewHealthcheck constructs a new Healthcheck which will use the given
// function to update its status.
func NewHealthcheck(f func(Healthcheck), labels Labels) Healthcheck {
	if UseNilMetrics {
		return NilHealthcheck{}
	}
	return &StandardHealthcheck{nil, f, deepCopyLabels(labels)}
}

// NilHealthcheck is a no-op.
type NilHealthcheck struct{}

// Check is a no-op.
func (NilHealthcheck) Check() {}

// Error is a no-op.
func (NilHealthcheck) Error() error { return nil }

// Healthy is a no-op.
func (NilHealthcheck) Healthy() {}

// Unhealthy is a no-op.
func (NilHealthcheck) Unhealthy(error) {}

// Labels is a no-op.
func (NilHealthcheck) Labels() Labels { return Labels{} }

// WithLabels is a no-op.
func (NilHealthcheck) WithLabels(Labels) Healthcheck { return NilHealthcheck{} }

// StandardHealthcheck is the standard implementation of a Healthcheck and
// stores the pstatus and a function to call to update the status.
type StandardHealthcheck struct {
	err    error
	f      func(Healthcheck)
	labels Labels
}

// Check runs the healthcheck function to update the healthcheck's status.
func (h *StandardHealthcheck) Check() {
	h.f(h)
}

// Error returns the healthcheck's status, which will be nil if it is healthy.
func (h *StandardHealthcheck) Error() error {
	return h.err
}

// Healthy marks the healthcheck as healthy.
func (h *StandardHealthcheck) Healthy() {
	h.err = nil
}

// Unhealthy marks the healthcheck as unhealthy.  The error is stored and
// may be retrieved by the Error method.
func (h *StandardHealthcheck) Unhealthy(err error) {
	h.err = err
}

// Labels returns a deep copy of the healthcheck's labels.
func (h *StandardHealthcheck) Labels() Labels {
	return deepCopyLabels(h.labels)
}

// WithLabels returns a copy of the Healthcheck with the given labels appended
// to the current list of labels.
func (h *StandardHealthcheck) WithLabels(labels Labels) Healthcheck {
	newLabels := h.labels
	for k, v := range labels {
		newLabels[k] = v
	}
	return &StandardHealthcheck{h.err, h.f, newLabels}
}
