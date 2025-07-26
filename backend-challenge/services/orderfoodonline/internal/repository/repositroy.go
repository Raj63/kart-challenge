package repository

import (
	"context"
	"fmt"
	"orderfoodonline/internal/config"
	"orderfoodonline/internal/metrics"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository provides database access for the orderfoodonline service.
type Repository struct {
	db     *mongo.Database
	client *mongo.Client
}

// NewRepository creates a new Repository instance with MongoDB connection.
func NewRepository(ctx context.Context, cfg *config.DbConfig) (*Repository, error) {
	start := time.Now()

	mongoURI := fmt.Sprintf("%s://%s:%d", cfg.Type, cfg.Host, cfg.Port)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		metrics.RecordDatabaseQuery("connect", "database", "error", time.Since(start).Seconds())
		return nil, fmt.Errorf("error connecting to mongodb: %v", err)
	}

	// Test the connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		metrics.RecordDatabaseQuery("ping", "database", "error", time.Since(start).Seconds())
		return nil, fmt.Errorf("error pinging mongodb: %v", err)
	}

	metrics.RecordDatabaseQuery("connect", "database", "success", time.Since(start).Seconds())
	metrics.SetActiveConnections(1) // Set initial connection count

	return &Repository{client: mongoClient, db: mongoClient.Database(cfg.DatabaseName)}, nil
}

// Close disconnects from the MongoDB database.
func (r *Repository) Close(ctx context.Context) error {
	start := time.Now()

	err := r.client.Disconnect(ctx)
	if err != nil {
		metrics.RecordDatabaseQuery("disconnect", "database", "error", time.Since(start).Seconds())
		return err
	}

	metrics.RecordDatabaseQuery("disconnect", "database", "success", time.Since(start).Seconds())
	metrics.SetActiveConnections(0) // Set connection count to 0

	return nil
}
