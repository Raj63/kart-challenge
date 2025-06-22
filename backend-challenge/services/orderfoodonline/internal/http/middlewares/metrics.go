// Package middlewares provides HTTP middleware components for the Order Food Online service.
package middlewares

import (
	"orderfoodonline/internal/metrics"
	"time"

	"github.com/gin-gonic/gin"
)

// metricsMiddleware provides HTTP request metrics collection using Prometheus.
// It records request counts, durations, and status codes for monitoring and observability.
type metricsMiddleware struct{}

// NewMetricsMiddleware creates a new instance of MetricsMiddleware.
func NewMetricsMiddleware() MetricsMiddleware {
	return &metricsMiddleware{}
}

// RecordMetrics is a Gin middleware that records HTTP request metrics.
// It captures request method, endpoint, status code, and duration for each HTTP request.
func (m *metricsMiddleware) RecordMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Extract endpoint (remove query parameters)
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		// Record metrics
		statusCode := string(rune(c.Writer.Status()))
		metrics.RecordHTTPRequest(
			c.Request.Method,
			endpoint,
			statusCode,
			duration,
		)
	}
}
