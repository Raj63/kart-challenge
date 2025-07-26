package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileWriter(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)
	defer writer.Close()

	assert.NotNil(t, writer)
	assert.Equal(t, filename, writer.filename)
	assert.NotNil(t, writer.file)

	// Verify file was created
	_, err = os.Stat(filename)
	assert.NoError(t, err)
}

func TestNewFileWriter_CreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	nestedDir := filepath.Join(tempDir, "nested", "deep")
	filename := filepath.Join(nestedDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)
	defer writer.Close()

	assert.NotNil(t, writer)

	// Verify directory was created
	_, err = os.Stat(nestedDir)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(filename)
	assert.NoError(t, err)
}

func TestNewFileWriter_InvalidPath(t *testing.T) {
	// Try to create file in a non-existent directory that we can't create
	invalidPath := "/invalid/path/that/should/not/exist/test.log"

	writer, err := NewFileWriter(invalidPath)
	assert.Error(t, err)
	assert.Nil(t, writer)
	assert.Contains(t, err.Error(), "failed to create log directory")
}

func TestFileWriter_Write(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)
	defer writer.Close()

	// Test writing data
	testData := []byte("test log message\n")
	n, err := writer.Write(testData)
	assert.NoError(t, err)
	assert.Equal(t, len(testData), n)

	// Verify data was written to file
	content, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, testData, content)
}

func TestFileWriter_WriteMultiple(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)
	defer writer.Close()

	// Write multiple messages
	messages := []string{
		"first message\n",
		"second message\n",
		"third message\n",
	}

	for _, msg := range messages {
		n, err := writer.Write([]byte(msg))
		assert.NoError(t, err)
		assert.Equal(t, len(msg), n)
	}

	// Verify all messages were written
	content, err := os.ReadFile(filename)
	assert.NoError(t, err)
	expectedContent := "first message\nsecond message\nthird message\n"
	assert.Equal(t, expectedContent, string(content))
}

func TestFileWriter_WriteAfterClose(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)

	// Close the writer
	err = writer.Close()
	assert.NoError(t, err)

	// Try to write after close
	testData := []byte("test message\n")
	_, err = writer.Write(testData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file already closed")
}

func TestFileWriter_Flush(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)
	defer writer.Close()

	// Write some data
	testData := []byte("test message\n")
	_, err = writer.Write(testData)
	assert.NoError(t, err)

	// Flush the data
	err = writer.Flush()
	assert.NoError(t, err)

	// Verify data was flushed to disk
	content, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, testData, content)
}

func TestFileWriter_FlushAfterClose(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)

	// Close the writer
	err = writer.Close()
	assert.NoError(t, err)

	// Try to flush after close
	err = writer.Flush()
	assert.Error(t, err) // Flush should error when file is nil
}

func TestFileWriter_Close(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)

	// Write some data
	testData := []byte("test message\n")
	_, err = writer.Write(testData)
	assert.NoError(t, err)

	// Close the writer
	err = writer.Close()
	assert.NoError(t, err)

	// Verify file was closed by checking if we can't write anymore
	_, err = writer.Write([]byte("another message\n"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file already closed")
}

func TestFileWriter_CloseMultiple(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)

	// Close multiple times should not error
	err = writer.Close()
	assert.NoError(t, err)

	err = writer.Close()
	assert.Error(t, err)
}

func TestFileWriter_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewFileWriter(filename)
	require.NoError(t, err)
	defer writer.Close()

	// Check file permissions
	info, err := os.Stat(filename)
	assert.NoError(t, err)

	// File should have 0600 permissions (owner read/write only)
	mode := info.Mode()
	assert.Equal(t, os.FileMode(0600), mode&os.ModePerm)
}
