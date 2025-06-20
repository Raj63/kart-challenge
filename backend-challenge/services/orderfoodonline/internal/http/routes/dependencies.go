package routes

import (
	"orderfoodonline/internal/http/handlers"
	"orderfoodonline/internal/http/middlewares"
)

// Dependencies holds the dependencies required by the router.
type Dependencies struct {
	SwaggerHandler handlers.SwaggerHandler
	AuthMiddleware middlewares.AuthMiddleware
	ProductHandler handlers.ProductHandler
}
