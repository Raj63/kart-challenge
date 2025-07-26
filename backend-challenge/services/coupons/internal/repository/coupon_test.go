package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"coupons/internal/repository/mocks"
	"coupons/internal/repository/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"
)

func TestCouponRepository_AddCoupons_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	codes := []string{"COUPON1", "COUPON2", "COUPON3"}
	fileName := "test-file.gz"
	ctx := context.Background()

	// Mock collection behavior for bulk write
	mockCouponCollection.EXPECT().BulkWrite(ctx, gomock.Any()).Return(&mongo.BulkWriteResult{}, nil)

	// When: Adding coupons
	err := repo.AddCoupons(ctx, fileName, codes)

	// Then: Should succeed without error
	require.NoError(t, err)
}

func TestNewCouponRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	// Then: Repository should be created successfully
	assert.NotNil(t, repo)
	assert.IsType(t, &couponRepository{}, repo)
}

func TestCouponRepository_AddCoupons_EmptyCodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}
	emptyCodes := []string{}
	fileName := "test-file.gz"
	ctx := context.Background()

	// When: Adding empty coupon codes
	err := repo.AddCoupons(ctx, fileName, emptyCodes)

	// Then: Should succeed without error (early return)
	require.NoError(t, err)
	// Should not call BulkWrite for empty codes
}

func TestCouponRepository_AddCoupons_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	codes := []string{"COUPON1", "COUPON2"}
	fileName := "test-file.gz"
	ctx := context.Background()

	expectedError := errors.New("database connection failed")

	// Mock collection behavior to return error
	mockCouponCollection.EXPECT().BulkWrite(ctx, gomock.Any()).Return(nil, expectedError)

	// When: Adding coupons with database error
	err := repo.AddCoupons(ctx, fileName, codes)

	// Then: Should return error
	require.Error(t, err)
	assert.ErrorContains(t, err, expectedError.Error())
}

func TestCouponRepository_AddCoupons_SingleCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	codes := []string{"SINGLECOUPON"}
	fileName := "test-file.gz"
	ctx := context.Background()

	// Mock collection behavior for bulk write
	mockCouponCollection.EXPECT().BulkWrite(ctx, gomock.Any()).Return(&mongo.BulkWriteResult{}, nil)

	// When: Adding single coupon
	err := repo.AddCoupons(ctx, fileName, codes)

	// Then: Should succeed without error
	require.NoError(t, err)
}

func TestCouponRepository_DeactivateCoupons_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	codes := []string{"COUPON1", "COUPON2", "COUPON3"}
	fileName := "test-file.gz"
	ctx := context.Background()

	// Mock collection behavior for update many
	mockCouponCollection.EXPECT().UpdateMany(ctx, gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)

	// When: Deactivating coupons
	err := repo.DeactivateCoupons(ctx, fileName, codes)

	// Then: Should succeed without error
	require.NoError(t, err)
}

func TestCouponRepository_DeactivateCoupons_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	codes := []string{"COUPON1", "COUPON2"}
	fileName := "test-file.gz"
	ctx := context.Background()

	expectedError := errors.New("database connection failed")

	// Mock collection behavior to return error
	mockCouponCollection.EXPECT().UpdateMany(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedError)

	// When: Deactivating coupons with database error
	err := repo.DeactivateCoupons(ctx, fileName, codes)

	// Then: Should return error
	require.Error(t, err)
	assert.ErrorContains(t, err, expectedError.Error())
}

func TestCouponRepository_IsFileProcessed_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	fileName := "test-file.gz"
	isAdd := true
	ctx := context.Background()

	expectedError := errors.New("Registry cannot be nil")

	// Mock collection behavior to return error
	mockProcessedFilesCollection.EXPECT().FindOne(ctx, gomock.Any()).Return(&mongo.SingleResult{})

	// When: Checking if file is processed with database error
	file, err := repo.IsFileProcessed(ctx, isAdd, fileName)

	// Then: Should return error and nil file
	require.Error(t, err)
	assert.Nil(t, file)
	assert.ErrorContains(t, err, expectedError.Error())
}

func TestCouponRepository_InsertProcessedFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}
	file := &models.ProcessedCouponFile{
		ID:              "file-123",
		MD5Hash:         "abc123",
		FileName:        "test-file.gz",
		IsAdd:           true,
		Size:            1024,
		CouponCodeCount: 0,
		Datetime:        time.Now().Unix(),
		Status:          "initiated",
	}

	ctx := context.Background()

	// Mock collection behavior
	mockProcessedFilesCollection.EXPECT().InsertOne(ctx, file).Return(&mongo.InsertOneResult{}, nil)

	// When: Inserting processed file
	err := repo.InsertProcessedFile(ctx, file)

	// Then: Should succeed without error
	require.NoError(t, err)
}

func TestCouponRepository_InsertProcessedFile_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}
	file := &models.ProcessedCouponFile{
		ID:              "file-123",
		MD5Hash:         "abc123",
		FileName:        "test-file.gz",
		IsAdd:           true,
		Size:            1024,
		CouponCodeCount: 0,
		Datetime:        time.Now().Unix(),
		Status:          "initiated",
	}

	expectedError := errors.New("database connection failed")
	ctx := context.Background()

	// Mock collection behavior to return error
	mockProcessedFilesCollection.EXPECT().InsertOne(ctx, file).Return(nil, expectedError)

	// When: Inserting processed file with database error
	err := repo.InsertProcessedFile(ctx, file)

	// Then: Should return error
	require.Error(t, err)
	assert.ErrorContains(t, err, expectedError.Error())
}

func TestCouponRepository_UpdateProcessingStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	fileID := "file-123"
	status := "completed"
	total := int64(100)
	ctx := context.Background()

	// Mock collection behavior
	mockProcessedFilesCollection.EXPECT().UpdateOne(ctx, gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)

	// When: Updating processing status
	err := repo.UpdateProcessingStatus(ctx, fileID, status, total)

	// Then: Should succeed without error
	require.NoError(t, err)
}

func TestCouponRepository_UpdateProcessingStatus_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	fileID := "file-123"
	status := "failed"
	total := int64(50)
	ctx := context.Background()

	expectedError := errors.New("database connection failed")

	// Mock collection behavior to return error
	mockProcessedFilesCollection.EXPECT().UpdateOne(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedError)

	// When: Updating processing status with database error
	err := repo.UpdateProcessingStatus(ctx, fileID, status, total)

	// Then: Should return error
	require.Error(t, err)
	assert.ErrorContains(t, err, expectedError.Error())
}

func TestCouponRepository_InterfaceCompliance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Given: A coupon repository with mock dependencies
	mockCouponCollection := mocks.NewMockCollection(ctrl)
	mockProcessedFilesCollection := mocks.NewMockCollection(ctrl)

	repo := &couponRepository{
		couponCollection:         mockCouponCollection,
		processedFilesCollection: mockProcessedFilesCollection,
	}

	// When: Checking interface compliance
	var couponRepo CouponRepository = repo

	// Then: Should implement the interface
	assert.NotNil(t, couponRepo)

	// Test that all methods can be called
	ctx := context.Background()

	// Test AddCoupons
	mockCouponCollection.EXPECT().BulkWrite(ctx, gomock.Any()).Return(&mongo.BulkWriteResult{}, nil)
	err := couponRepo.AddCoupons(ctx, "test.gz", []string{"CODE1"})
	require.NoError(t, err)

	// Test DeactivateCoupons
	mockCouponCollection.EXPECT().UpdateMany(ctx, gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	err = couponRepo.DeactivateCoupons(ctx, "test.gz", []string{"CODE1"})
	require.NoError(t, err)

	// Test InsertProcessedFile
	testFile := &models.ProcessedCouponFile{ID: "test"}
	mockProcessedFilesCollection.EXPECT().InsertOne(ctx, testFile).Return(&mongo.InsertOneResult{}, nil)
	err = couponRepo.InsertProcessedFile(ctx, testFile)
	require.NoError(t, err)

	// Test UpdateProcessingStatus
	mockProcessedFilesCollection.EXPECT().UpdateOne(ctx, gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	err = couponRepo.UpdateProcessingStatus(ctx, "test", "completed", 10)
	require.NoError(t, err)
}
