package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileWriter handles simple log file writing without rotation.
type FileWriter struct {
	filename string
	file     *os.File
	mu       sync.Mutex
}

// NewFileWriter creates a new FileWriter for the specified filename.
func NewFileWriter(filename string) (*FileWriter, error) {
	writer := &FileWriter{
		filename: filename,
	}

	if err := writer.openFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

// openFile opens the log file
func (fw *FileWriter) openFile() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(fw.filename)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open file in append mode
	file, err := os.OpenFile(fw.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	fw.file = file
	return nil
}

// Write writes data to the log file.
func (fw *FileWriter) Write(data []byte) (int, error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.file == nil {
		return 0, fmt.Errorf("file is not open")
	}

	return fw.file.Write(data)
}

// Flush flushes the file buffer to disk.
func (fw *FileWriter) Flush() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.file != nil {
		return fw.file.Sync()
	}
	return nil
}

// Close closes the FileWriter and the underlying file.
func (fw *FileWriter) Close() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.file != nil {
		return fw.file.Close()
	}
	return nil
}
