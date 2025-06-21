package config

import (
	"coupons/internal/constants"
	"library/config"
	"library/logger"
)

// Config holds the application configuration.
type Config struct {
	Env       string            `json:"env"`       // Environment name (e.g., "development", "production")
	Logger    *logger.LogConfig `json:"logger"`    // Logger configuration
	Database  *DbConfig         `json:"database"`  // Database configuration
	Processor *ProcessorConfig  `json:"processor"` // Processor configuration
}

// DbConfig holds database connection configuration.
type DbConfig struct {
	Host         string `json:"host"`     // Database host address
	Port         int    `json:"port"`     // Database port number
	User         string `json:"user"`     // Database username
	Password     string `json:"password"` // Database password
	DatabaseName string `json:"dbname"`   // Database name
	Type         string `json:"type"`     // Database type (e.g., "mongodb")
}

// ProcessorConfig holds processor-specific configuration.
type ProcessorConfig struct {
	BatchSize     int    `json:"batch_size"`     // Number of coupon codes to process in each batch
	DataDirectory string `json:"data_directory"` // Directory to watch for coupon files
}

// NewConfig creates a new Config from a config manager.
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
			BatchSize: configManager.GetInt("processor.batch_size"),
		},
	}

	cfg.Logger.Version = constants.Version
	cfg.Logger.Commit = constants.CommitHash

	return &cfg, nil
}
