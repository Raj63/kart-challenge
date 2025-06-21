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
	collection     *mongo.Collection
	processedFiles *mongo.Collection
}

// NewCouponRepository creates a new CouponRepository using the given MongoDB database.
func NewCouponRepository(repo *Repository) CouponRepository {
	return &couponRepository{
		collection:     repo.db.Collection("coupons"),
		processedFiles: repo.db.Collection("processed-coupon-files"),
	}
}

// AddCoupons inserts new coupon codes into the database as active.
// If a coupon with the same code and filename already exists, it updates the datetime instead of inserting a duplicate.
func (r *couponRepository) AddCoupons(ctx context.Context, fileName string, codes []string) error {
	if len(codes) == 0 {
		return nil
	}

	var models []mongo.WriteModel
	now := time.Now().Unix()

	for _, code := range codes {
		filter := bson.M{"couponcode": code, "filename": fileName}
		update := bson.M{
			"$set": bson.M{
				"datetime": now,
				"isactive": true,
			},
			"$setOnInsert": bson.M{
				"id":         uuid.New().String(),
				"couponcode": code,
				"filename":   fileName,
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

	_, err := r.collection.BulkWrite(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to upsert coupons: %w", err)
	}

	return nil
}

// DeactivateCoupons marks coupon codes as inactive in the database.
func (r *couponRepository) DeactivateCoupons(ctx context.Context, fileName string, codes []string) error {
	filter := bson.M{"file_name": fileName, "coupon_code": bson.M{"$in": codes}}
	update := bson.M{"$set": bson.M{"isactive": false}}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to deactivate coupons: %w", err)
	}
	return nil
}

// IsFileProcessed checks if a file with the given isAdd and filename is already processed.
func (r *couponRepository) IsFileProcessed(ctx context.Context, isAdd bool, filename string) (bool, error) {
	filter := bson.M{"$and": []bson.M{{"isadd": isAdd}, {"filename": filename}}}
	count, err := r.processedFiles.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check processed files: %w", err)
	}
	return count > 0, nil
}

// InsertProcessedFile inserts a record for a processed file.
func (r *couponRepository) InsertProcessedFile(ctx context.Context, file *models.ProcessedCouponFile) error {
	_, err := r.processedFiles.InsertOne(ctx, file)
	if err != nil {
		return fmt.Errorf("failed to insert processed file: %w", err)
	}
	return nil
}
