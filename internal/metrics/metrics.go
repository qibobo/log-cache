package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics registers Counter and Gauge metrics.
type Initializer interface {
	// NewCounter returns a function to increment for the given metric.
	NewCounter(name string) func(delta uint64)

	// NewPerNodeCounter returns a function to increment for the given metric,
	// but requires a node index for initialization.
	NewPerNodeCounter(name string, nodeIndex int) func(delta uint64)

	// NewGauge returns a function to set the value for the given metric.
	NewGauge(name, unit string) func(value float64)
}

// NullMetrics are the default metrics.
type NullMetrics struct{}

func (m NullMetrics) NewCounter(name string) func(uint64) {
	return func(uint64) {}
}

func (m NullMetrics) NewPerNodeCounter(name string, nodeIndex int) func(uint64) {
	return func(uint64) {}
}

func (m NullMetrics) NewGauge(name, unit string) func(float64) {
	return func(float64) {}
}

// Metrics stores health metrics for the process. It has a gauge and counter
// metrics.
type Metrics struct {
	Registry *prometheus.Registry
}

// New returns a new Metrics.
func New() *Metrics {
	return &Metrics{
		Registry: prometheus.NewRegistry(),
	}
}

// NewCounter returns a func to be used increment the counter total.
func (m *Metrics) NewCounter(name string) func(delta uint64) {
	prometheusCounterMetric := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
	})
	m.Registry.MustRegister(prometheusCounterMetric)

	return func(d uint64) {
		prometheusCounterMetric.Add(float64(d))
	}
}

// NewCounter returns a func to be used increment the counter total.
func (m *Metrics) NewPerNodeCounter(name string, nodeIndex int) func(delta uint64) {
	prometheusCounterMetric := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        name,
		ConstLabels: prometheus.Labels{"nodeIndex": strconv.Itoa(nodeIndex)},
	})
	m.Registry.MustRegister(prometheusCounterMetric)

	return func(d uint64) {
		prometheusCounterMetric.Add(float64(d))
	}
}

// NewGauge returns a func to be used to set the value of a gauge metric.
func (m *Metrics) NewGauge(name, unit string) func(value float64) {
	prometheusGaugeMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		ConstLabels: map[string]string{
			"unit": unit,
		},
	})
	m.Registry.MustRegister(prometheusGaugeMetric)

	return prometheusGaugeMetric.Set
}

func (m *Metrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
