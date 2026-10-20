package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	RequestLatency *prometheus.HistogramVec
	HttpCollector  *prometheus.CounterVec
}

type MetricLabelNames string

var (
	metricsCacheMu sync.Mutex
	metricsCache   = make(map[string]*Metrics)
)

func CreateMetrics(namespace, subsystem string, labels prometheus.Labels) *Metrics {
	cacheKey := namespace + "/" + subsystem

	metricsCacheMu.Lock()
	defer metricsCacheMu.Unlock()

	if cached, ok := metricsCache[cacheKey]; ok {
		return cached
	}

	var (
		requestLatency = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "http_latency",
				Help:        "The API methods latency info",
				ConstLabels: labels,
			},
			[]string{"service", "method", "success"})
		httpCollector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   namespace,
				Subsystem:   subsystem,
				Name:        "http_collector",
				Help:        "The number of http requests",
				ConstLabels: labels,
			},
			[]string{"service", "method", "success"},
		)
	)

	prometheus.MustRegister(httpCollector)
	prometheus.MustRegister(requestLatency)

	created := &Metrics{
		HttpCollector:  httpCollector,
		RequestLatency: requestLatency,
	}
	metricsCache[cacheKey] = created

	return created
}

func NewMetrics(serviceName string) *Metrics {
	return CreateMetrics("arm", serviceName, nil)
}
