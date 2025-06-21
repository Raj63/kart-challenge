package service

import (
	"context"
	"orderfoodonline/internal/repository/models"
)

// ProductService defines business logic for products.
type ProductService interface {
	ListProducts(ctx context.Context) ([]models.Product, error)
	FindProductByID(ctx context.Context, id string) (*models.Product, error)
}

// OrderService defines business logic for orders.
type OrderService interface {
	PlaceOrder(ctx context.Context, req *models.OrderCreateRequest) (*models.Order, error)
}
