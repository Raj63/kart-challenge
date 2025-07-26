// Package handlers provides HTTP handlers for the Order Food Online service.
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"errors"
	"orderfoodonline/internal/repository/models"
	"orderfoodonline/internal/service"
	servicemocks "orderfoodonline/internal/service/mocks"
)

func TestOrderHandler_PlaceOrder_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockOrderService(ctrl)
	h := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/order", bytes.NewBuffer([]byte("not-json")))

	h.PlaceOrder(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestOrderHandler_PlaceOrder_EmptyItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockOrderService(ctrl)
	h := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(models.OrderCreateRequest{Items: []models.OrderItem{}})
	c.Request, _ = http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.PlaceOrder(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestOrderHandler_PlaceOrder_ValidationException(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockOrderService(ctrl)
	h := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(models.OrderCreateRequest{Items: []models.OrderItem{{ProductID: "p1", Quantity: 0}}})
	c.Request, _ = http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New(service.InvalidProductOrQuantity))

	h.PlaceOrder(c)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), "Validation exception")
}

func TestOrderHandler_PlaceOrder_ProductNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockOrderService(ctrl)
	h := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(models.OrderCreateRequest{Items: []models.OrderItem{{ProductID: "notfound", Quantity: 1}}})
	c.Request, _ = http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New(service.ProductNotFound))

	h.PlaceOrder(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
}

func TestOrderHandler_PlaceOrder_InvalidPromoCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockOrderService(ctrl)
	h := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(models.OrderCreateRequest{Items: []models.OrderItem{{ProductID: "p1", Quantity: 1}}, CouponCode: "badcode"})
	c.Request, _ = http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New(service.InvalidPromoCode))

	h.PlaceOrder(c)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), "Validation exception")
}

func TestOrderHandler_PlaceOrder_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockOrderService(ctrl)
	h := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(models.OrderCreateRequest{Items: []models.OrderItem{{ProductID: "p1", Quantity: 1}}})
	c.Request, _ = http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

	h.PlaceOrder(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to place an order")
}

func TestOrderHandler_PlaceOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockOrderService(ctrl)
	h := NewOrderHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	orderReq := models.OrderCreateRequest{Items: []models.OrderItem{{ProductID: "p1", Quantity: 1}}}
	body, _ := json.Marshal(orderReq)
	c.Request, _ = http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	expectedOrder := &models.Order{ID: "order1"}
	mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(expectedOrder, nil)

	h.PlaceOrder(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Order
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedOrder.ID, resp.ID)
}
