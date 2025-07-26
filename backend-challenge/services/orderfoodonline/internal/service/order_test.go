package service

import (
	"context"
	"errors"
	"testing"

	"orderfoodonline/internal/repository/models"

	libmocks "library/logger/mocks"
	"orderfoodonline/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewOrderService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	// When: Creating a new order service
	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	// Then: Service should be created successfully
	assert.NotNil(t, service)
	assert.IsType(t, &orderService{}, service)
}

func TestOrderService_PlaceOrder_Success(t *testing.T) {
	// Given: An order service with mock repositories and valid order request
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 2},
			{ProductID: "2", Quantity: 1},
		},
	}

	wafflePrd := models.Product{ID: "1", Name: "Chicken Waffle", Price: 12.99, Category: "Waffle"}
	burgerPrd := models.Product{ID: "2", Name: "Beef Burger", Price: 15.50, Category: "Burger"}
	expectedProducts := []models.Product{
		wafflePrd,
		burgerPrd,
	}

	expectedOrder := &models.Order{
		ID:       "order-123",
		Items:    request.Items,
		Products: expectedProducts,
	}

	ctx := context.Background()

	// Mock repository behaviors
	mockProductRepo.EXPECT().FindProductByID(ctx, "1").Return(&wafflePrd, nil)
	mockProductRepo.EXPECT().FindProductByID(ctx, "2").Return(&burgerPrd, nil)
	mockOrderRepo.EXPECT().PlaceOrder(ctx, gomock.Any()).Return(expectedOrder, nil)

	// When: Placing an order
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return order without error
	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, expectedOrder.ID, order.ID)
	assert.Len(t, order.Items, 2)
	assert.Len(t, order.Products, 2)
	assert.Equal(t, request.Items[0].ProductID, order.Items[0].ProductID)
	assert.Equal(t, request.Items[0].Quantity, order.Items[0].Quantity)

}

func TestOrderService_PlaceOrder_WithValidCoupon(t *testing.T) {
	// Given: An order service with valid coupon code
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		CouponCode: "SAVE20OFF",
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 1},
		},
	}

	expectedProduct := &models.Product{
		ID: "1", Name: "Chicken Waffle", Price: 12.99, Category: "Waffle",
	}

	expectedOrder := &models.Order{
		ID:       "order-123",
		Items:    request.Items,
		Products: []models.Product{*expectedProduct},
	}

	ctx := context.Background()

	// Mock repository behaviors
	mockCouponRepo.EXPECT().ValidateCouponCode(ctx, "SAVE20OFF").Return(true, nil)
	mockProductRepo.EXPECT().FindProductByID(ctx, "1").Return(expectedProduct, nil)
	mockOrderRepo.EXPECT().PlaceOrder(ctx, gomock.Any()).Return(expectedOrder, nil)

	// When: Placing an order with valid coupon
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return order without error
	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, expectedOrder.ID, order.ID)
}

func TestOrderService_PlaceOrder_WithInvalidCouponLength(t *testing.T) {
	// Given: An order service with invalid coupon code length
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	// Test coupon code too short
	request := &models.OrderCreateRequest{
		CouponCode: "SHORT",
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 1},
		},
	}

	ctx := context.Background()

	// When: Placing an order with invalid coupon length
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return error for invalid promo code
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, InvalidPromoCode, err.Error())

	// Test coupon code too long
	request.CouponCode = "VERYLONGCOUPONCODE"
	order, err = service.PlaceOrder(ctx, request)

	// Then: Should return error for invalid promo code
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, InvalidPromoCode, err.Error())
}

func TestOrderService_PlaceOrder_WithInvalidCouponValidation(t *testing.T) {
	// Given: An order service with invalid coupon code
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		CouponCode: "INVALID20",
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 1},
		},
	}

	ctx := context.Background()

	// Mock repository behavior for invalid coupon
	mockCouponRepo.EXPECT().ValidateCouponCode(ctx, "INVALID20").Return(false, nil)

	// When: Placing an order with invalid coupon
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return error for invalid promo code
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, InvalidPromoCode, err.Error())
}

func TestOrderService_PlaceOrder_WithCouponValidationError(t *testing.T) {
	// Given: An order service with coupon validation error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		CouponCode: "SAVE20OFF",
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 1},
		},
	}

	expectedError := errors.New("database error")
	ctx := context.Background()

	// Mock repository behavior for coupon validation error
	mockCouponRepo.EXPECT().ValidateCouponCode(ctx, "SAVE20OFF").Return(false, expectedError)

	// When: Placing an order with coupon validation error
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return the original error
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, expectedError.Error(), err.Error())

}

func TestOrderService_PlaceOrder_WithEmptyProductID(t *testing.T) {
	// Given: An order service with empty product ID
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		Items: []models.OrderItem{
			{ProductID: "", Quantity: 1},
		},
	}

	ctx := context.Background()

	// When: Placing an order with empty product ID
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return error for invalid product or quantity
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, InvalidProductOrQuantity, err.Error())
}

func TestOrderService_PlaceOrder_WithInvalidQuantity(t *testing.T) {
	// Given: An order service with invalid quantity
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 0},
		},
	}

	ctx := context.Background()

	// When: Placing an order with invalid quantity
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return error for invalid product or quantity
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, InvalidProductOrQuantity, err.Error())
}

func TestOrderService_PlaceOrder_WithNegativeQuantity(t *testing.T) {
	// Given: An order service with negative quantity
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: -1},
		},
	}

	ctx := context.Background()

	// When: Placing an order with negative quantity
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return error for invalid product or quantity
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, InvalidProductOrQuantity, err.Error())
}

func TestOrderService_PlaceOrder_WithProductNotFound(t *testing.T) {
	// Given: An order service with non-existent product
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		Items: []models.OrderItem{
			{ProductID: "999", Quantity: 1},
		},
	}

	ctx := context.Background()

	// Mock repository behavior for non-existent product
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
	mockProductRepo.EXPECT().FindProductByID(ctx, "999").Return(nil, errors.New("product not found"))

	// When: Placing an order with non-existent product
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return error for product not found
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, FindProductByIDError, err.Error())
}

func TestOrderService_PlaceOrder_WithProductRepositoryError(t *testing.T) {
	// Given: An order service with product repository error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 1},
		},
	}

	expectedError := errors.New("database connection failed")
	ctx := context.Background()

	// Mock repository behavior for product repository error
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())
	mockProductRepo.EXPECT().FindProductByID(ctx, "1").Return(nil, expectedError)

	// When: Placing an order with product repository error
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return error for product fetch failure
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, FindProductByIDError, err.Error())
}

func TestOrderService_PlaceOrder_WithOrderRepositoryError(t *testing.T) {
	// Given: An order service with order repository error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 1},
		},
	}

	expectedProduct := &models.Product{
		ID: "1", Name: "Chicken Waffle", Price: 12.99, Category: "Waffle",
	}

	expectedError := errors.New("failed to save order")
	ctx := context.Background()

	// Mock repository behaviors
	mockProductRepo.EXPECT().FindProductByID(ctx, "1").Return(expectedProduct, nil)
	mockOrderRepo.EXPECT().PlaceOrder(ctx, gomock.Any()).Return(nil, expectedError)

	// When: Placing an order with order repository error
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return the original error
	require.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, expectedError.Error(), err.Error())
}

func TestOrderService_PlaceOrder_WithWhitespaceCouponCode(t *testing.T) {
	// Given: An order service with whitespace-only coupon code
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given: Mock repositories
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockCouponRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockCouponRepo, mockLogger)

	request := &models.OrderCreateRequest{
		CouponCode: "   ", // Whitespace only
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 1},
		},
	}

	expectedProduct := &models.Product{
		ID: "1", Name: "Chicken Waffle", Price: 12.99, Category: "Waffle",
	}

	expectedOrder := &models.Order{
		ID:       "order-123",
		Items:    request.Items,
		Products: []models.Product{*expectedProduct},
	}

	ctx := context.Background()

	// Mock repository behaviors (coupon validation should not be called for whitespace)
	mockProductRepo.EXPECT().FindProductByID(ctx, "1").Return(expectedProduct, nil)
	mockOrderRepo.EXPECT().PlaceOrder(ctx, gomock.Any()).Return(expectedOrder, nil)

	// When: Placing an order with whitespace coupon code
	order, err := service.PlaceOrder(ctx, request)

	// Then: Should return order without error (whitespace is trimmed)
	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, expectedOrder.ID, order.ID)
}
