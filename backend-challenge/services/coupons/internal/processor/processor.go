package processor

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"coupons/internal/config"
	"coupons/internal/repository"
	"coupons/internal/repository/models"
	"library/logger"

	"github.com/google/uuid"

	"github.com/fsnotify/fsnotify"
)

// CouponProcessor watches directories for coupon files and processes them into the database.
type CouponProcessor struct {
	repo            repository.CouponRepository
	log             *logger.Logger
	processorConfig *config.ProcessorConfig
}

// NewCouponProcessor creates a new CouponProcessor.
func NewCouponProcessor(repo repository.CouponRepository, processorConfig *config.ProcessorConfig, log *logger.Logger) *CouponProcessor {
	return &CouponProcessor{repo: repo, processorConfig: processorConfig, log: log}
}

// Run starts the directory watcher and processes coupon files as they appear.
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
	p.processExistingFiles(ctx, addDir, true)
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
func (p *CouponProcessor) processExistingFiles(ctx context.Context, dir string, isAdd bool) {
	files, err := filepath.Glob(filepath.Join(dir, "*.gz"))
	if err != nil {
		p.log.Error("failed to list files in %s: %v", dir, err)
		return
	}
	for _, f := range files {
		p.handleGzFile(ctx, f, isAdd)
	}
}

// handleGzFile extracts coupon codes from a .gz file and adds or deactivates them in the database.
func (p *CouponProcessor) handleGzFile(ctx context.Context, path string, isAdd bool) {
	p.log.Info("Processing file: %s", path)
	file, err := os.Open(path)
	if err != nil {
		p.log.Error("failed to open file: %v", err)
		return
	}
	defer file.Close()

	// Compute MD5 and size
	hash := md5.New()
	stat, err := file.Stat()
	if err != nil {
		p.log.Error("failed to stat file: %v", err)
		return
	}
	size := stat.Size()
	if _, err := io.Copy(hash, file); err != nil {
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
	if alreadyProcessed {
		p.log.Info("File %s already processed, skipping", fileName)
		return
	}

	batchSize := 1000
	if p.processorConfig.BatchSize > 0 {
		batchSize = p.processorConfig.BatchSize
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		p.log.Error("failed to create gzip reader: %v", err)
		return
	}
	defer gz.Close()
	tr := tar.NewReader(gz)

	var (
		codes []string
		total int
	)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			p.log.Error("tar read error: %v", err)
			return
		}
		if hdr.Typeflag == tar.TypeReg {
			buf := new(strings.Builder)
			if _, err := io.Copy(buf, tr); err != nil {
				p.log.Error("failed to read tar file: %v", err)
				return
			}
			for _, line := range strings.Split(buf.String(), "\n") {
				line = strings.TrimSpace(line)
				if line != "" {
					codes = append(codes, line)
					if len(codes) >= batchSize {
						if isAdd {
							err = p.repo.AddCoupons(ctx, fileName, codes)
						} else {
							err = p.repo.DeactivateCoupons(ctx, fileName, codes)
						}
						if err != nil {
							p.log.Error("failed to process batch: %v", err)
							return
						}
						total += len(codes)
						codes = codes[:0]
					}
				}
			}
		}
	}
	// Process any remaining codes
	if len(codes) > 0 {
		if isAdd {
			err = p.repo.AddCoupons(ctx, fileName, codes)
		} else {
			err = p.repo.DeactivateCoupons(ctx, fileName, codes)
		}
		if err != nil {
			p.log.Error("failed to process final batch: %v", err)
			return
		}
		total += len(codes)
	}

	p.log.Info("Processed %d coupons from %s", total, fileName)

	// Insert processed file record
	processed := &models.ProcessedCouponFile{
		ID:              uuid.New().String(),
		MD5Hash:         md5sum,
		FileName:        fileName,
		Size:            size,
		CouponCodeCount: total,
		Datetime:        time.Now().Unix(),
	}
	if err := p.repo.InsertProcessedFile(ctx, processed); err != nil {
		p.log.Error("failed to record processed file: %v", err)
	}
}
