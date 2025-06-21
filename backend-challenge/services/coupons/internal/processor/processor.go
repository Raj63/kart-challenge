// Package processor provides file processing functionality for the Coupons service.
// It handles watching directories for coupon files, processing gzipped files,
// and managing batch operations with database persistence.
package processor

import (
	"bufio"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"coupons/internal/config"
	"coupons/internal/repository"
	"coupons/internal/repository/models"
	"library/logger"

	"github.com/google/uuid"

	"github.com/fsnotify/fsnotify"
)

// CouponProcessor handles the processing of coupon files from watched directories.
// It monitors add/remove directories for .gz files, processes them with optimized
// batch operations, and maintains processing state for resume functionality.
type CouponProcessor struct {
	repo            repository.CouponRepository // Repository for database operations
	log             *logger.Logger              // Application logger
	processorConfig *config.ProcessorConfig     // Processor configuration settings
}

// NewCouponProcessor creates a new CouponProcessor instance with the provided dependencies.
// It initializes the processor with repository, configuration, and logger components.
func NewCouponProcessor(repo repository.CouponRepository, processorConfig *config.ProcessorConfig, log *logger.Logger) *CouponProcessor {
	return &CouponProcessor{repo: repo, processorConfig: processorConfig, log: log}
}

// Run starts the directory watcher and processes coupon files as they appear.
// It creates add/remove directories, sets up file system watchers, processes existing files,
// and continuously monitors for new files until the context is cancelled.
func (p *CouponProcessor) Run(ctx context.Context) error {
	addDir := fmt.Sprintf("%s/add", p.processorConfig.DataDirectory)
	removeDir := fmt.Sprintf("%s/remove", p.processorConfig.DataDirectory)
	if err := os.MkdirAll(addDir, 0755); err != nil {
		return fmt.Errorf("failed to create add dir: %w", err)
	}
	if err := os.MkdirAll(removeDir, 0755); err != nil {
		return fmt.Errorf("failed to create remove dir: %w", err)
	}
	p.log.Info("Watching directories: %s, %s", addDir, removeDir)
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer w.Close()
	if err := w.Add(addDir); err != nil {
		return fmt.Errorf("failed to watch add dir: %w", err)
	}
	if err := w.Add(removeDir); err != nil {
		return fmt.Errorf("failed to watch remove dir: %w", err)
	}
	// Initial scan
	p.log.Info("processExistingFiles add-dir: %s", addDir)
	p.processExistingFiles(ctx, addDir, true)
	p.log.Info("processExistingFiles remove-dir: %s", removeDir)
	p.processExistingFiles(ctx, removeDir, false)

	for {
		select {
		case event := <-w.Events:
			if event.Op&(fsnotify.Create|fsnotify.Rename) != 0 {
				if strings.HasSuffix(event.Name, ".gz") {
					if strings.Contains(event.Name, "/add/") {
						p.handleGzFile(ctx, event.Name, true)
					} else if strings.Contains(event.Name, "/remove/") {
						p.handleGzFile(ctx, event.Name, false)
					}
				}
			}
		case err := <-w.Errors:
			p.log.Error("watcher error: %v", err)
		case <-ctx.Done():
			return nil
		}
	}
}

// processExistingFiles processes all .gz files in the given directory.
// It scans the directory for existing files and processes them with the specified
// operation type (add or remove). This is called during startup to handle
// files that may have been placed before the service started.
func (p *CouponProcessor) processExistingFiles(ctx context.Context, dir string, isAdd bool) {
	files, err := filepath.Glob(filepath.Join(dir, "*.gz"))
	if err != nil {
		p.log.Error("failed to list files in %s: %v", dir, err)
		return
	}
	p.log.Info("processExistingFiles %s, files %+v", dir, files)
	for _, f := range files {
		p.handleGzFile(ctx, f, isAdd)
	}
}

// batchProcessor handles database operations for coupon batches.
// It encapsulates the logic for processing batches of coupon codes
// with the appropriate repository operation (add or deactivate).
type batchProcessor struct {
	repo     repository.CouponRepository // Repository for database operations
	isAdd    bool                        // Flag indicating if this is an add operation
	fileName string                      // Name of the source file being processed
	log      *logger.Logger              // Application logger
}

// processBatch processes a batch of coupon codes using the appropriate repository operation.
// If the batch is empty, it returns immediately. Otherwise, it calls either AddCoupons
// or DeactivateCoupons based on the isAdd flag.
func (bp *batchProcessor) processBatch(ctx context.Context, codes []string) error {
	if len(codes) == 0 {
		return nil
	}

	if bp.isAdd {
		return bp.repo.AddCoupons(ctx, bp.fileName, codes)
	}
	return bp.repo.DeactivateCoupons(ctx, bp.fileName, codes)
}

// handleGzFile extracts coupon codes from a .gz file and adds or deactivates them in the database.
// Optimized for large files (1-2 GB) with parallel processing and efficient memory usage.
func (p *CouponProcessor) handleGzFile(ctx context.Context, path string, isAdd bool) {
	p.log.Info("Processing file: %s", path)

	// Open file and compute hash efficiently
	file, err := os.Open(path)
	if err != nil {
		p.log.Error("failed to open file: %v", err)
		return
	}
	defer file.Close()

	// Compute MD5 and size with larger buffer for better performance
	hash := md5.New()
	stat, err := file.Stat()
	if err != nil {
		p.log.Error("failed to stat file: %v", err)
		return
	}
	size := stat.Size()

	// Use larger buffer for hashing
	buf := make([]byte, 64*1024) // 64KB buffer
	if _, err := io.CopyBuffer(hash, file, buf); err != nil {
		p.log.Error("failed to hash file: %v", err)
		return
	}
	md5sum := hex.EncodeToString(hash.Sum(nil))
	file.Seek(0, io.SeekStart) // Reset file pointer for reading

	fileName := filepath.Base(path)
	alreadyProcessed, err := p.repo.IsFileProcessed(ctx, isAdd, fileName)
	if err != nil {
		p.log.Error("failed to check processed files: %v", err)
		return
	}

	var resumeCount int64
	processedFileID := uuid.New().String()
	if alreadyProcessed != nil {
		if alreadyProcessed.Status == "completed" || alreadyProcessed.Status == "initiated" {
			p.log.Info("File %s already processed/under processing, skipping", fileName)
			return
		}

		if alreadyProcessed.Status == "failed" && alreadyProcessed.CouponCodeCount > 0 {
			resumeCount = alreadyProcessed.CouponCodeCount
			processedFileID = alreadyProcessed.ID
			p.log.Info("Resuming %s from line %d", fileName, alreadyProcessed.CouponCodeCount+1)
		}
	}

	// Optimize batch size for large files
	batchSize := 5000 // Increased from 1000
	if p.processorConfig.BatchSize > 0 {
		batchSize = p.processorConfig.BatchSize
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		p.log.Error("failed to create gzip reader: %v", err)
		return
	}
	defer gz.Close()

	// Create batch processor
	bp := &batchProcessor{
		repo:     p.repo,
		isAdd:    isAdd,
		fileName: fileName,
		log:      p.log,
	}

	// Insert processed file record with initiated status
	processed := &models.ProcessedCouponFile{
		ID:              processedFileID,
		MD5Hash:         md5sum,
		FileName:        fileName,
		Size:            size,
		CouponCodeCount: 0,
		Datetime:        time.Now().Unix(),
		IsAdd:           isAdd,
		Status:          "initiated",
	}

	if err := p.repo.InsertProcessedFile(ctx, processed); err != nil {
		p.log.Error("failed to record processed file: %v", err)
	}

	var total int64
	status := "failed"
	defer func() {
		if err := p.repo.UpdateProcessingStatus(ctx, processed.ID, status, total); err != nil {
			p.log.Error("failed to record processed file: %v", err)
		}
	}()

	// Use optimized processing with worker pool
	total, err = p.processFileOptimized(ctx, gz, bp, int(batchSize), resumeCount, fileName)
	if err != nil {
		p.log.Error("failed to process file %s: %v", fileName, err)
		return
	}

	p.log.Info("Processed %d coupons from %s", total, fileName)
	status = "completed"
}

// processFileOptimized processes the file with optimized performance for large files.
// It uses parallel processing with worker pools, larger buffers, and efficient memory
// management to handle files up to 1-2 GB in size. The function returns the total
// number of processed coupons and any error that occurred during processing.
func (p *CouponProcessor) processFileOptimized(ctx context.Context, gz *gzip.Reader, bp *batchProcessor, batchSize int, resumeCount int64, fileName string) (int64, error) {
	var total int64

	// Create channels for batch processing
	batchChan := make(chan []string, 10) // Buffer for 10 batches
	errorChan := make(chan error, 1)

	// Start worker goroutines for database operations
	numWorkers := 4 // Adjust based on your system capabilities
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for codes := range batchChan {
				select {
				case <-ctx.Done():
					return
				default:
					if err := bp.processBatch(ctx, codes); err != nil {
						select {
						case errorChan <- fmt.Errorf("worker %d failed to process batch: %w", workerID, err):
						default:
						}
						return
					}
				}
			}
		}(i)
	}

	// Process file in chunks with larger scanner buffer
	scanner := bufio.NewScanner(gz)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer

	var (
		codes   = make([]string, 0, batchSize)
		lineNum int64
	)

	// Process lines
	for scanner.Scan() {
		lineNum++
		if resumeCount > 0 && lineNum <= resumeCount {
			continue // skip already processed lines
		}

		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			codes = append(codes, line)
			if len(codes) >= batchSize {
				// Send batch to workers
				select {
				case batchChan <- codes:
					total += int64(len(codes))
					codes = make([]string, 0, batchSize) // Pre-allocate new slice
				case err := <-errorChan:
					p.log.Error("failed to process batch: %v", err)
					return total, err
				case <-ctx.Done():
					return total, ctx.Err()
				}
			}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return total, fmt.Errorf("scanner error: %w", err)
	}

	// Process any remaining codes
	if len(codes) > 0 {
		select {
		case batchChan <- codes:
			total += int64(len(codes))
		case err := <-errorChan:
			p.log.Error("failed to process final batch: %v", err)
			return total, err
		case <-ctx.Done():
			return total, ctx.Err()
		}
	}

	// Close batch channel and wait for workers
	close(batchChan)
	wg.Wait()

	// Check for any errors from workers
	select {
	case err := <-errorChan:
		p.log.Error("failed to process batch: %v", err)
		return total, err
	default:
	}

	return total, nil
}
