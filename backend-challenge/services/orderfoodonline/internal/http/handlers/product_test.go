// Package handlers provides HTTP handlers for the Order Food Online service.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"orderfoodonline/internal/repository/models"
	servicemocks "orderfoodonline/internal/service/mocks"
)

func TestProductHandler_ListProducts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := servicemocks.NewMockProductService(ctrl)
	h := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expected := []models.Product{{ID: "1", Name: "Pizza", Price: 123, Category: "food"}}
	mockService.EXPECT().ListProducts(gomock.Any()).Return(expected, nil)

	h.ListProducts(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.Product
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expected, resp)
}

func TestProductHandler_ListProducts_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockProductService(ctrl)
	h := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockService.EXPECT().ListProducts(gomock.Any()).Return(nil, errors.New("db error"))

	h.ListProducts(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to fetch products")
}

func TestProductHandler_GetProductByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockProductService(ctrl)
	h := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "productId", Value: "1"}}

	expected := &models.Product{ID: "1", Name: "Pizza"}
	mockService.EXPECT().FindProductByID(gomock.Any(), "1").Return(expected, nil)

	h.GetProductByID(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Product
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, *expected, resp)
}

func TestProductHandler_GetProductByID_EmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockProductService(ctrl)
	h := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "productId", Value: "   "}}

	h.GetProductByID(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid ID supplied")
}

func TestProductHandler_GetProductByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockProductService(ctrl)
	h := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "productId", Value: "notfound"}}

	mockService.EXPECT().FindProductByID(gomock.Any(), "notfound").Return(nil, nil)

	h.GetProductByID(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
}

func TestProductHandler_GetProductByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := servicemocks.NewMockProductService(ctrl)
	h := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "productId", Value: "1"}}

	mockService.EXPECT().FindProductByID(gomock.Any(), "1").Return(nil, errors.New("db error"))

	h.GetProductByID(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to fetch product")
}
