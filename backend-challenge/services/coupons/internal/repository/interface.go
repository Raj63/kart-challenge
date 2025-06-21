// Package repository provides data access layer interfaces for the Coupons processor service.
package repository

import (
	"context"

	"coupons/internal/repository/models"
)

// CouponRepository defines methods for managing coupons in the database.
// It provides operations for adding and deactivating coupon codes, as well as
// tracking processed files to support resume functionality and prevent duplicate processing.
type CouponRepository interface {
	// AddCoupons adds a batch of coupon codes to the database with the specified filename.
	// The filename is used for tracking which file the coupons came from.
	AddCoupons(ctx context.Context, fileName string, codes []string) error

	// DeactivateCoupons deactivates a batch of coupon codes in the database.
	// This operation marks coupons as inactive rather than deleting them.
	DeactivateCoupons(ctx context.Context, fileName string, codes []string) error

	// IsFileProcessed checks if a file has already been processed by querying the database.
	// Returns the processed file record if found, nil otherwise.
	// The isAdd parameter distinguishes between add and remove operations.
	IsFileProcessed(ctx context.Context, isAdd bool, filename string) (*models.ProcessedCouponFile, error)

	// InsertProcessedFile records a new processed file in the database.
	// This is used to track file processing status and support resume functionality.
	InsertProcessedFile(ctx context.Context, file *models.ProcessedCouponFile) error

	// UpdateProcessingStatus updates the status and coupon count of a processed file.
	// Used to track processing progress and final status (completed, failed, etc.).
	UpdateProcessingStatus(ctx context.Context, id, status string, total int64) error
}
