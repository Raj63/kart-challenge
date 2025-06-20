package config

import (
	"library/config"
	"library/logger"
	"time"
)

// Config holds the application configuration.
type Config struct {
	Env     string            `json:"env"`
	Swagger *SwaggerConfig    `json:"swagger"`
	Server  *ServerConfig     `json:"server"`
	Logger  *logger.LogConfig `json:"logger"`
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
	}
	return &cfg, nil
}
