package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// These are the counters we will use for our request stats
var (
	responses = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "http",
		Name:      "responses",
		Help:      "The count of requests received for each endpoint and the response code returned",
	}, []string{"code", "path"})
	latencies = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "latency",
		Help:      "The latency of the requests per endpoint",
		Buckets:   prometheus.ExponentialBuckets(.1, 3, 10),
	}, []string{"endpoint"})
	requests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "http",
		Name:      "active_requests",
		Help:      "The number of active requests",
	}, []string{"path"})
)

// MetricsRecorder is the middleware that records metrics for the request
func MetricsRecorder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Track the number of concurrent requests. Note that httpsnoop here is calling through to the next middleware
		requests.WithLabelValues(r.URL.Path).Inc()
		m := httpsnoop.CaptureMetrics(next, w, r)
		requests.WithLabelValues(r.URL.Path).Dec()

		// Count the response codes per endpoint
		responses.WithLabelValues(strconv.Itoa(m.Code), r.URL.Path).Inc()

		// Track the latency of the requests in milliseconds
		latency := float64(m.Duration.Nanoseconds()) / float64(time.Millisecond)
		latencies.WithLabelValues(r.URL.Path).Observe(latency)
	})
}
