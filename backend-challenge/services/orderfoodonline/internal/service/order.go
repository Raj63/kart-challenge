package service

import (
	"context"
	"errors"
	"library/logger"
	"orderfoodonline/internal/repository"
	"orderfoodonline/internal/repository/models"
	"strings"
)

type orderService struct {
	repo        repository.OrderRepository
	productRepo repository.ProductRepository
	couponRepo  repository.CouponRepository
	logger      *logger.Logger
}

// NewOrderService creates a new OrderService.
func NewOrderService(repo repository.OrderRepository, productRepo repository.ProductRepository,
	couponRepo repository.CouponRepository, logger *logger.Logger) OrderService {
	return &orderService{repo: repo, productRepo: productRepo, couponRepo: couponRepo, logger: logger}
}

// PlaceOrder creates a new order based on the given request.
// It validates the input, applies business logic, and persists the order.
// Returns the created order or an error if the operation fails.
func (s *orderService) PlaceOrder(ctx context.Context, req *models.OrderCreateRequest) (*models.Order, error) {
	// Validate coupon code if provided
	if strings.TrimSpace(req.CouponCode) != "" {
		// Rule 1: Check length between 8 and 10 characters
		if len(strings.TrimSpace(req.CouponCode)) < 8 || len(strings.TrimSpace(req.CouponCode)) > 10 {
			return nil, errors.New(InvalidPromoCode)
		}
		isValid, err := s.couponRepo.ValidateCouponCode(ctx, req.CouponCode)
		if err != nil {
			return nil, err
		}
		if !isValid {
			return nil, errors.New(InvalidPromoCode)
		}
	}

	var products []models.Product
	for _, item := range req.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return nil, errors.New(InvalidProductOrQuantity)
		}
		prod, err := s.productRepo.FindProductByID(ctx, item.ProductID)
		if err != nil {
			s.logger.Error("%s: %v", FindProductByIDError, err)
			return nil, errors.New(FindProductByIDError)
		}
		if prod == nil {
			return nil, errors.New(ProductNotFound + item.ProductID)
		}
		products = append(products, *prod)
	}
	order := &models.Order{
		Items:    req.Items,
		Products: products,
	}
	return s.repo.PlaceOrder(ctx, order)
}
