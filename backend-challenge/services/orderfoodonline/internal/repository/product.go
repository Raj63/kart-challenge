package repository

import (
	"context"

	"orderfoodonline/internal/repository/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// productRepository provides MongoDB-backed access to product data.
type productRepository struct {
	collection *mongo.Collection
}

// NewProductRepository creates a new ProductRepository using the given Repository.
func NewProductRepository(repo *Repository) (ProductRepository, error) {
	collection := repo.db.Collection("products")
	return &productRepository{collection: collection}, nil
}

// ListProducts returns all products from the database.
func (r *productRepository) ListProducts(ctx context.Context) ([]models.Product, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var products []models.Product
	for cur.Next(ctx) {
		var p models.Product
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// FindProductByID returns a product by its ID, or nil if not found.
func (r *productRepository) FindProductByID(ctx context.Context, id string) (*models.Product, error) {
	var p models.Product
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
