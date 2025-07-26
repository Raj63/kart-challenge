package repository

import (
	"context"
	"coupons/internal/repository/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// couponRepository provides MongoDB-backed access to coupon data.
type couponRepository struct {
	couponCollection         Collection
	processedFilesCollection Collection
}

// NewCouponRepository creates a new CouponRepository using the given MongoDB database.
func NewCouponRepository(repo *Repository) CouponRepository {
	return &couponRepository{
		couponCollection:         repo.db.Collection("coupons"),
		processedFilesCollection: repo.db.Collection("processed-coupon-files"),
	}
}

// AddCoupons inserts new coupon codes into the database as active.
// If a coupon with the same code and filename already exists, it updates the datetime instead of inserting a duplicate.
func (c *couponRepository) AddCoupons(ctx context.Context, fileName string, codes []string) error {
	if len(codes) == 0 {
		return nil
	}

	var models []mongo.WriteModel
	now := time.Now().Unix()

	for _, code := range codes {
		filter := bson.M{"coupon_code": code, "file_name": fileName}
		update := bson.M{
			"$set": bson.M{
				"datetime": now,
				"isactive": true,
			},
			"$setOnInsert": bson.M{
				"id":          uuid.New().String(),
				"coupon_code": code,
				"file_name":   fileName,
			},
		}

		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true))
	}

	if len(models) == 0 {
		return nil
	}

	_, err := c.couponCollection.BulkWrite(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to upsert coupons: %w", err)
	}

	return nil
}

// DeactivateCoupons marks coupon codes as inactive in the database.
func (c *couponRepository) DeactivateCoupons(ctx context.Context, fileName string, codes []string) error {
	filter := bson.M{"file_name": fileName, "coupon_code": bson.M{"$in": codes}}
	update := bson.M{"$set": bson.M{"isactive": false}}
	_, err := c.couponCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to deactivate coupons: %w", err)
	}
	return nil
}

// IsFileProcessed checks if a file with the given isAdd and filename is already processed.
func (c *couponRepository) IsFileProcessed(ctx context.Context, isAdd bool, filename string) (*models.ProcessedCouponFile, error) {
	filter := bson.M{"$and": []bson.M{{"isadd": isAdd}, {"file_name": filename}}}
	processedFile := &models.ProcessedCouponFile{}
	err := c.processedFilesCollection.FindOne(ctx, filter).Decode(&processedFile)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return processedFile, nil
}

// InsertProcessedFile inserts a record for a processed file.
func (c *couponRepository) InsertProcessedFile(ctx context.Context, file *models.ProcessedCouponFile) error {
	_, err := c.processedFilesCollection.InsertOne(ctx, file)
	if err != nil {
		return fmt.Errorf("failed to insert processed file: %w", err)
	}
	return nil
}

// UpdateProcessingStatus updates the status of an existing ProcessedCouponFile
func (c *couponRepository) UpdateProcessingStatus(ctx context.Context, id, status string, total int64) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"status": status, "coupon_code_counts": total}}
	_, err := c.processedFilesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update processing status: %w", err)
	}
	return nil
}
