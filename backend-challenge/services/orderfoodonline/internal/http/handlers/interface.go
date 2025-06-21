package handlers

import "github.com/gin-gonic/gin"

// SwaggerHandler defines handlers for serving Swagger documentation.
type SwaggerHandler interface {
	GetSwaggerJSONHandler(*gin.Context)
}

// ProductHandler defines handlers for product-related endpoints.
type ProductHandler interface {
	ListProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
}

// OrderHandler defines handlers for order-related endpoints.
type OrderHandler interface {
	PlaceOrder(c *gin.Context)
}
