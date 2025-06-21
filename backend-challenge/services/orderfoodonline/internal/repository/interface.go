// Package repository provides data access layer interfaces for the Order Food Online service.
package repository

import (
	"context"
	"orderfoodonline/internal/repository/models"
)

// ProductRepository defines methods for accessing and managing product data from the database.
// It provides operations for listing products, finding specific products by ID,
// bulk inserting products, and managing database migrations.
type ProductRepository interface {
	// ListProducts retrieves all available products from the database.
	ListProducts(ctx context.Context) ([]models.Product, error)

	// FindProductByID retrieves a specific product by its unique identifier.
	FindProductByID(ctx context.Context, id string) (*models.Product, error)

	// BulkInsertProducts inserts multiple products into the database in a single operation.
	BulkInsertProducts(ctx context.Context, products []models.Product) error

	// GetAppliedMigrations retrieves all applied database migrations.
	GetAppliedMigrations(ctx context.Context) ([]models.Migration, error)

	// InsertMigration records a new migration in the database.
	InsertMigration(ctx context.Context, migration *models.Migration) error

	// UpdateMigration updates an existing migration record in the database.
	UpdateMigration(ctx context.Context, migration *models.Migration) error
}

// OrderRepository defines methods for placing and managing orders in the database.
// It provides operations for creating new orders and retrieving order information.
type OrderRepository interface {
	// PlaceOrder creates a new order in the database and returns the created order with its ID.
	PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error)
}

// CouponRepository defines methods for validating and managing coupon codes.
// It provides operations for checking if coupon codes are valid and active.
type CouponRepository interface {
	// ValidateCouponCode checks if a given coupon code is valid and can be applied to orders.
	// Returns true if the coupon is valid, false otherwise.
	ValidateCouponCode(ctx context.Context, couponCode string) (bool, error)
}
