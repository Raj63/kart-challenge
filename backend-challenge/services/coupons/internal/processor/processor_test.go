package processor

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"coupons/internal/config"
	"coupons/internal/repository"
	"coupons/internal/repository/mocks"
	"coupons/internal/repository/models"

	libmocks "library/logger/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewCouponProcessor(t *testing.T) {
	// Given: Mock dependencies
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	config := &config.ProcessorConfig{
		DataDirectory: "/test/data",
		BatchSize:     1000,
	}

	// When: Creating a new coupon processor
	processor := NewCouponProcessor(mockRepo, config, mockLogger)

	// Then: Processor should be created successfully
	assert.NotNil(t, processor)
	assert.Equal(t, mockRepo, processor.repo)
	assert.Equal(t, config, processor.processorConfig)
	assert.Equal(t, mockLogger, processor.log)
}

func TestBatchProcessor_ProcessBatch_AddOperation(t *testing.T) {
	// Given: A batch processor for add operation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	bp := &batchProcessor{
		repo:     mockRepo,
		isAdd:    true,
		fileName: "test-file.gz",
		log:      mockLogger,
	}

	codes := []string{"COUPON1", "COUPON2", "COUPON3"}
	ctx := context.Background()

	// Mock repository behavior
	mockRepo.EXPECT().AddCoupons(ctx, "test-file.gz", codes).Return(nil)

	// When: Processing batch for add operation
	err := bp.processBatch(ctx, codes)

	// Then: Should succeed without error
	require.NoError(t, err)
}

func TestBatchProcessor_ProcessBatch_RemoveOperation(t *testing.T) {
	// Given: A batch processor for remove operation
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	bp := &batchProcessor{
		repo:     mockRepo,
		isAdd:    false,
		fileName: "test-file.gz",
		log:      mockLogger,
	}

	codes := []string{"COUPON1", "COUPON2", "COUPON3"}
	ctx := context.Background()

	// Mock repository behavior
	mockRepo.EXPECT().DeactivateCoupons(ctx, "test-file.gz", codes).Return(nil)

	// When: Processing batch for remove operation
	err := bp.processBatch(ctx, codes)

	// Then: Should succeed without error
	require.NoError(t, err)
}

func TestBatchProcessor_ProcessBatch_EmptyCodes(t *testing.T) {
	// Given: A batch processor with empty codes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	bp := &batchProcessor{
		repo:     mockRepo,
		isAdd:    true,
		fileName: "test-file.gz",
		log:      mockLogger,
	}

	emptyCodes := []string{}
	ctx := context.Background()

	// When: Processing empty batch
	err := bp.processBatch(ctx, emptyCodes)

	// Then: Should succeed without error (early return)
	require.NoError(t, err)
}

func TestBatchProcessor_ProcessBatch_DatabaseError(t *testing.T) {
	// Given: A batch processor that returns database error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	bp := &batchProcessor{
		repo:     mockRepo,
		isAdd:    true,
		fileName: "test-file.gz",
		log:      mockLogger,
	}

	codes := []string{"COUPON1", "COUPON2"}
	ctx := context.Background()
	expectedError := errors.New("database connection failed")

	// Mock repository behavior to return error
	mockRepo.EXPECT().AddCoupons(ctx, "test-file.gz", codes).Return(expectedError)

	// When: Processing batch with database error
	err := bp.processBatch(ctx, codes)

	// Then: Should return error
	require.Error(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())
}

func createGzipFile(filePath string, promoCodeCount int) ([]string, error) {
	// Create the file
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(f)
	defer gzWriter.Close()

	codes := make([]string, promoCodeCount)
	for i := 0; i < promoCodeCount; i++ {
		codes[i] = fmt.Sprintf("COUPONS%d", i)
		// Write the promo code (with newline)
		_, err = gzWriter.Write([]byte(codes[i] + "\n"))
		if err != nil {
			return nil, fmt.Errorf("failed to write to gzip file: %v", err)
		}
	}

	return codes, nil
}

func TestCouponProcessor_HandleGzFile_FileAlreadyProcessed(t *testing.T) {
	// Given: A coupon processor with already processed file
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	//prepare temp gzip file for processing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	fileName := "test-file.gz"
	filePath := filepath.Join(tmpDir, fileName)

	codes, err := createGzipFile(filePath, 10)
	require.NoError(t, err)
	assert.Equal(t, 10, len(codes))

	config := &config.ProcessorConfig{
		DataDirectory: tmpDir,
		BatchSize:     1000,
	}
	processor := NewCouponProcessor(mockRepo, config, mockLogger)

	ctx := context.Background()

	isAdd := true

	// Mock already processed file
	alreadyProcessed := &models.ProcessedCouponFile{
		ID:              "file-123",
		MD5Hash:         "abc123",
		FileName:        fileName,
		IsAdd:           isAdd,
		Status:          "completed",
		CouponCodeCount: 10,
		Datetime:        time.Now().Unix(),
	}

	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(2)
	// Mock repository behaviors
	mockRepo.EXPECT().IsFileProcessed(ctx, isAdd, fileName).Return(alreadyProcessed, nil)

	// When: Handling already processed file
	processor.handleGzFile(ctx, filePath, isAdd)

	// Then: Should skip processing
	// No additional assertions needed as the method should return early
}

func TestCouponProcessor_HandleGzFile_FileUnderProcessing(t *testing.T) {
	// Given: A coupon processor with file under processing
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	//prepare temp gzip file for processing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	fileName := "test-file.gz"
	filePath := filepath.Join(tmpDir, fileName)

	codes, err := createGzipFile(filePath, 10)
	require.NoError(t, err)
	assert.Equal(t, 10, len(codes))

	config := &config.ProcessorConfig{
		DataDirectory: tmpDir,
		BatchSize:     1000,
	}

	processor := NewCouponProcessor(mockRepo, config, mockLogger)

	ctx := context.Background()
	isAdd := true

	// Mock file under processing
	underProcessing := &models.ProcessedCouponFile{
		ID:              "file-123",
		MD5Hash:         "abc123",
		FileName:        fileName,
		IsAdd:           isAdd,
		Status:          "initiated",
		CouponCodeCount: 10,
		Datetime:        time.Now().Unix(),
	}

	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(2)
	// Mock repository behaviors
	mockRepo.EXPECT().IsFileProcessed(ctx, isAdd, fileName).Return(underProcessing, nil)

	// When: Handling file under processing
	processor.handleGzFile(ctx, filePath, isAdd)

	// Then: Should skip processing
	// No additional assertions needed as the method should return early
}

func TestCouponProcessor_HandleGzFile_ResumeFromFailure(t *testing.T) {
	// Given: A coupon processor with failed file to resume
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	//prepare temp gzip file for processing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	fileName := "test-file.gz"
	filePath := filepath.Join(tmpDir, fileName)

	codes, err := createGzipFile(filePath, 10)
	require.NoError(t, err)
	assert.Equal(t, 10, len(codes))

	config := &config.ProcessorConfig{
		DataDirectory: tmpDir,
		BatchSize:     1000,
	}

	processor := NewCouponProcessor(mockRepo, config, mockLogger)

	ctx := context.Background()
	isAdd := true

	// Mock failed file with partial processing
	failedFile := &models.ProcessedCouponFile{
		ID:              "file-123",
		MD5Hash:         "abc123",
		FileName:        fileName,
		IsAdd:           isAdd,
		Status:          "failed",
		Size:            163,
		CouponCodeCount: 5, // Resume from line 6
		Datetime:        time.Now().Unix(),
	}

	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(3)
	// Mock repository behaviors for resume scenario
	mockRepo.EXPECT().IsFileProcessed(ctx, isAdd, fileName).Return(failedFile, nil)
	mockRepo.EXPECT().UpdateProcessingStatus(ctx, "file-123", "initiated", int64(5)).Return(nil) // resume update with resume count
	mockRepo.EXPECT().AddCoupons(ctx, fileName, codes[5:]).Return(nil)
	mockRepo.EXPECT().UpdateProcessingStatus(ctx, "file-123", "completed", int64(10)).Return(nil) // final update with total count
	// Note: This test would need a real gzip file to test the full resume functionality
	// For now, we just verify the resume logic is triggered

	// When: Handling failed file for resume
	processor.handleGzFile(ctx, filePath, isAdd)

	// Then: Should attempt to resume processing
	// No additional assertions needed as this is just testing the method can be called
}

func TestCouponProcessor_HandleGzFile_NewFileProcessing(t *testing.T) {
	// Given: A coupon processor with new file to process
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	//prepare temp gzip file for processing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	fileName := "test-file.gz"
	filePath := filepath.Join(tmpDir, fileName)

	codes, err := createGzipFile(filePath, 10)
	require.NoError(t, err)
	assert.Equal(t, 10, len(codes))

	config := &config.ProcessorConfig{
		DataDirectory: tmpDir,
		BatchSize:     1000,
	}

	processor := NewCouponProcessor(mockRepo, config, mockLogger)

	ctx := context.Background()
	isAdd := true

	// Mock new file (not processed before)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(2)
	mockRepo.EXPECT().IsFileProcessed(ctx, isAdd, fileName).Return(nil, nil)
	mockRepo.EXPECT().InsertProcessedFile(ctx, gomock.Any()).Return(nil) // fresh insert as no record exists
	mockRepo.EXPECT().AddCoupons(ctx, fileName, codes).Return(nil)
	mockRepo.EXPECT().UpdateProcessingStatus(ctx, gomock.Any(), "completed", int64(10)).Return(nil) // final update with total count
	// Note: This test would need a real gzip file to test the full processing functionality
	// For now, we just verify the new file logic is triggered

	// When: Handling new file
	processor.handleGzFile(ctx, filePath, isAdd)

	// Then: Should attempt to process new file
	// No additional assertions needed as this is just testing the method can be called
}

func TestCouponProcessor_HandleGzFile_DatabaseError(t *testing.T) {
	// Given: A coupon processor with database error
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	//prepare temp gzip file for processing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	fileName := "test-file.gz"
	filePath := filepath.Join(tmpDir, fileName)

	codes, err := createGzipFile(filePath, 10)
	require.NoError(t, err)
	assert.Equal(t, 10, len(codes))

	config := &config.ProcessorConfig{
		DataDirectory: tmpDir,
		BatchSize:     1000,
	}

	processor := NewCouponProcessor(mockRepo, config, mockLogger)

	ctx := context.Background()
	isAdd := true
	expectedError := errors.New("database connection failed")

	// Mock repository behavior to return error
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(1)
	mockLogger.EXPECT().Error("failed to check processed files: %v", expectedError) // the error format with the expected error
	mockRepo.EXPECT().IsFileProcessed(ctx, isAdd, fileName).Return(nil, expectedError)

	// When: Handling file with database error
	processor.handleGzFile(ctx, filePath, isAdd)

	// Then: Should handle error gracefully
	// No additional assertions needed as this is just testing the method can be called
}

func TestCouponProcessor_ProcessExistingFiles_EmptyDirectory(t *testing.T) {
	// Given: A coupon processor with empty directory
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	//prepare temp gzip file for processing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	config := &config.ProcessorConfig{
		DataDirectory: tmpDir,
		BatchSize:     1000,
	}
	processor := NewCouponProcessor(mockRepo, config, mockLogger)

	ctx := context.Background()
	isAdd := true

	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(1)
	// When: Processing existing files in empty directory
	processor.processExistingFiles(ctx, tmpDir, isAdd)

	// Then: Should handle empty directory gracefully
	// No assertions needed as this is just testing the method can be called
}

func TestCouponProcessor_ProcessExistingFiles_WithFiles(t *testing.T) {
	// Given: A coupon processor with files in directory
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)
	//prepare temp gzip file for processing
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	fileName := "test-file.gz"
	filePath := filepath.Join(tmpDir, fileName)

	codes, err := createGzipFile(filePath, 10)
	require.NoError(t, err)
	assert.Equal(t, 10, len(codes))

	config := &config.ProcessorConfig{
		DataDirectory: tmpDir,
		BatchSize:     1000,
	}

	processor := NewCouponProcessor(mockRepo, config, mockLogger)

	ctx := context.Background()
	isAdd := true

	// Note: This test would need to create actual files or mock filepath.Glob
	// For now, we just verify the method can be called
	// Mock new file (not processed before)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(3)
	mockRepo.EXPECT().IsFileProcessed(ctx, isAdd, fileName).Return(nil, nil)
	mockRepo.EXPECT().InsertProcessedFile(ctx, gomock.Any()).Return(nil) // fresh insert as no record exists
	mockRepo.EXPECT().AddCoupons(ctx, fileName, codes).Return(nil)
	mockRepo.EXPECT().UpdateProcessingStatus(ctx, gomock.Any(), "completed", int64(10)).Return(nil) // final update with total count

	// When: Processing existing files
	processor.processExistingFiles(ctx, tmpDir, isAdd)

	// Then: Should process files if they exist
	// No assertions needed as this is just testing the method can be called
}

func TestCouponProcessor_InterfaceCompliance(t *testing.T) {
	// Given: A mock coupon repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)

	// When: Checking interface compliance
	var repo repository.CouponRepository = mockRepo

	// Then: Should implement the interface
	assert.NotNil(t, repo)

	// Test that all methods can be called
	ctx := context.Background()

	// Test AddCoupons
	mockRepo.EXPECT().AddCoupons(ctx, "test.gz", []string{"CODE1"}).Return(nil)
	err := repo.AddCoupons(ctx, "test.gz", []string{"CODE1"})
	require.NoError(t, err)

	// Test DeactivateCoupons
	mockRepo.EXPECT().DeactivateCoupons(ctx, "test.gz", []string{"CODE1"}).Return(nil)
	err = repo.DeactivateCoupons(ctx, "test.gz", []string{"CODE1"})
	require.NoError(t, err)

	// Test IsFileProcessed
	mockRepo.EXPECT().IsFileProcessed(ctx, true, "test.gz").Return(nil, nil)
	file, err := repo.IsFileProcessed(ctx, true, "test.gz")
	require.NoError(t, err)
	assert.Nil(t, file)

	// Test InsertProcessedFile
	testFile := &models.ProcessedCouponFile{ID: "test"}
	mockRepo.EXPECT().InsertProcessedFile(ctx, testFile).Return(nil)
	err = repo.InsertProcessedFile(ctx, testFile)
	require.NoError(t, err)

	// Test UpdateProcessingStatus
	mockRepo.EXPECT().UpdateProcessingStatus(ctx, "test", "completed", int64(10)).Return(nil)
	err = repo.UpdateProcessingStatus(ctx, "test", "completed", 10)
	require.NoError(t, err)
}

func TestCouponProcessor_ConfigurationHandling(t *testing.T) {
	// Given: Different processor configurations
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCouponRepository(ctrl)
	mockLogger := libmocks.NewMockILogger(ctrl)

	testCases := []struct {
		name          string
		config        *config.ProcessorConfig
		expectedBatch int
	}{
		{
			name: "Default batch size",
			config: &config.ProcessorConfig{
				DataDirectory: "/test/data",
				BatchSize:     0, // Should use default
			},
			expectedBatch: 5000, // Default value
		},
		{
			name: "Custom batch size",
			config: &config.ProcessorConfig{
				DataDirectory: "/test/data",
				BatchSize:     2000,
			},
			expectedBatch: 2000,
		},
		{
			name: "Large batch size",
			config: &config.ProcessorConfig{
				DataDirectory: "/test/data",
				BatchSize:     10000,
			},
			expectedBatch: 10000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// When: Creating processor with specific config
			processor := NewCouponProcessor(mockRepo, tc.config, mockLogger)

			// Then: Should use correct batch size
			assert.Equal(t, tc.expectedBatch, processor.processorConfig.BatchSize)
		})
	}
}
