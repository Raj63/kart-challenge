// Package routes provides HTTP routing configuration for the Order Food Online service.
package routes

import (
	"orderfoodonline/internal/http/handlers"
	"orderfoodonline/internal/http/middlewares"
)

// Dependencies holds all the dependencies required by the router for proper operation.
// It encapsulates all HTTP handlers and middleware components needed to set up
// the complete routing configuration for the Order Food Online service.
type Dependencies struct {
	SwaggerHandler    handlers.SwaggerHandler       // Handler for serving Swagger documentation
	AuthMiddleware    middlewares.AuthMiddleware    // Middleware for authentication and authorization
	MetricsMiddleware middlewares.MetricsMiddleware // Middleware for Prometheus metrics collection
	ProductHandler    handlers.ProductHandler       // Handler for product-related endpoints
	OrderHandler      handlers.OrderHandler         // Handler for order-related endpoints
}
