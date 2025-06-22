package service

import (
	"context"
	"errors"
	"library/logger"
	"orderfoodonline/internal/repository"
	"orderfoodonline/internal/repository/models"
)

// ProductResponse is the response model for a product (for Swagger docs).
type ProductResponse struct {
	ID       string  `json:"id" example:"10"`
	Name     string  `json:"name" example:"Chicken Waffle"`
	Price    float64 `json:"price" example:"1"`
	Category string  `json:"category" example:"Waffle"`
}

type productService struct {
	repo   repository.ProductRepository
	logger logger.ILogger
}

// NewProductService creates a new ProductService.
func NewProductService(repo repository.ProductRepository, logger logger.ILogger) ProductService {
	return &productService{repo: repo, logger: logger}
}

// ListProducts retrieves all available products from the catalog.
// Returns a slice of Product models or an error if the operation fails.
func (s *productService) ListProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		s.logger.Error("%s: %v", ProductListingError, err)
		return nil, errors.New(ProductListingError)
	}
	return products, nil
}

// FindProductByID fetches a single product by its unique ID.
// Returns the Product model or an error if not found or on failure.
func (s *productService) FindProductByID(ctx context.Context, id string) (*models.Product, error) {
	return s.repo.FindProductByID(ctx, id)
}
