package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"library/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigManager(t *testing.T) {
	configPath := "/tmp/test-config.json"
	manager := NewConfigManager(configPath)

	assert.NotNil(t, manager)
	assert.Equal(t, configPath, manager.configPath)
	assert.NotNil(t, manager.config)
}

func TestManager_Load_JSON_FileExists(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create a test JSON config file
	jsonContent := `{
		"env": "test",
		"server": {
			"port": 8080,
			"host": "localhost"
		},
		"logging": {
			"level": "DEBUG",
			"output_to_file": true
		}
	}`

	err := os.WriteFile(configPath, []byte(jsonContent), 0600)
	require.NoError(t, err)

	manager := NewConfigManager(configPath)
	err = manager.Load()
	require.NoError(t, err)

	// Test that values were loaded correctly
	assert.Equal(t, "test", manager.GetString("env"))
	assert.Equal(t, 8080, manager.GetInt("server.port"))
	assert.Equal(t, "localhost", manager.GetString("server.host"))
	assert.Equal(t, "DEBUG", manager.GetString("logging.level"))
	assert.True(t, manager.GetBool("logging.output_to_file"))
}

func TestManager_Load_YAML_FileExists(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a test YAML config file
	yamlContent := `
env: test
server:
  port: 8080
  host: localhost
logging:
  level: DEBUG
  output_to_file: true
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0600)
	require.NoError(t, err)

	manager := NewConfigManager(configPath)
	err = manager.Load()
	require.NoError(t, err)

	// Test that values were loaded correctly
	assert.Equal(t, "test", manager.GetString("env"))
	assert.Equal(t, 8080, manager.GetInt("server.port"))
	assert.Equal(t, "localhost", manager.GetString("server.host"))
	assert.Equal(t, "DEBUG", manager.GetString("logging.level"))
	assert.True(t, manager.GetBool("logging.output_to_file"))
}

func TestManager_Load_FileNotExists_CreatesDefault(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := NewConfigManager(configPath)
	err := manager.Load()
	require.NoError(t, err)

	// Test that default values were created
	assert.Equal(t, "local", manager.GetString("env"))
	assert.Equal(t, 8080, manager.GetInt("server.port"))
	assert.False(t, manager.GetBool("logging.output_to_file"))
	assert.True(t, manager.GetBool("logging.output_to_stdio"))
	assert.Equal(t, "INFO", manager.GetString("logging.level"))

	// Verify file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err)
}

func TestManager_Load_UnsupportedFormat(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.txt")

	manager := NewConfigManager(configPath)
	err := manager.Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported config file format")
}

func TestManager_Load_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create invalid JSON
	jsonContent := `{
		"env": "test",
		"server": {
			"port": 8080,
			"host": "localhost"
		}
		"logging": {
			"level": "DEBUG"
		}
	}`

	err := os.WriteFile(configPath, []byte(jsonContent), 0600)
	require.NoError(t, err)

	manager := NewConfigManager(configPath)
	err = manager.Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON config")
}

func TestManager_Load_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create invalid YAML
	yamlContent := `
env: test
server:
  port: 8080
  host: localhost
logging:
  level: DEBUG
  output_to_file: true
  invalid: [unclosed: bracket
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0600)
	require.NoError(t, err)

	manager := NewConfigManager(configPath)
	err = manager.Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse YAML config")
}

func TestManager_Get(t *testing.T) {
	manager := NewConfigManager("/tmp/test.json")
	manager.config = map[string]interface{}{
		"string_val": "test",
		"int_val":    42,
		"bool_val":   true,
		"nested": map[string]interface{}{
			"key": "value",
		},
	}

	tests := []struct {
		name     string
		key      string
		expected interface{}
	}{
		{"Simple key", "string_val", "test"},
		{"Integer key", "int_val", 42},
		{"Boolean key", "bool_val", true},
		{"Nested key", "nested.key", "value"},
		{"Non-existent key", "nonexistent", nil},
		{"Non-existent nested key", "nested.nonexistent", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.Get(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_GetString(t *testing.T) {
	manager := NewConfigManager("/tmp/test.json")
	manager.config = map[string]interface{}{
		"string_val": "test",
		"int_val":    42,
		"bool_val":   true,
		"nil_val":    nil,
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"String value", "string_val", "test"},
		{"Integer value", "int_val", "42"},
		{"Boolean value", "bool_val", "true"},
		{"Nil value", "nil_val", ""},
		{"Non-existent key", "nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.GetString(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_GetInt(t *testing.T) {
	manager := NewConfigManager("/tmp/test.json")
	manager.config = map[string]interface{}{
		"int_val":        42,
		"float_val":      42.5,
		"string_val":     "42",
		"string_invalid": "not_a_number",
		"nil_val":        nil,
	}

	tests := []struct {
		name     string
		key      string
		expected int
	}{
		{"Integer value", "int_val", 42},
		{"Float value", "float_val", 42},
		{"String number", "string_val", 42},
		{"Invalid string", "string_invalid", 0},
		{"Nil value", "nil_val", 0},
		{"Non-existent key", "nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.GetInt(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_GetBool(t *testing.T) {
	manager := NewConfigManager("/tmp/test.json")
	manager.config = map[string]interface{}{
		"bool_true":    true,
		"bool_false":   false,
		"string_true":  "true",
		"string_false": "false",
		"string_True":  "True",
		"int_one":      1,
		"int_zero":     0,
		"nil_val":      nil,
	}

	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"Boolean true", "bool_true", true},
		{"Boolean false", "bool_false", false},
		{"String true", "string_true", true},
		{"String false", "string_false", false},
		{"String True", "string_True", true},
		{"Integer one", "int_one", true},
		{"Integer zero", "int_zero", false},
		{"Nil value", "nil_val", false},
		{"Non-existent key", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.GetBool(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_GetDuration(t *testing.T) {
	manager := NewConfigManager("/tmp/test.json")
	manager.config = map[string]interface{}{
		"duration_5s":      "5s",
		"duration_1m":      "1m",
		"duration_1h":      "1h",
		"invalid_duration": "invalid",
		"nil_val":          nil,
	}

	tests := []struct {
		name     string
		key      string
		expected time.Duration
	}{
		{"5 seconds", "duration_5s", 5 * time.Second},
		{"1 minute", "duration_1m", 1 * time.Minute},
		{"1 hour", "duration_1h", 1 * time.Hour},
		{"Invalid duration", "invalid_duration", 0},
		{"Nil value", "nil_val", 0},
		{"Non-existent key", "nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.GetDuration(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_GetLogConfig(t *testing.T) {
	manager := NewConfigManager("/tmp/test.json")
	manager.config = map[string]interface{}{
		"logging": map[string]interface{}{
			"output_to_file":    true,
			"output_to_stdio":   false,
			"log_file_path":     "/tmp/test.log",
			"log_dir":           "/logs",
			"file_writer_type":  "simple",
			"level":             "DEBUG",
			"include_timestamp": true,
			"include_level":     true,
			"include_caller":    false,
			"use_colors":        true,
			"timestamp_format":  "2006-01-02 15:04:05",
			"log_format":        "[{timestamp}] {message}",
			"buffer_size":       1024,
			"flush_interval":    "5s",
		},
	}

	logConfig := manager.GetLogConfig()
	assert.NotNil(t, logConfig)
	assert.True(t, logConfig.OutputToFile)
	assert.False(t, logConfig.OutputToStdio)
	assert.Equal(t, "/tmp/test.log", logConfig.LogFilePath)
	assert.Equal(t, "/logs", logConfig.LogDir)
	assert.Equal(t, "simple", logConfig.FileWriterType)
	assert.Equal(t, logger.DEBUG, logConfig.Level)
	assert.True(t, logConfig.IncludeTimestamp)
	assert.True(t, logConfig.IncludeLevel)
	assert.False(t, logConfig.IncludeCaller)
	assert.True(t, logConfig.UseColors)
	assert.Equal(t, "2006-01-02 15:04:05", logConfig.TimestampFormat)
	assert.Equal(t, "[{timestamp}] {message}", logConfig.LogFormat)
	assert.Equal(t, 1024, logConfig.BufferSize)
	assert.Equal(t, 5*time.Second, logConfig.FlushInterval)
}

func TestManager_ParseLogLevel(t *testing.T) {
	manager := NewConfigManager("/tmp/test.json")

	tests := []struct {
		name     string
		level    string
		expected logger.LogLevel
	}{
		{"DEBUG", "DEBUG", logger.DEBUG},
		{"debug", "debug", logger.DEBUG},
		{"INFO", "INFO", logger.INFO},
		{"info", "info", logger.INFO},
		{"WARN", "WARN", logger.WARN},
		{"warn", "warn", logger.WARN},
		{"ERROR", "ERROR", logger.ERROR},
		{"error", "error", logger.ERROR},
		{"FATAL", "FATAL", logger.FATAL},
		{"fatal", "fatal", logger.FATAL},
		{"Unknown", "UNKNOWN", logger.INFO},
		{"Empty", "", logger.INFO},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.parseLogLevel(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertMapStringInterface(t *testing.T) {
	// Test with map[interface{}]interface{} (YAML style)
	input := map[string]interface{}{
		"key1": map[interface{}]interface{}{
			"nested_key1": "value1",
			"nested_key2": 42,
		},
		"key2": map[string]interface{}{
			"nested_key3": "value3",
		},
	}

	convertMapStringInterface(input)

	// Verify conversion
	nested1, ok := input["key1"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value1", nested1["nested_key1"])
	assert.Equal(t, 42, nested1["nested_key2"])

	nested2, ok := input["key2"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value3", nested2["nested_key3"])
}

func TestMergeConfig(t *testing.T) {
	dst := map[string]interface{}{
		"existing_key": "existing_value",
		"nested": map[string]interface{}{
			"existing_nested": "existing_nested_value",
		},
	}

	src := map[string]interface{}{
		"new_key":      "new_value",
		"existing_key": "overwritten_value",
		"nested": map[string]interface{}{
			"new_nested":      "new_nested_value",
			"existing_nested": "overwritten_nested_value",
		},
	}

	mergeConfig(dst, src)

	// Verify merging behavior
	assert.Equal(t, "overwritten_value", dst["existing_key"])
	assert.Equal(t, "new_value", dst["new_key"])

	nested, ok := dst["nested"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "overwritten_nested_value", nested["existing_nested"])
	assert.Equal(t, "new_nested_value", nested["new_nested"])
}
