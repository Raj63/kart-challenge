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
// @Failure 500 {object} map[string]string "error":"failed to fetch products"
// @Router /api/product [get]
func (h *productHandler) ListProducts(c *gin.Context) {
	products, err := h.service.ListProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
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
// @Failure 400 {object} map[string]string "error":"productId is required"
// @Failure 404 {object} map[string]string "error":"product not found"
// @Failure 500 {object} map[string]string "error":"failed to fetch product"
// @Router /api/product/{productId} [get]
func (h *productHandler) GetProductByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("productId"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}
	product, err := h.service.FindProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product"})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}
