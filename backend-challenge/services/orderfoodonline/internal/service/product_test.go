package service

import (
	"context"
	"errors"
	libmocks "library/logger/mocks"
	"testing"

	"orderfoodonline/internal/repository/mocks"
	"orderfoodonline/internal/repository/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewProductService(t *testing.T) {
	// Given: A mock repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	// When: Creating a new product service
	service := NewProductService(mockProductRepo, mockLogger)

	// Then: Service should be created successfully
	assert.NotNil(t, service)
	assert.IsType(t, &productService{}, service)
}

func TestProductService_ListProducts_Success(t *testing.T) {
	// Given: A product service with mock repository and expected products
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	// When: Creating a new product service
	service := NewProductService(mockProductRepo, mockLogger)

	expectedProducts := []models.Product{
		{ID: "1", Name: "Chicken Waffle", Price: 12.99, Category: "Waffle"},
		{ID: "2", Name: "Beef Burger", Price: 15.50, Category: "Burger"},
	}

	ctx := context.Background()

	// Mock repository behavior
	mockProductRepo.EXPECT().ListProducts(ctx).Return(expectedProducts, nil)

	// When: Listing products
	products, err := service.ListProducts(ctx)

	// Then: Should return products without error
	require.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, expectedProducts[0].ID, products[0].ID)
	assert.Equal(t, expectedProducts[0].Name, products[0].Name)
	assert.Equal(t, expectedProducts[0].Price, products[0].Price)
	assert.Equal(t, expectedProducts[0].Category, products[0].Category)
}

func TestProductService_ListProducts_RepositoryError(t *testing.T) {
	// Given: A product service with mock repository that returns an error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	// When: Creating a new product service
	service := NewProductService(mockProductRepo, mockLogger)

	expectedError := errors.New("database connection failed")
	ctx := context.Background()

	// Mock repository behavior to return error
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
	mockProductRepo.EXPECT().ListProducts(ctx).Return([]models.Product{}, expectedError)

	// When: Listing products
	products, err := service.ListProducts(ctx)

	// Then: Should return error and empty products
	require.Error(t, err)
	assert.Empty(t, products)
	assert.Equal(t, ProductListingError, err.Error())
}

func TestProductService_ListProducts_EmptyResult(t *testing.T) {
	// Given: A product service with mock repository that returns empty list
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	// When: Creating a new product service
	service := NewProductService(mockProductRepo, mockLogger)

	ctx := context.Background()

	// Mock repository behavior to return empty list
	mockProductRepo.EXPECT().ListProducts(ctx).Return([]models.Product{}, nil)

	// When: Listing products
	products, err := service.ListProducts(ctx)

	// Then: Should return empty list without error
	require.NoError(t, err)
	assert.Empty(t, products)
}

func TestProductService_FindProductByID_Success(t *testing.T) {
	// Given: A product service with mock repository and expected product
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	// When: Creating a new product service
	service := NewProductService(mockProductRepo, mockLogger)

	expectedProduct := &models.Product{
		ID:       "1",
		Name:     "Chicken Waffle",
		Price:    12.99,
		Category: "Waffle",
	}

	ctx := context.Background()
	productID := "1"

	// Mock repository behavior
	mockProductRepo.EXPECT().FindProductByID(ctx, productID).Return(expectedProduct, nil)

	// When: Finding product by ID
	product, err := service.FindProductByID(ctx, productID)

	// Then: Should return product without error
	require.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Name, product.Name)
	assert.Equal(t, expectedProduct.Price, product.Price)
	assert.Equal(t, expectedProduct.Category, product.Category)
}

func TestProductService_FindProductByID_NotFound(t *testing.T) {
	// Given: A product service with mock repository that returns not found
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	// When: Creating a new product service
	service := NewProductService(mockProductRepo, mockLogger)

	ctx := context.Background()
	productID := "999"

	// Mock repository behavior to return not found
	mockProductRepo.EXPECT().FindProductByID(ctx, productID).Return(nil, errors.New("product not found"))

	// When: Finding product by ID
	product, err := service.FindProductByID(ctx, productID)

	// Then: Should return error and nil product
	require.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "product not found", err.Error())
}

func TestProductService_FindProductByID_RepositoryError(t *testing.T) {
	// Given: A product service with mock repository that returns database error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	// When: Creating a new product service
	service := NewProductService(mockProductRepo, mockLogger)

	expectedError := errors.New("database connection failed")
	ctx := context.Background()
	productID := "1"

	// Mock repository behavior to return error
	mockProductRepo.EXPECT().FindProductByID(ctx, productID).Return(nil, expectedError)

	// When: Finding product by ID
	product, err := service.FindProductByID(ctx, productID)

	// Then: Should return error and nil product
	require.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, expectedError.Error(), err.Error())
}

func TestProductService_FindProductByID_EmptyID(t *testing.T) {
	// Given: A product service with mock repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	// When: Creating a new product service
	service := NewProductService(mockProductRepo, mockLogger)

	ctx := context.Background()
	emptyProductID := ""

	// Mock repository behavior for empty ID
	mockProductRepo.EXPECT().FindProductByID(ctx, emptyProductID).Return(nil, errors.New("invalid product ID"))

	// When: Finding product with empty ID
	product, err := service.FindProductByID(ctx, emptyProductID)

	// Then: Should return error and nil product
	require.Error(t, err)
	assert.Nil(t, product)

}
