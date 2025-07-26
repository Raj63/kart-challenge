package handlers

import (
	"net/http"
	"orderfoodonline/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(service service.ProductService) ProductHandler {
	return &productHandler{service: service}
}

// ListProducts godoc
// @Summary List products
// @Description Get all products
// @Tags product
// @Produce json
// @Success 200 {array} service.ProductResponse
// @Failure 500 {object} map[string]string "error":"Failed to fetch products"
// @Router /api/product [get]
func (h *productHandler) ListProducts(c *gin.Context) {
	products, err := h.service.ListProducts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags product
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} service.ProductResponse
// @Failure 400 {object} map[string]string "error":"Invalid ID supplied"
// @Failure 404 {object} map[string]string "error":"Product not found"
// @Failure 500 {object} map[string]string "error":"Failed to fetch product"
// @Router /api/product/{productId} [get]
func (h *productHandler) GetProductByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("productId"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID supplied"})
		return
	}
	product, err := h.service.FindProductByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}
