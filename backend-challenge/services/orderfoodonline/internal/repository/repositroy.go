package repository

import (
	"context"
	"fmt"
	"orderfoodonline/internal/config"

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
	mongoURI := fmt.Sprintf("%s://%s:%d", cfg.Type, cfg.Host, cfg.Port)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb: %v", err)
	}
	return &Repository{client: mongoClient, db: mongoClient.Database(cfg.DatabaseName)}, nil
}

// Close disconnects from the MongoDB database.
func (r *Repository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
