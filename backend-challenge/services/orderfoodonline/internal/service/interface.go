// Package service provides business logic layer interfaces for the Order Food Online service.
package service

import (
	"context"
	"orderfoodonline/internal/repository/models"
)

// ProductService defines business logic operations for product management.
// It provides high-level operations for retrieving and managing product information
// with business rules and validation.
type ProductService interface {
	// ListProducts retrieves all available products with any necessary business logic
	// such as filtering, sorting, or access control.
	ListProducts(ctx context.Context) ([]models.Product, error)

	// FindProductByID retrieves a specific product by ID with business validation
	// such as checking if the product is active or available.
	FindProductByID(ctx context.Context, id string) (*models.Product, error)
}

// OrderService defines business logic operations for order management.
// It provides high-level operations for creating and managing orders
// with business rules, validation, and cross-service coordination.
type OrderService interface {
	// PlaceOrder creates a new order with business validation including
	// product availability, coupon validation, pricing calculations, and inventory updates.
	PlaceOrder(ctx context.Context, req *models.OrderCreateRequest) (*models.Order, error)
}
