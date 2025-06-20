package handlers

import (
	"fmt"
	"library/logger"
	"net/http"
	"orderfoodonline/internal/config"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swaggerHandler implements the SwaggerHandler interface for serving Swagger docs.
type swaggerHandler struct {
	jsonData []byte
}

// NewSwaggerHandler creates a new SwaggerHandler for serving Swagger documentation.
func NewSwaggerHandler(config *config.SwaggerConfig, logger *logger.Logger) (SwaggerHandler, error) {
	// Read full Swagger JSON
	jsonData, err := os.ReadFile(config.FilePath)
	if err != nil {
		logger.Error("error reading swagger.json: %v", err)
		return nil, fmt.Errorf("error reading swagger.json: %w", err)
	}
	return &swaggerHandler{
		jsonData: jsonData,
	}, nil
}

// GetSwaggerJSONHandler serves the Swagger JSON documentation.
func (s *swaggerHandler) GetSwaggerJSONHandler(c *gin.Context) {
	if len(s.jsonData) > 0 {
		c.Data(200, "application/json", s.jsonData)
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to serve Swagger Json"})
}

// GetSwaggerUIHandler serves the Swagger UI documentation.
func (s *swaggerHandler) GetSwaggerUIHandler(c *gin.Context) {
	if len(s.jsonData) > 0 {
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/swagger.json"))
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to serve Swagger Json"})
}
