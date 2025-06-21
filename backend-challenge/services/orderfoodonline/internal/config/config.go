package config

import (
	"library/config"
	"library/logger"
	"orderfoodonline/internal/constants"
	"time"
)

// Config holds the application configuration.
type Config struct {
	Env      string            `json:"env"`
	Swagger  *SwaggerConfig    `json:"swagger"`
	Server   *ServerConfig     `json:"server"`
	Logger   *logger.LogConfig `json:"logger"`
	Database *DbConfig         `json:"database"`
}

// SwaggerConfig holds configuration for Swagger documentation.
type SwaggerConfig struct {
	FilePath string `json:"file_path"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
	IdleTimeout    time.Duration `json:"idle_timeout"`
	MaxConnections int           `json:"max_connections"`
}

// DbConfig holds database related configuration.
type DbConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DatabaseName string `json:"dbname"`
	Type         string `json:"type"`
}

// NewConfig creates a new Config from a config manager.
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
