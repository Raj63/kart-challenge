package repository

import (
	"context"
	"orderfoodonline/internal/repository/models"
)

// ProductRepository defines methods for accessing product data from the database.
type ProductRepository interface {
	ListProducts(ctx context.Context) ([]models.Product, error)
	FindProductByID(ctx context.Context, id string) (*models.Product, error)
}

// OrderRepository defines methods for placing and retrieving orders from the database.
type OrderRepository interface {
	PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error)
}
