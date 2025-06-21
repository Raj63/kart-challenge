package main

import (
	"context"
	libConfig "library/config"
	"library/logger"
	"log"
	"orderfoodonline/internal/config"
	"orderfoodonline/internal/http/handlers"
	"orderfoodonline/internal/http/middlewares"
	"orderfoodonline/internal/http/routes"
	"orderfoodonline/internal/repository"
	"orderfoodonline/internal/service"
)

// main is the entry point for the orderfoodonline service.
// @title Order Food Online
// @version 2.0
// @description Documentation's Order Food Online
// @termsOfService http://swagger.io/terms/

// @contact.name Rajesh Kumar Biswas
// @contact.url https://github.com/Raj63
// @contact.email biswas.rajesh63@gmail.com

// @host localhost:8080

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api_key
// @description Provide your API key as the value of the 'api_key' header to authenticate

// router.Init is a function that contains all routes of the application
//
// @description Contains all routes of the application
// @summary Contains all routes of the application
// @tags [Application]
// @accept json
// @produce json
func main() {
	cfgManager := libConfig.NewConfigManager("./config.json")
	err := cfgManager.Load()
	if err != nil {
		log.Fatalf("failed to initialize config-manager: %v", err)
	}

	appConfig, err := config.NewConfig(cfgManager)
	if err != nil {
		log.Fatalf("failed to create app config: %v", err)
	}

	appLogger, err := logger.NewLogger(appConfig.Logger)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	ctx := context.Background()

	repo, err := repository.NewRepository(ctx, appConfig.Database)
	if err != nil {
		appLogger.Error("failed to initialize repository: %v", err)
		log.Fatalf("failed to initialize repository: %v", err)
	}
	defer repo.Close(context.Background())

	productRepository, err := repository.NewProductRepository(repo)
	if err != nil {
		appLogger.Error("failed to initialize product repository: %v", err)
		log.Fatalf("failed to initialize product repository: %v", err)
	}

	orderRepository := repository.NewOrderRepository(repo)
	orderService := service.NewOrderService(orderRepository, productRepository)
	orderHandler := handlers.NewOrderHandler(orderService)

	productService := service.NewProductService(productRepository, appLogger)

	swaggerHandler, err := handlers.NewSwaggerHandler(appConfig.Swagger, appLogger)
	if err != nil {
		appLogger.Error("failed to initialize swagger handler: %v", err)
		log.Fatalf("failed to initialize swagger handler: %v", err)

	}
	productHandler := handlers.NewProductHandler(productService)

	dep := routes.Dependencies{
		AuthMiddleware: middlewares.NewAuthMiddleware(appLogger),
		SwaggerHandler: swaggerHandler,
		ProductHandler: productHandler,
		OrderHandler:   orderHandler,
	}
	// create a new http router
	router := routes.NewRouter(appConfig, appLogger)

	// initialize the router with middlewares and routes
	err = router.Init(dep)
	if err != nil {
		appLogger.Error("failed to initialize routes: %v", err)
		log.Fatalf("failed to initialize routes: %v", err)
	}

	// Start the server
	router.Run()
}
