package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"library/logger"

	"gopkg.in/yaml.v3"
)

// Manager handles configuration loading and management
type Manager struct {
	configPath string
	config     map[string]interface{}
	watchers   []Watcher
	mu         sync.RWMutex
}

// Watcher is called when configuration changes
type Watcher func(key string, oldValue, newValue interface{})

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *Manager {
	return &Manager{
		configPath: configPath,
		config:     make(map[string]interface{}),
		watchers:   make([]Watcher, 0),
	}
}

// Load loads configuration from file
func (cm *Manager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Create default config
		if err := cm.createDefaultConfig(); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
		log.Printf("Created and using default configuration file: %s", cm.configPath)
		return nil
	}
	// Generate default config and load
	cm.config = cm.generateDefaultConfig()

	// Load config based on file extension
	ext := strings.ToLower(filepath.Ext(cm.configPath))
	switch ext {
	case ".json":
		return cm.loadJSON()
	case ".yaml", ".yml":
		return cm.loadYAML()
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
}

// loadJSON loads JSON configuration
func (cm *Manager) loadJSON() error {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

	mergeConfig(cm.config, config)
	log.Printf("Loaded configuration from: %s", cm.configPath)
	return nil
}

// loadYAML loads YAML configuration

func (cm *Manager) loadYAML() error {
	file, err := os.Open(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close config file: %v", err)
		}
	}()

	var config map[string]interface{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("failed to parse YAML config: %w", err)
	}

	convertMapStringInterface(config) // normalize nested maps

	mergeConfig(cm.config, config)
	log.Printf("Loaded YAML configuration from: %s", cm.configPath)
	return nil
}

// convertMapStringInterface  normalize these before merging since YAML lib parses things where
// nested maps may come out as map[interface{}]interface{}, not map[string]interface{}
func convertMapStringInterface(m map[string]interface{}) {
	for k, v := range m {
		switch child := v.(type) {
		case map[interface{}]interface{}:
			// Convert child to map[string]interface{}
			mapped := make(map[string]interface{})
			for key, val := range child {
				strKey := fmt.Sprintf("%v", key)
				mapped[strKey] = val
			}
			convertMapStringInterface(mapped) // recurse
			m[k] = mapped
		case map[string]interface{}:
			convertMapStringInterface(child)
		}
	}
}

// createDefaultConfig creates a default configuration file
func (cm *Manager) createDefaultConfig() error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Generate default config
	cm.config = cm.generateDefaultConfig()

	// Write config based on file extension
	ext := strings.ToLower(filepath.Ext(cm.configPath))
	switch ext {
	case ".json":
		return cm.writeJSONConfig(cm.config)
	case ".yaml", ".yml":
		return cm.writeYAMLConfig(cm.config)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
}

// generateDefaultConfig generates default configuration
func (cm *Manager) generateDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"env": "local",
		"logging": map[string]interface{}{
			"output_to_file":    false,
			"output_to_stdio":   true,
			"log_dir":           "/logs",
			"file_writer_type":  "simple", // "none", "simple"
			"level":             "INFO",
			"include_timestamp": true,
			"include_level":     true,
			"include_caller":    false,
			"use_colors":        true,
			"timestamp_format":  "2006-01-02 15:04:05.000",
			"log_format":        "[{timestamp}] [{level}] {message}",
			"buffer_size":       1024,
			"async_logging":     false,
			"flush_interval":    "5s",
		},
		"server": map[string]interface{}{
			"host":            "",
			"port":            8080,
			"read_timeout":    "30s",
			"write_timeout":   "30s",
			"idle_timeout":    "60s",
			"max_connections": 10000,
		},
		"swagger": map[string]interface{}{
			"file_path": ".",
		},
	}
}

// mergeConfig will merge new configuration data into an existing map[string]interface{}, such that it:
// - Preserves existing defaults,
// - Adds new keys, and
// - Overrides existing keys with new values
func mergeConfig(dst, src map[string]interface{}) {
	for key, val := range src {
		// If value is a nested map, merge recursively
		if srcMap, ok := val.(map[string]interface{}); ok {
			if dstMap, exists := dst[key].(map[string]interface{}); exists {
				mergeConfig(dstMap, srcMap)
				continue
			}
		}
		// Otherwise, overwrite or set the value
		dst[key] = val
	}
}

// writeJSONConfig writes JSON configuration
func (cm *Manager) writeJSONConfig(config map[string]interface{}) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON config: %w", err)
	}

	return os.WriteFile(cm.configPath, data, 0600)
}

// writeYAMLConfig writes YAML configuration
func (cm *Manager) writeYAMLConfig(config map[string]interface{}) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML config: %w", err)
	}

	return os.WriteFile(cm.configPath, data, 0600)
}

// Get retrieves a configuration value
func (cm *Manager) Get(key string) interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	keys := strings.Split(key, ".")
	return cm.getNestedValue(cm.config, keys)
}

// GetString retrieves a string configuration value
func (cm *Manager) GetString(key string) string {
	value := cm.Get(key)
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// GetInt retrieves an integer configuration value
func (cm *Manager) GetInt(key string) int {
	value := cm.Get(key)
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// GetInt64 retrieves an int64 configuration value
func (cm *Manager) GetInt64(key string) int64 {
	value := cm.Get(key)
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

// GetBool retrieves a boolean configuration value
func (cm *Manager) GetBool(key string) bool {
	value := cm.Get(key)
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true"
	case int:
		return v != 0
	}
	return false
}

// GetDuration retrieves a duration configuration value
func (cm *Manager) GetDuration(key string) time.Duration {
	value := cm.GetString(key)
	if value == "" {
		return 0
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Invalid duration format for key %s: %s", key, value)
		return 0
	}

	return duration
}

// GetStringSlice retrieves a string slice configuration value
func (cm *Manager) GetStringSlice(key string) []string {
	value := cm.Get(key)
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	}
	return nil
}

// Set sets a configuration value
func (cm *Manager) Set(key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	oldValue := cm.Get(key)
	keys := strings.Split(key, ".")
	cm.setNestedValue(cm.config, keys, value)

	// Notify watchers
	for _, watcher := range cm.watchers {
		watcher(key, oldValue, value)
	}
}

// Save saves the current configuration to file
func (cm *Manager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	ext := strings.ToLower(filepath.Ext(cm.configPath))
	switch ext {
	case ".json":
		return cm.writeJSONConfig(cm.config)
	case ".yaml", ".yml":
		return cm.writeYAMLConfig(cm.config)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
}

// AddWatcher adds a configuration change watcher
func (cm *Manager) AddWatcher(watcher Watcher) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.watchers = append(cm.watchers, watcher)
}

// getNestedValue retrieves a nested value from a map
func (cm *Manager) getNestedValue(data map[string]interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		return data
	}

	key := keys[0]
	value, exists := data[key]
	if !exists {
		return nil
	}

	if len(keys) == 1 {
		return value
	}

	if nested, ok := value.(map[string]interface{}); ok {
		return cm.getNestedValue(nested, keys[1:])
	}

	return nil
}

// setNestedValue sets a nested value in a map
func (cm *Manager) setNestedValue(data map[string]interface{}, keys []string, value interface{}) {
	if len(keys) == 0 {
		return
	}

	key := keys[0]
	if len(keys) == 1 {
		data[key] = value
		return
	}

	// Create nested map if it doesn't exist
	if _, exists := data[key]; !exists {
		data[key] = make(map[string]interface{})
	}

	if nested, ok := data[key].(map[string]interface{}); ok {
		cm.setNestedValue(nested, keys[1:], value)
	}
}

// GetLogConfig returns a logger configuration
func (cm *Manager) GetLogConfig() *logger.LogConfig {
	return &logger.LogConfig{
		OutputToFile:     cm.GetBool("logging.output_to_file"),
		OutputToStdio:    cm.GetBool("logging.output_to_stdio"),
		LogFilePath:      cm.GetString("logging.log_file_path"),
		LogDir:           cm.GetString("logging.log_dir"),
		FileWriterType:   cm.GetString("logging.file_writer_type"),
		Level:            cm.parseLogLevel(cm.GetString("logging.level")),
		IncludeTimestamp: cm.GetBool("logging.include_timestamp"),
		IncludeLevel:     cm.GetBool("logging.include_level"),
		IncludeCaller:    cm.GetBool("logging.include_caller"),
		UseColors:        cm.GetBool("logging.use_colors"),
		TimestampFormat:  cm.GetString("logging.timestamp_format"),
		LogFormat:        cm.GetString("logging.log_format"),
		BufferSize:       cm.GetInt("logging.buffer_size"),
		AsyncLogging:     cm.GetBool("logging.async_logging"),
		FlushInterval:    cm.GetDuration("logging.flush_interval"),
	}
}

// parseLogLevel parses a log level string
func (cm *Manager) parseLogLevel(level string) logger.LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return logger.DEBUG
	case "INFO":
		return logger.INFO
	case "WARN":
		return logger.WARN
	case "ERROR":
		return logger.ERROR
	case "FATAL":
		return logger.FATAL
	default:
		return logger.INFO
	}
}
