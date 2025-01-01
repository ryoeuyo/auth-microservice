package metric

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	AuthRequests        *prometheus.CounterVec
	AuthRequestDuration *prometheus.HistogramVec
	AuthFailedAttempts  *prometheus.CounterVec
}

func NewMetric() *Metric {
	authRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_requests_total",
			Help: "Total number of authentication requests",
		},
		[]string{"method"}, // Для фильтрации по методу (например, login, register)
	)

	authRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_request_duration_seconds",
			Help:    "Duration of authentication requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	authFailedAttempts := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_failed_attempts_total",
			Help: "Total number of failed authentication attempts",
		},
		[]string{"method"},
	)

	prometheus.MustRegister(authRequests, authRequestDuration, authFailedAttempts)

	return &Metric{
		AuthRequests:        authRequests,
		AuthRequestDuration: authRequestDuration,
		AuthFailedAttempts:  authFailedAttempts,
	}
}
