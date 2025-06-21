package repository

import (
	"context"

	"coupons/internal/repository/models"
)

// CouponRepository defines methods for managing coupons in the database.
type CouponRepository interface {
	AddCoupons(ctx context.Context, fileName string, codes []string) error
	DeactivateCoupons(ctx context.Context, fileName string, codes []string) error
	IsFileProcessed(ctx context.Context, isAdd bool, filename string) (bool, error)
	InsertProcessedFile(ctx context.Context, file *models.ProcessedCouponFile) error
}
