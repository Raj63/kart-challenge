package logger

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected string
	}{
		{"DEBUG level", DEBUG, "DEBUG"},
		{"INFO level", INFO, "INFO"},
		{"WARN level", WARN, "WARN"},
		{"ERROR level", ERROR, "ERROR"},
		{"FATAL level", FATAL, "FATAL"},
		{"Unknown level", LogLevel(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.level.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogLevel_Color(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected string
	}{
		{"DEBUG level", DEBUG, "\033[36m"},
		{"INFO level", INFO, "\033[32m"},
		{"WARN level", WARN, "\033[33m"},
		{"ERROR level", ERROR, "\033[31m"},
		{"FATAL level", FATAL, "\033[35m"},
		{"Unknown level", LogLevel(99), "\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.level.Color()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultLogConfig(t *testing.T) {
	config := DefaultLogConfig()

	assert.NotNil(t, config)
	assert.False(t, config.OutputToFile)
	assert.True(t, config.OutputToStdio)
	assert.Equal(t, "/logs", config.LogDir)
	assert.Equal(t, "simple", config.FileWriterType)
	assert.Equal(t, INFO, config.Level)
	assert.True(t, config.IncludeTimestamp)
	assert.True(t, config.IncludeLevel)
	assert.False(t, config.IncludeCaller)
	assert.True(t, config.UseColors)
	assert.Equal(t, "2006-01-02 15:04:05.000", config.TimestampFormat)
	assert.Equal(t, "[{timestamp}] [{level}] [{version}-{commit}] {message}", config.LogFormat)
	assert.Equal(t, 1024, config.BufferSize)
	assert.Equal(t, 5*time.Second, config.FlushInterval)
}

func TestNewLogger_StdioOnly(t *testing.T) {
	config := &LogConfig{
		OutputToStdio: true,
		OutputToFile:  false,
		Level:         INFO,
	}

	logger, err := NewLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	assert.NotNil(t, logger)
	assert.Equal(t, config, logger.config)
	assert.NotNil(t, logger.stdioWriter)
	assert.Nil(t, logger.fileWriter)
	assert.Nil(t, logger.buffer)
}

func TestNewLogger_WithFile(t *testing.T) {
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)
	logFile := filepath.Join(tempDir, "test.log")

	config := &LogConfig{
		OutputToStdio:  true,
		OutputToFile:   true,
		LogDir:         tempDir,
		LogFilePath:    logFile,
		FileWriterType: "simple",
		Level:          INFO,
	}

	logger, err := NewLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	assert.NotNil(t, logger)
	assert.NotNil(t, logger.fileWriter)

	// Verify file was created
	_, err = os.Stat(logFile)
	assert.NoError(t, err)
}

func TestNewLogger_InvalidFileWriterType(t *testing.T) {
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	config := &LogConfig{
		OutputToStdio:  true,
		OutputToFile:   true,
		LogDir:         tempDir,
		FileWriterType: "invalid",
		Level:          INFO,
	}

	logger, err := NewLogger(config)
	assert.Error(t, err)
	assert.Nil(t, logger)
	assert.Contains(t, err.Error(), "unknown file writer type")
}

func TestLogger_LogMethods(t *testing.T) {
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	// Prepare log config
	config := &LogConfig{
		OutputToStdio:    false,
		OutputToFile:     true,
		LogDir:           tempDir,
		LogFilePath:      filepath.Join(tempDir, "test.log"),
		FileWriterType:   "simple",
		Level:            DEBUG,
		IncludeTimestamp: true,
		IncludeLevel:     true,
		UseColors:        false,
		LogFormat:        "[{timestamp}] [{level}] [{version}-{commit}] {message}",
	}

	logger, err := NewLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	// Log all levels except Fatal
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	// Do not call logger.Fatal("fatal message")

	// Open file in append mode
	file, err := os.OpenFile(config.LogFilePath, os.O_RDONLY, 0600)
	require.NoError(t, err)
	defer file.Close()

	buf, err := io.ReadAll(file)
	require.NoError(t, err)
	output := string(buf)

	assert.Contains(t, output, "debug message")
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
	// Optionally assert that fatal is not present, since it's not logged
	assert.NotContains(t, output, "fatal message")
}

func TestLogger_LevelFiltering(t *testing.T) {
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	// Prepare log config
	config := &LogConfig{
		OutputToStdio:    false,
		OutputToFile:     true,
		LogDir:           tempDir,
		LogFilePath:      filepath.Join(tempDir, "test.log"),
		FileWriterType:   "simple",
		Level:            WARN,
		IncludeTimestamp: true,
		IncludeLevel:     true,
		UseColors:        false,
		LogFormat:        "[{timestamp}] [{level}] [{version}-{commit}] {message}",
	}

	logger, err := NewLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	// Test level filtering
	logger.Debug("debug message") // Should not appear
	logger.Info("info message")   // Should not appear
	logger.Warn("warn message")   // Should appear
	logger.Error("error message") // Should appear

	// Open file in append mode
	file, err := os.OpenFile(config.LogFilePath, os.O_RDONLY, 0600)
	require.NoError(t, err)
	defer file.Close()

	buf, err := io.ReadAll(file)
	require.NoError(t, err)
	output := string(buf)

	assert.NotContains(t, output, "debug message")
	assert.NotContains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
	assert.NotContains(t, output, "fatal message")
}

func TestLogger_FormatMessage(t *testing.T) {
	config := &LogConfig{
		OutputToStdio:    true,
		OutputToFile:     false,
		Level:            INFO,
		IncludeTimestamp: true,
		IncludeLevel:     true,
		IncludeCaller:    true,
		UseColors:        false,
		TimestampFormat:  "2006-01-02 15:04:05",
		LogFormat:        "[{timestamp}] [{level}] {message}",
		Version:          "1.0.0",
		Commit:           "abc123",
	}

	logger, err := NewLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	entry := logEntry{
		level:   INFO,
		message: "test message",
		time:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		caller:  "test.go:10",
	}

	formatted := logger.formatMessage(entry)
	assert.Contains(t, formatted, "2023-01-01 12:00:00")
	assert.Contains(t, formatted, "INFO")
	assert.Contains(t, formatted, "test message")
}

func TestLogger_SetAndGetLevel(t *testing.T) {
	config := &LogConfig{
		OutputToStdio: true,
		OutputToFile:  false,
		Level:         INFO,
	}

	logger, err := NewLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	// Test initial level
	assert.Equal(t, INFO, logger.GetLevel())

	// Test setting new level
	logger.SetLevel(ERROR)
	assert.Equal(t, ERROR, logger.GetLevel())

	// Test IsLevelEnabled
	assert.False(t, logger.IsLevelEnabled(DEBUG))
	assert.False(t, logger.IsLevelEnabled(INFO))
	assert.False(t, logger.IsLevelEnabled(WARN))
	assert.True(t, logger.IsLevelEnabled(ERROR))
	assert.True(t, logger.IsLevelEnabled(FATAL))
}

func TestLogger_Close(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	config := &LogConfig{
		OutputToStdio:  true,
		OutputToFile:   true,
		LogDir:         tempDir,
		LogFilePath:    logFile,
		FileWriterType: "simple",
		BufferSize:     10,
		Level:          INFO,
	}

	logger, err := NewLogger(config)
	require.NoError(t, err)

	// Test close
	err = logger.Close()
	assert.NoError(t, err)
}

func TestReplacePlaceholder(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		placeholder string
		value       string
		expected    string
	}{
		{
			name:        "Replace timestamp",
			format:      "[{timestamp}] {message}",
			placeholder: "{timestamp}",
			value:       "2023-01-01",
			expected:    "[2023-01-01] {message}",
		},
		{
			name:        "Replace level",
			format:      "[{level}] {message}",
			placeholder: "{level}",
			value:       "INFO",
			expected:    "[INFO] {message}",
		},
		{
			name:        "Replace message",
			format:      "[{level}] {message}",
			placeholder: "{message}",
			value:       "test message",
			expected:    "[{level}] test message",
		},
		{
			name:        "No placeholder",
			format:      "simple message",
			placeholder: "{timestamp}",
			value:       "2023-01-01",
			expected:    "simple message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replacePlaceholder(tt.format, tt.placeholder, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCaller(t *testing.T) {
	caller := getCaller()
	assert.NotEmpty(t, caller)
	assert.Contains(t, caller, ".go:")
}

func TestLogger_Flush(t *testing.T) {
	config := &LogConfig{
		OutputToStdio: true,
		OutputToFile:  false,
		BufferSize:    5,
		Level:         INFO,
	}

	logger, err := NewLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	// Test flush
	logger.Info("message to flush")
	logger.flush()
}
