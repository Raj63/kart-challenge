package routes

import (
	"context"
	"fmt"
	"library/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"orderfoodonline/internal/config"
	"orderfoodonline/internal/constants"
	"orderfoodonline/internal/http/middlewares"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router handles HTTP routing for the product service
type Router struct {
	engine *gin.Engine
	config *config.Config
	logger *logger.Logger
}

// NewRouter creates a new Gin engine
func NewRouter(config *config.Config, logger *logger.Logger) *Router {
	// Set up Gin router
	// Use gin.ReleaseMode for production
	// Use gin.DebugMode for local development
	engine := gin.Default()

	router := &Router{
		engine: engine,
		config: config,
		logger: logger,
	}

	return router
}

// Run starts the HTTP server and handles graceful shutdown.
func (r *Router) Run() {
	serverConfig := r.config.Server

	r.logger.Info("Server is running on port: %v", serverConfig.Port)

	// Create HTTP server
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler:           r.GetEngine(),
		ReadHeaderTimeout: serverConfig.ReadTimeout,
		WriteTimeout:      serverConfig.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	r.logger.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		r.logger.Fatal("Server forced to shutdown: %v", err)
	}

	r.logger.Info("Server exited")

}

// GetEngine returns the Gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// Init initializes the router with dependencies.
func (r *Router) Init(dep Dependencies) error {
	// setup middlewares first
	r.setupMiddleware(dep)

	// setup api routes
	if err := r.setupAPIRoutes(dep); err != nil {
		return fmt.Errorf("failed to setup api routes: %w", err)
	}

	return nil
}

// setupMiddleware sets up all required middlewares for the router.
func (r *Router) setupMiddleware(dep Dependencies) {
	// RateLimiter middleware
	r.engine.Use(middlewares.RateLimiterHandler())

	// CORS middleware
	r.engine.Use(middlewares.CorsHandler())

	// OPTIONS method handler
	r.engine.Use(middlewares.OptionsHandler())

	// Health handler
	r.engine.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Healthy")
	})

	// VERSION handler
	r.engine.GET("/api/version", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"version":     constants.Version,
			"commit_hash": constants.CommitHash,
		})
	})

	// Swagger Page handler
	if r.config.Env == "local" {
		// Serve filtered Swagger JSON
		r.engine.GET("/api/swagger.json", dep.SwaggerHandler.GetSwaggerJSONHandler)

		// Serve Swagger UI (pointing to filtered doc.json)
		r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/swagger.json")))
	} else {
		log.Println("Swagger is disabled in dev/production")
	}
}
