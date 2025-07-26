// Package main implements the entry point for the coupons processor service.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"coupons/internal/config"
	"coupons/internal/processor"
	"coupons/internal/repository"
	libConfig "library/config"
	"library/logger"
)

// main is the entry point for the coupons processor service.
func main() {
	// Load configuration (implement config package as needed)
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

	couponRepository := repository.NewCouponRepository(repo)

	// Start the coupon processor service
	proc := processor.NewCouponProcessor(couponRepository, appConfig.Processor, appLogger)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := proc.Run(ctx); err != nil {
		appLogger.Error("processor exited with error: %v", err)
		log.Fatalf("processor exited with error: %v", err)
	}
}
