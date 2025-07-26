package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OptionsHandler Middleware handles OPTIONS requests
func OptionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
