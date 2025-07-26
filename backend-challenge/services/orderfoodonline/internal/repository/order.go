package repository

import (
	"context"
	"orderfoodonline/internal/metrics"
	"orderfoodonline/internal/repository/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// orderRepository provides MongoDB-backed access to order data.
type orderRepository struct {
	collection *mongo.Collection
}

// NewOrderRepository creates a new OrderRepository using the given Repository.
func NewOrderRepository(repo *Repository) OrderRepository {
	collection := repo.db.Collection("orders")
	return &orderRepository{collection: collection}
}

// PlaceOrder inserts a new order into the database and returns the created order.
func (r *orderRepository) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	start := time.Now()

	order.ID = uuid.New().String()
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		metrics.RecordDatabaseQuery("insert_one", "orders", "error", time.Since(start).Seconds())
		return nil, err
	}

	metrics.RecordDatabaseQuery("insert_one", "orders", "success", time.Since(start).Seconds())
	return order, nil
}
