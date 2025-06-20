package middlewares

import (
	"library/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// auth provides authentication and authorization middleware.
type auth struct {
	logger *logger.Logger
}

// NewAuthMiddleware creates a new instance of auth middleware which implements AuthMiddleware.
func NewAuthMiddleware(logger *logger.Logger) AuthMiddleware {
	return &auth{logger: logger}
}

// Authentication middleware to protect routes
func (a *auth) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Authorize middleware to protect resources from bad access
func (a *auth) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// logic to verify user access could go here (if necessary)

		// If everything is fine, continue to the next handler
		c.Next()
	}
}
