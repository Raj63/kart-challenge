// Package config provides configuration management for the Order Food Online service.
package config

import (
	"library/config"
	"library/logger"
	"orderfoodonline/internal/constants"
	"time"
)

// Config holds the complete application configuration including environment,
// server settings, database connection, logging, and Swagger documentation.
type Config struct {
	Env      string            `json:"env"`      // Environment name (e.g., "development", "production")
	Swagger  *SwaggerConfig    `json:"swagger"`  // Swagger documentation configuration
	Server   *ServerConfig     `json:"server"`   // HTTP server configuration
	Logger   *logger.LogConfig `json:"logger"`   // Logging configuration
	Database *DbConfig         `json:"database"` // Database connection configuration
}

// SwaggerConfig holds configuration for Swagger documentation generation and serving.
type SwaggerConfig struct {
	FilePath string `json:"file_path"` // Path to the Swagger JSON file
}

// ServerConfig holds HTTP server-related configuration including timeouts and connection limits.
type ServerConfig struct {
	Host           string        `json:"host"`            // Server host address
	Port           int           `json:"port"`            // Server port number
	ReadTimeout    time.Duration `json:"read_timeout"`    // Request read timeout
	WriteTimeout   time.Duration `json:"write_timeout"`   // Response write timeout
	IdleTimeout    time.Duration `json:"idle_timeout"`    // Connection idle timeout
	MaxConnections int           `json:"max_connections"` // Maximum number of concurrent connections
}

// DbConfig holds database connection configuration including credentials and connection details.
type DbConfig struct {
	Host         string `json:"host"`     // Database host address
	Port         int    `json:"port"`     // Database port number
	User         string `json:"user"`     // Database username
	Password     string `json:"password"` // Database password
	DatabaseName string `json:"dbname"`   // Database name
	Type         string `json:"type"`     // Database type (e.g., "mongodb", "postgres")
}

// NewConfig creates a new Config instance from a configuration manager.
// It populates all configuration fields from the provided config manager
// and sets version information from constants.
func NewConfig(configManager *config.Manager) (*Config, error) {
	cfg := Config{
		Env: configManager.GetString("env"),
		Swagger: &SwaggerConfig{
			FilePath: configManager.GetString("swagger.file_path"),
		},
		Server: &ServerConfig{
			Host:           configManager.GetString("server.host"),
			Port:           configManager.GetInt("server.port"),
			ReadTimeout:    configManager.GetDuration("server.read_timeout"),
			WriteTimeout:   configManager.GetDuration("server.write_timeout"),
			IdleTimeout:    configManager.GetDuration("server.idle_timeout"),
			MaxConnections: configManager.GetInt("server.max_connections"),
		},
		Logger: configManager.GetLogConfig(),
		Database: &DbConfig{
			Host:         configManager.GetString("database.host"),
			Port:         configManager.GetInt("database.port"),
			User:         configManager.GetString("database.user"),
			Password:     configManager.GetString("database.password"),
			DatabaseName: configManager.GetString("database.dbname"),
			Type:         configManager.GetString("database.type"),
		},
	}
	cfg.Logger.Version = constants.Version
	cfg.Logger.Commit = constants.CommitHash
	return &cfg, nil
}
