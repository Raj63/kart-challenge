// Package metrics provides Prometheus metrics collection for the Order Food Online service.
// It defines and manages metrics for HTTP requests, database operations, and application performance.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestTotal tracks the total number of HTTP requests by method and status code
	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// HTTPRequestDuration tracks the duration of HTTP requests by method and endpoint
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// DatabaseQueryDuration tracks the duration of database queries by operation
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "collection"},
	)

	// DatabaseQueryTotal tracks the total number of database queries by operation and status
	DatabaseQueryTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "collection", "status"},
	)

	// ActiveConnections tracks the number of active database connections
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_active_connections",
			Help: "Number of active database connections",
		},
	)

	// OrderProcessingDuration tracks the duration of order processing
	OrderProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_processing_duration_seconds",
			Help:    "Duration of order processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	// OrdersTotal tracks the total number of orders by status
	OrdersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_total",
			Help: "Total number of orders",
		},
		[]string{"status"},
	)
)

// RecordHTTPRequest records an HTTP request with method, endpoint, status code, and duration
func RecordHTTPRequest(method, endpoint, statusCode string, duration float64) {
	HTTPRequestTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// RecordDatabaseQuery records a database query with operation, collection, status, and duration
func RecordDatabaseQuery(operation, collection, status string, duration float64) {
	DatabaseQueryTotal.WithLabelValues(operation, collection, status).Inc()
	DatabaseQueryDuration.WithLabelValues(operation, collection).Observe(duration)
}

// SetActiveConnections sets the number of active database connections
func SetActiveConnections(count float64) {
	ActiveConnections.Set(count)
}

// RecordOrderProcessing records order processing with status and duration
func RecordOrderProcessing(status string, duration float64) {
	OrderProcessingDuration.WithLabelValues(status).Observe(duration)
}

// RecordOrder records an order with status
func RecordOrder(status string) {
	OrdersTotal.WithLabelValues(status).Inc()
}
