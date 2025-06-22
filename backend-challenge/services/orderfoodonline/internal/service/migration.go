package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"library/logger"
	"orderfoodonline/internal/repository"
	"orderfoodonline/internal/repository/models"

	"github.com/google/uuid"
)

// MigrationService handles database migrations and seeding.
type MigrationService struct {
	repo repository.ProductRepository
	log  logger.ILogger
}

// NewMigrationService creates a new MigrationService.
func NewMigrationService(repo repository.ProductRepository, log logger.ILogger) *MigrationService {
	return &MigrationService{
		repo: repo,
		log:  log,
	}
}

// RunMigrations executes all pending migrations in order.
func (m *MigrationService) RunMigrations(ctx context.Context, migrationsDir string) error {
	m.log.Info("Starting database migrations from directory: %s", migrationsDir)

	// Get all migration files
	migrationFiles, err := m.getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	if len(migrationFiles) == 0 {
		m.log.Info("No migration files found")
		return nil
	}

	// Sort migration files by version
	sort.Strings(migrationFiles)

	// Get applied migrations
	appliedMigrations, err := m.repo.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Execute pending migrations
	for _, migrationFile := range migrationFiles {
		version := m.extractVersionFromFilename(migrationFile)

		// Check if migration is already applied
		if m.isMigrationApplied(appliedMigrations, version) {
			m.log.Info("Migration %s already applied, skipping", version)
			continue
		}

		m.log.Info("Applying migration: %s", migrationFile)
		if err := m.applyMigration(ctx, filepath.Join(migrationsDir, migrationFile)); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migrationFile, err)
		}
	}

	m.log.Info("Database migrations completed successfully")
	return nil
}

// getMigrationFiles returns a list of migration JSON files sorted by version.
func (m *MigrationService) getMigrationFiles(migrationsDir string) ([]string, error) {
	var migrationFiles []string

	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			// Check if filename matches migration pattern (e.g., 0001_init_product_catalog.json)
			if m.isMigrationFilename(d.Name()) {
				migrationFiles = append(migrationFiles, d.Name())
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk migrations directory: %w", err)
	}

	return migrationFiles, nil
}

// isMigrationFilename checks if a filename follows the migration naming pattern.
func (m *MigrationService) isMigrationFilename(filename string) bool {
	// Pattern: 0001_init_product_catalog.json
	parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
	if len(parts) < 2 {
		return false
	}

	// First part should be a 4-digit number
	if len(parts[0]) != 4 {
		return false
	}

	_, err := strconv.Atoi(parts[0])
	return err == nil
}

// extractVersionFromFilename extracts the version number from a migration filename.
func (m *MigrationService) extractVersionFromFilename(filename string) string {
	parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// isMigrationApplied checks if a migration version has already been applied.
func (m *MigrationService) isMigrationApplied(appliedMigrations []models.Migration, version string) bool {
	for _, migration := range appliedMigrations {
		if migration.Version == version {
			return true
		}
	}
	return false
}

// applyMigration applies a single migration file.
func (m *MigrationService) applyMigration(ctx context.Context, migrationPath string) error {
	// Read migration file
	// #nosec G304
	data, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Parse migration data
	var migrationData models.MigrationData
	if err := json.Unmarshal(data, &migrationData); err != nil {
		return fmt.Errorf("failed to parse migration file: %w", err)
	}

	// Create migration record
	migration := &models.Migration{
		ID:          uuid.New().String(),
		Version:     migrationData.Version,
		Description: migrationData.Description,
		AppliedAt:   time.Now().Unix(),
		Status:      "in_progress",
	}

	// Insert migration record
	if err := m.repo.InsertMigration(ctx, migration); err != nil {
		return fmt.Errorf("failed to insert migration record: %w", err)
	}

	// Apply the migration (seed products)
	if err := m.seedProducts(ctx, migrationData.Products); err != nil {
		// Update migration status to failed
		migration.Status = "failed"
		if updateErr := m.repo.UpdateMigration(ctx, migration); updateErr != nil {
			m.log.Error("Failed to update migration status: %v", updateErr)
		}
		return fmt.Errorf("failed to seed products: %w", err)
	}

	// Update migration status to completed
	migration.Status = "completed"
	if err := m.repo.UpdateMigration(ctx, migration); err != nil {
		return fmt.Errorf("failed to update migration status: %w", err)
	}

	m.log.Info("Successfully applied migration %s: %s", migrationData.Version, migrationData.Description)
	return nil
}

// seedProducts seeds the database with product data.
func (m *MigrationService) seedProducts(ctx context.Context, products []models.Product) error {
	if len(products) == 0 {
		return nil
	}

	m.log.Info("Seeding %d products", len(products))

	// Insert products in batches
	batchSize := 100
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}

		batch := products[i:end]
		if err := m.repo.BulkInsertProducts(ctx, batch); err != nil {
			return fmt.Errorf("failed to insert product batch: %w", err)
		}

		m.log.Info("Inserted batch of %d products", len(batch))
	}

	return nil
}
