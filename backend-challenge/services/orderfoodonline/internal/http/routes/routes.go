// Package routes provides HTTP routing configuration for the Order Food Online service.
package routes

import (
	"fmt"
)

// validateDependencies checks that all required dependencies are provided.
// It validates that all necessary handlers and middleware are properly initialized
// before setting up the routes to prevent runtime errors.
func validateDependencies(d Dependencies) error {
	if d.AuthMiddleware == nil {
		return fmt.Errorf("authMiddleware cannot be nil")
	}
	if d.ProductHandler == nil {
		return fmt.Errorf("productHandler cannot be nil")
	}
	if d.SwaggerHandler == nil {
		return fmt.Errorf("swaggerHandler cannot be nil")
	}
	return nil
}

// setupAPIRoutes sets up API routes using the provided dependencies.
// It configures all REST API endpoints under the /api prefix with authentication
// middleware applied to all routes. Routes include product listing, product details,
// and order placement.
func (r *Router) setupAPIRoutes(di Dependencies) error {
	if err := validateDependencies(di); err != nil {
		return err
	}
	// Group all routes under /api
	api := r.engine.Group("/api")

	// Public routes goes here(no auth required)

	// Protect all /api routes with auth middleware
	api.Use(di.AuthMiddleware.Authenticate())
	{
		api.GET("/product", di.ProductHandler.ListProducts)
		api.GET("/product/:productId", di.ProductHandler.GetProductByID)
		api.POST("/order", di.OrderHandler.PlaceOrder)
	}

	return nil
}
