package repository

import (
	"context"
	"fmt"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// couponRepository provides MongoDB-backed access to coupon data.
type couponRepository struct {
	collection *mongo.Collection
}

// NewCouponRepository creates a new CouponRepository using the given Repository.
func NewCouponRepository(repo *Repository) (CouponRepository, error) {
	collection := repo.db.Collection("coupons")

	couponRepo := &couponRepository{collection: collection}
	if err := couponRepo.createCouponIndex(context.Background()); err != nil {
		return nil, err
	}
	return couponRepo, nil
}

func (c *couponRepository) createCouponIndex(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "coupon_code", Value: 1},
			{Key: "file_name", Value: 1},
		},
		Options: options.Index().
			SetPartialFilterExpression(bson.D{{Key: "isactive", Value: true}}).
			SetName("coupon_code_file_name_active_idx"),
	}

	// Set a timeout context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

// ValidateCouponCode validates a coupon code according to the business rules:
// 1. Must be found in at least two files (couponcode, filename distinct combo > 1)
func (c *couponRepository) ValidateCouponCode(ctx context.Context, couponCode string) (bool, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"coupon_code": couponCode,
			"isactive":    true,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": "$file_name",
		}}},
		{{Key: "$limit", Value: 2}},
	}

	cursor, err := c.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return false, fmt.Errorf("aggregation error: %w", err)
	}
	defer cursor.Close(ctx)

	// We just need to check if there are at least 2 results
	var count int
	for cursor.Next(ctx) {
		count++
		if count >= 2 {
			return true, nil
		}
	}
	if err := cursor.Err(); err != nil {
		return false, fmt.Errorf("cursor error: %w", err)
	}
	return false, nil
}
