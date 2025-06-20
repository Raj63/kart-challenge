package main

import (
	libConfig "library/config"
	"library/logger"
	"orderfoodonline/internal/config"
	"orderfoodonline/internal/http/handlers"
	"orderfoodonline/internal/http/middlewares"
	"orderfoodonline/internal/http/routes"
)

// main is the entry point for the orderfoodonline service.
func main() {
	cfgManager := libConfig.NewConfigManager("./config.json")
	err := cfgManager.Load()
	if err != nil {
		panic(err)
	}

	appConfig, err := config.NewConfig(cfgManager)
	if err != nil {
		panic(err)
	}

	appLogger, err := logger.NewLogger(appConfig.Logger)
	if err != nil {
		panic(err)
	}

	swaggerHandler, err := handlers.NewSwaggerHandler(appConfig.Swagger, appLogger)
	if err != nil {
		appLogger.Error("failed to initialize swagger handler: %v", err)
		panic(err)
	}

	dep := routes.Dependencies{
		AuthMiddleware: middlewares.NewAuthMiddleware(appLogger),
		SwaggerHandler: swaggerHandler,
	}
	// create a new http router
	router := routes.NewRouter(appConfig, appLogger)

	// initialize the router with middlewares and routes
	err = router.Init(dep)
	if err != nil {
		appLogger.Error("failed to initialize routes: %v", err)
		panic(err)
	}

	// Start the server
	router.Run()
}
