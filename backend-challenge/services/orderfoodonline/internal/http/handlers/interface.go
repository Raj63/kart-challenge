// Package handlers provides HTTP request handlers for the Order Food Online service.
package handlers

import "github.com/gin-gonic/gin"

// SwaggerHandler defines HTTP handlers for serving Swagger documentation.
// It provides endpoints for accessing API documentation and OpenAPI specifications.
type SwaggerHandler interface {
	// GetSwaggerJSONHandler serves the Swagger JSON specification for API documentation.
	GetSwaggerJSONHandler(*gin.Context)
}

// ProductHandler defines HTTP handlers for product-related endpoints.
// It provides REST API operations for managing products including listing and retrieval.
type ProductHandler interface {
	// ListProducts handles HTTP GET requests to retrieve all available products.
	// Returns a JSON response with the list of products.
	ListProducts(c *gin.Context)

	// GetProductByID handles HTTP GET requests to retrieve a specific product by ID.
	// Returns a JSON response with the product details or a 404 error if not found.
	GetProductByID(c *gin.Context)
}

// OrderHandler defines HTTP handlers for order-related endpoints.
// It provides REST API operations for creating and managing orders.
type OrderHandler interface {
	// PlaceOrder handles HTTP POST requests to create new orders.
	// Validates the order request, applies business rules, and returns the created order.
	PlaceOrder(c *gin.Context)
}
