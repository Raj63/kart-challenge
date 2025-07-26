// Package middlewares provides HTTP middleware components for the Order Food Online service.
package middlewares

import "github.com/gin-gonic/gin"

// AuthMiddleware defines authentication and authorization middleware methods.
// It provides middleware functions for securing API endpoints through
// API key validation and user authorization.
type AuthMiddleware interface {
	// Authenticate creates a middleware function that validates API keys
	// and authenticates incoming requests.
	Authenticate() gin.HandlerFunc

	// Authorize creates a middleware function that checks user permissions
	// and authorizes access to protected endpoints.
	Authorize() gin.HandlerFunc
}

// MetricsMiddleware provides HTTP request metrics collection using Prometheus.
type MetricsMiddleware interface {
	// RecordMetrics is a Gin middleware that records HTTP request metrics.
	RecordMetrics() gin.HandlerFunc
}
