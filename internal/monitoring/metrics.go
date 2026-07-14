package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"runtime"
)

// метрики приложения
type Metrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
	OperationsTotal  *prometheus.CounterVec
	ErrorsTotal      *prometheus.CounterVec
	GoroutinesCount  prometheus.GaugeFunc
}

// конструктор метрик
func NewMetrics() *Metrics {
	m := &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),

		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),

		RequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),

		OperationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_operations_total",
				Help: "Total number of operations processed by type",
			},
			[]string{"operation_type"},
		),

		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_errors_total",
				Help: "Total number of errors occurred by type",
			},
			[]string{"error_type", "handler"},
		),
	}

	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "go_goroutines_count",
			Help: "Number of goroutines",
		},
		func() float64 {
			return float64(runtime.NumGoroutine())
		},
	))

	return m
}
