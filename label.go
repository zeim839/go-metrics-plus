package metrics

// Labels is a map of key-values describing a metric.
type Labels map[string]string

func deepCopyLabels(labels Labels) Labels {
	copy := Labels{}
	for key, value := range labels {
		copy[key] = value
	}
	return copy
}
