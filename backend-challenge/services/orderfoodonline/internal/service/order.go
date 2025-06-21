package service

import (
	"context"
	"errors"
	"orderfoodonline/internal/repository"
	"orderfoodonline/internal/repository/models"
)

type orderService struct {
	repo        repository.OrderRepository
	productRepo repository.ProductRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(repo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{repo: repo, productRepo: productRepo}
}

// PlaceOrder creates a new order based on the given request.
// It validates the input, applies business logic, and persists the order.
// Returns the created order or an error if the operation fails.
func (s *orderService) PlaceOrder(ctx context.Context, req *models.OrderCreateRequest) (*models.Order, error) {
	var products []models.Product
	for _, item := range req.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return nil, errors.New("invalid productId or quantity")
		}
		prod, err := s.productRepo.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if prod == nil {
			return nil, errors.New("product not found: " + item.ProductID)
		}
		products = append(products, *prod)
	}
	order := &models.Order{
		Items:    req.Items,
		Products: products,
	}
	return s.repo.PlaceOrder(ctx, order)
}
