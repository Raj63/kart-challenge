package routes

import (
	"fmt"
)

// validateDependencies checks that all required dependencies are provided.
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
	}

	return nil
}
