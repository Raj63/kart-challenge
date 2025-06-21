package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel defines the severity level for a log message.
type LogLevel int

const (
	// DEBUG represents detailed debug information, typically only useful for developers.
	DEBUG LogLevel = iota

	// INFO represents general operational entries about what's happening inside the application.
	INFO

	// WARN indicates a potential issue or important situation that is not necessarily an error.
	WARN

	// ERROR indicates a runtime error or unexpected condition that should be investigated.
	ERROR

	// FATAL indicates a severe issue that will likely cause the application to terminate.
	FATAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color returns the ANSI color code for the log level
func (l LogLevel) Color() string {
	switch l {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	case FATAL:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// FileWriterType defines the type of file writer used for logging output.
type FileWriterType int

const (
	// FileWriterNone disables file logging.
	FileWriterNone FileWriterType = iota

	// FileWriterSimple writes logs to a single file without rotation.
	FileWriterSimple

	// FileWriterRotating writes logs to files with rotation based on size, time, or both.
	FileWriterRotating
)

// LogConfig holds the logging configuration
type LogConfig struct {
	// App specifics Version & commit-hash
	Version string `json:"-"`
	Commit  string `json:"-"`

	// Output configuration
	OutputToFile   bool   `json:"output_to_file"`
	OutputToStdio  bool   `json:"output_to_stdio"`
	LogFilePath    string `json:"log_file_path"`
	LogDir         string `json:"log_dir"`
	FileWriterType string `json:"file_writer_type"` // "none", "simple"

	// Log level and format configuration
	Level            LogLevel `json:"level"`
	IncludeTimestamp bool     `json:"include_timestamp"`
	IncludeLevel     bool     `json:"include_level"`
	IncludeCaller    bool     `json:"include_caller"`
	UseColors        bool     `json:"use_colors"`

	// Format configuration
	TimestampFormat string `json:"timestamp_format"`
	LogFormat       string `json:"log_format"`

	// Performance configuration
	BufferSize    int           `json:"buffer_size"`
	FlushInterval time.Duration `json:"flush_interval"`
}

// DefaultLogConfig returns the default conversion options
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		// Basic Options
		OutputToFile:     false,
		OutputToStdio:    true,
		LogFilePath:      "",
		LogDir:           "/logs",
		FileWriterType:   "simple",
		Level:            INFO,
		IncludeTimestamp: true,
		IncludeLevel:     true,
		IncludeCaller:    false,
		UseColors:        true,
		TimestampFormat:  "2006-01-02 15:04:05.000",
		LogFormat:        "[{timestamp}] [{level}] [{version}-{commit}] {message}",
		BufferSize:       1024,
		FlushInterval:    5 * time.Second,
	}
}

// Logger represents a logger instance
type Logger struct {
	config      *LogConfig
	fileWriter  io.WriteCloser
	stdioWriter io.Writer
	mu          sync.Mutex
	buffer      chan logEntry
	stopChan    chan struct{}
}

// logEntry represents a log entry
type logEntry struct {
	level   LogLevel
	message string
	time    time.Time
	caller  string
}

// NewLogger creates a new logger instance
func NewLogger(config *LogConfig) (*Logger, error) {
	logger := &Logger{
		config:   config,
		stopChan: make(chan struct{}),
	}

	// Initialize stdio writer
	if config.OutputToStdio {
		logger.stdioWriter = os.Stdout
	}

	// Initialize file writer
	if config.OutputToFile {
		if err := logger.initFileWriter(); err != nil {
			return nil, fmt.Errorf("failed to initialize file writer: %w", err)
		}
	}

	return logger, nil
}

// initFileWriter initializes the appropriate file writer
func (l *Logger) initFileWriter() error {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(l.config.LogDir, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Generate log file path if not provided
	if l.config.LogFilePath == "" {
		timestamp := time.Now().Format("2006-01-02")
		l.config.LogFilePath = filepath.Join(l.config.LogDir, fmt.Sprintf("app-%s.log", timestamp))
	}

	// Create appropriate file writer based on type
	switch l.config.FileWriterType {
	case "simple":
		fileWriter, err := NewFileWriter(l.config.LogFilePath)
		if err != nil {
			return fmt.Errorf("failed to create simple file writer: %w", err)
		}
		l.fileWriter = fileWriter

	case "none":
		// No file writer
		return nil

	default:
		return fmt.Errorf("unknown file writer type: %s", l.config.FileWriterType)
	}

	return nil
}

// writeLogEntry writes a log entry to all configured outputs
func (l *Logger) writeLogEntry(entry logEntry) {
	if entry.level < l.config.Level {
		return
	}

	formattedMessage := l.formatMessage(entry)

	l.mu.Lock()
	defer l.mu.Unlock()

	// Write to stdio
	if l.config.OutputToStdio && l.stdioWriter != nil {
		if _, err := fmt.Fprintln(l.stdioWriter, formattedMessage); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to stdio: %v\n", err)
		}
	}

	// Write to file
	if l.config.OutputToFile && l.fileWriter != nil {
		if _, err := l.fileWriter.Write([]byte(formattedMessage + "\n")); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
		}
	}
}

// formatMessage formats the log message according to the configuration
func (l *Logger) formatMessage(entry logEntry) string {
	message := l.config.LogFormat

	// if version is not specified then set to default
	if l.config.Version == "" {
		l.config.Version = "0.1.0"
	}

	message = replacePlaceholder(message, "{version}", l.config.Version)
	message = replacePlaceholder(message, "{commit}", l.config.Commit)

	if l.config.IncludeTimestamp {
		timestamp := entry.time.Format(l.config.TimestampFormat)
		message = replacePlaceholder(message, "{timestamp}", timestamp)
	}

	if l.config.IncludeLevel {
		levelStr := entry.level.String()
		if l.config.UseColors {
			levelStr = entry.level.Color() + levelStr + "\033[0m"
		}
		message = replacePlaceholder(message, "{level}", levelStr)
	}

	if l.config.IncludeCaller {
		message = replacePlaceholder(message, "{caller}", entry.caller)
	}

	message = replacePlaceholder(message, "{message}", entry.message)
	return message
}

// replacePlaceholder replaces a placeholder in the format string
func replacePlaceholder(format, placeholder, value string) string {
	return strings.ReplaceAll(format, placeholder, value)
}

// Debug logs a message at the DEBUG level.
// Use this for verbose output useful for debugging during development.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs a message at the INFO level.
// Use this for general application events or high-level progress reporting.
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a message at the WARN level.
// Use this when something unexpected happens, but the app can recover or continue.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs a message at the ERROR level.
// Use this for serious issues that should be investigated but don't require immediate shutdown.
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a message at the FATAL level and then exits the application with status code 1.
// Use this when a non-recoverable error occurs that requires the app to terminate.
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	entry := logEntry{
		level:   level,
		message: message,
		time:    time.Now(),
		caller:  getCaller(),
	}

	l.writeLogEntry(entry)
}

// getCaller returns the file, line, and function name of the caller.
func getCaller() string {
	// Skip 2 levels: getCaller → the logger function → actual caller
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}

	funcName := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		funcName = fn.Name()
	}

	return fmt.Sprintf("%s:%d [%s]", filepath.Base(file), line, filepath.Base(funcName))
}

// flush flushes any buffered data
func (l *Logger) flush() {
	if l.fileWriter != nil {
		if flusher, ok := l.fileWriter.(interface{ Flush() error }); ok {
			if err := flusher.Flush(); err != nil {
				fmt.Fprintf(os.Stderr, "Error flushing file writer: %v\n", err)
			}
		}
	}
}

// Close closes the logger and flushes any remaining data
func (l *Logger) Close() error {
	l.flush()

	if l.fileWriter != nil {
		return l.fileWriter.Close()
	}

	return nil
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Level = level
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.config.Level
}

// IsLevelEnabled checks if a log level is enabled
func (l *Logger) IsLevelEnabled(level LogLevel) bool {
	return level >= l.GetLevel()
}

// GetFileWriter returns the current file writer (for advanced operations)
func (l *Logger) GetFileWriter() io.WriteCloser {
	return l.fileWriter
}
