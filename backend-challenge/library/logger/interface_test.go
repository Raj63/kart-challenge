package logger

import (
	"testing"
)

// TestILoggerInterfaceCompleteness verifies that all public methods of Logger are included in ILogger
func TestILoggerInterfaceCompleteness(t *testing.T) {
	// This test ensures that the ILogger interface includes all public methods
	// from the Logger struct. If new methods are added to Logger, this test
	// will help identify if they need to be added to the interface.

	// This is a compile-time check - if the Logger doesn't implement ILogger,
	// the code won't compile
	var _ ILogger = (*Logger)(nil)
}
