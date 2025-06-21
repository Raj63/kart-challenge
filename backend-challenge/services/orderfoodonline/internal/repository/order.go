package repository

import (
	"context"
	"orderfoodonline/internal/repository/models"

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
	order.ID = uuid.New().String()
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}
