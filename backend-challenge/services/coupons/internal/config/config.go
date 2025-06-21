// Package config provides configuration management for the Coupons processor service.
package config

import (
	"coupons/internal/constants"
	"library/config"
	"library/logger"
)

// Config holds the complete application configuration for the Coupons processor service.
// It includes environment settings, logging configuration, database connection details,
// and processor-specific settings for batch processing and file monitoring.
type Config struct {
	Env       string            `json:"env"`       // Environment name (e.g., "development", "production")
	Logger    *logger.LogConfig `json:"logger"`    // Logger configuration with version and commit info
	Database  *DbConfig         `json:"database"`  // Database connection configuration
	Processor *ProcessorConfig  `json:"processor"` // Processor configuration for file processing
}

// DbConfig holds database connection configuration including credentials and connection details.
// It specifies the database type, host, port, and authentication information.
type DbConfig struct {
	Host         string `json:"host"`     // Database host address
	Port         int    `json:"port"`     // Database port number
	User         string `json:"user"`     // Database username
	Password     string `json:"password"` // Database password
	DatabaseName string `json:"dbname"`   // Database name
	Type         string `json:"type"`     // Database type (e.g., "mongodb", "postgres")
}

// ProcessorConfig holds processor-specific configuration for coupon file processing.
// It defines batch processing parameters and file monitoring directories.
type ProcessorConfig struct {
	BatchSize     int    `json:"batch_size"`     // Number of coupon codes to process in each batch
	DataDirectory string `json:"data_directory"` // Directory to watch for coupon files (add/remove subdirectories)
}

// NewConfig creates a new Config instance from a configuration manager.
// It populates all configuration fields from the provided config manager
// and sets version information from constants for logging.
func NewConfig(configManager *config.Manager) (*Config, error) {
	cfg := Config{
		Env:    configManager.GetString("env"),
		Logger: configManager.GetLogConfig(),
		Database: &DbConfig{
			Host:         configManager.GetString("database.host"),
			Port:         configManager.GetInt("database.port"),
			User:         configManager.GetString("database.user"),
			Password:     configManager.GetString("database.password"),
			DatabaseName: configManager.GetString("database.dbname"),
			Type:         configManager.GetString("database.type"),
		},
		Processor: &ProcessorConfig{
			BatchSize:     configManager.GetInt("processor.batch_size"),
			DataDirectory: configManager.GetString("processor.data_directory"),
		},
	}

	cfg.Logger.Version = constants.Version
	cfg.Logger.Commit = constants.CommitHash

	return &cfg, nil
}
