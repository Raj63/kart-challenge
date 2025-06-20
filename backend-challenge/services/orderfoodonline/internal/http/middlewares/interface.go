package middlewares

import "github.com/gin-gonic/gin"

// AuthMiddleware defines authentication and authorization middleware methods.
type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	Authorize() gin.HandlerFunc
}
