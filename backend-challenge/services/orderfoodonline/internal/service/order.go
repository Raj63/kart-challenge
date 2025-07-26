package service

import (
	"context"
	"errors"
	"library/logger"
	"orderfoodonline/internal/metrics"
	"orderfoodonline/internal/repository"
	"orderfoodonline/internal/repository/models"
	"strings"
	"time"
)

type orderService struct {
	repo        repository.OrderRepository
	productRepo repository.ProductRepository
	couponRepo  repository.CouponRepository
	logger      logger.ILogger
}

// NewOrderService creates a new OrderService.
func NewOrderService(repo repository.OrderRepository, productRepo repository.ProductRepository,
	couponRepo repository.CouponRepository, logger logger.ILogger) OrderService {
	return &orderService{repo: repo, productRepo: productRepo, couponRepo: couponRepo, logger: logger}
}

// PlaceOrder creates a new order based on the given request.
// It validates the input, applies business logic, and persists the order.
// Returns the created order or an error if the operation fails.
func (s *orderService) PlaceOrder(ctx context.Context, req *models.OrderCreateRequest) (*models.Order, error) {
	start := time.Now()

	// Validate coupon code if provided
	if strings.TrimSpace(req.CouponCode) != "" {
		// Rule 1: Check length between 8 and 10 characters
		if len(strings.TrimSpace(req.CouponCode)) < 8 || len(strings.TrimSpace(req.CouponCode)) > 10 {
			metrics.RecordOrderProcessing("validation_error", time.Since(start).Seconds())
			metrics.RecordOrder("validation_error")
			return nil, errors.New(InvalidPromoCode)
		}
		isValid, err := s.couponRepo.ValidateCouponCode(ctx, req.CouponCode)
		if err != nil {
			metrics.RecordOrderProcessing("coupon_validation_error", time.Since(start).Seconds())
			metrics.RecordOrder("coupon_validation_error")
			return nil, err
		}
		if !isValid {
			metrics.RecordOrderProcessing("invalid_coupon", time.Since(start).Seconds())
			metrics.RecordOrder("invalid_coupon")
			return nil, errors.New(InvalidPromoCode)
		}
	}

	var products []models.Product
	for _, item := range req.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			metrics.RecordOrderProcessing("invalid_product_or_quantity", time.Since(start).Seconds())
			metrics.RecordOrder("invalid_product_or_quantity")
			return nil, errors.New(InvalidProductOrQuantity)
		}
		prod, err := s.productRepo.FindProductByID(ctx, item.ProductID)
		if err != nil {
			s.logger.Error("%s: %v", FindProductByIDError, err)
			metrics.RecordOrderProcessing("product_lookup_error", time.Since(start).Seconds())
			metrics.RecordOrder("product_lookup_error")
			return nil, errors.New(FindProductByIDError)
		}
		if prod == nil {
			metrics.RecordOrderProcessing("product_not_found", time.Since(start).Seconds())
			metrics.RecordOrder("product_not_found")
			return nil, errors.New(ProductNotFound + item.ProductID)
		}
		products = append(products, *prod)
	}

	order := &models.Order{
		Items:    req.Items,
		Products: products,
	}

	result, err := s.repo.PlaceOrder(ctx, order)
	if err != nil {
		metrics.RecordOrderProcessing("database_error", time.Since(start).Seconds())
		metrics.RecordOrder("database_error")
		return nil, err
	}

	metrics.RecordOrderProcessing("success", time.Since(start).Seconds())
	metrics.RecordOrder("success")
	return result, nil
}
