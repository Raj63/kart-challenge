package handlers

import (
	"net/http"
	"orderfoodonline/internal/repository/models"
	"orderfoodonline/internal/service"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	service service.OrderService
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(service service.OrderService) OrderHandler {
	return &orderHandler{service: service}
}

// PlaceOrder godoc
// @Summary Place an order
// @Description Place a new order
// @Tags order
// @Accept json
// @Produce json
// @Param order body models.OrderCreateRequest true "Order request"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "error":"invalid request body"
// @Failure 400 {object} map[string]string "error":"no items in order"
// @Failure 422 {object} map[string]string "error":"invalid productId or quantity"
// @Failure 422 {object} map[string]string "error":"product not found"
// @Failure 500 {object} map[string]string "error":"failed to place order"
// @Router /order [post]
func (h *orderHandler) PlaceOrder(c *gin.Context) {
	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no items in order"})
	}
	order, err := h.service.PlaceOrder(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid productId or quantity" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if len(err.Error()) > 13 && err.Error()[:13] == "product not found" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to place order"})
		return
	}
	c.JSON(http.StatusOK, order)
}
