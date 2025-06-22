package logger

import (
	"io"
)

// ILogger defines the interface for logging operations.
// It provides a contract for all logging functionality including
// different log levels, configuration management, and resource cleanup.
type ILogger interface {
	// Debug logs a message at the DEBUG level.
	// Use this for verbose output useful for debugging during development.
	// The message is formatted using fmt.Sprintf with the provided format and arguments.
	Debug(format string, args ...interface{})

	// Info logs a message at the INFO level.
	// Use this for general application events or high-level progress reporting.
	// The message is formatted using fmt.Sprintf with the provided format and arguments.
	Info(format string, args ...interface{})

	// Warn logs a message at the WARN level.
	// Use this when something unexpected happens, but the app can recover or continue.
	// The message is formatted using fmt.Sprintf with the provided format and arguments.
	Warn(format string, args ...interface{})

	// Error logs a message at the ERROR level.
	// Use this for serious issues that should be investigated but don't require immediate shutdown.
	// The message is formatted using fmt.Sprintf with the provided format and arguments.
	Error(format string, args ...interface{})

	// Fatal logs a message at the FATAL level and then exits the application with status code 1.
	// Use this when a non-recoverable error occurs that requires the app to terminate.
	// The message is formatted using fmt.Sprintf with the provided format and arguments.
	Fatal(format string, args ...interface{})

	// Close closes the logger and flushes any remaining data.
	// This method should be called when the logger is no longer needed to ensure
	// all buffered data is written and resources are properly cleaned up.
	// Returns an error if the file writer fails to close properly.
	Close() error

	// SetLevel sets the logging level for the logger.
	// Only messages at or above this level will be logged.
	// The level parameter should be one of the predefined LogLevel constants.
	SetLevel(level LogLevel)

	// GetLevel returns the current logging level of the logger.
	// This can be used to check what level is currently set or to
	// conditionally log based on the current level.
	GetLevel() LogLevel

	// IsLevelEnabled checks if a specific log level is enabled.
	// Returns true if the provided level is at or above the current logging level,
	// false otherwise. This is useful for avoiding expensive operations when
	// a log level is disabled.
	IsLevelEnabled(level LogLevel) bool

	// GetFileWriter returns the current file writer for advanced operations.
	// This method provides access to the underlying io.WriteCloser for custom
	// file operations. Use with caution as direct manipulation of the file writer
	// may interfere with the logger's internal state.
	GetFileWriter() io.WriteCloser
}
