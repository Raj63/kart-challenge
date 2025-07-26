// Package middlewares provides HTTP middleware components for the Order Food Online service.
package middlewares

import (
	"library/logger/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewAuthMiddleware(t *testing.T) {
	// Given: A logger mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	// When: Creating a new auth middleware
	authMiddleware := NewAuthMiddleware(mockLogger)

	// Then: It should not be nil and implement the interface
	assert.NotNil(t, authMiddleware)
}

func TestAuthMiddleware_Authenticate_ValidAPIKey(t *testing.T) {
	// Given: A valid API key and auth middleware
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	authMiddleware := NewAuthMiddleware(mockLogger)
	handler := authMiddleware.Authenticate()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with valid API key
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("api_key", "apitest")
	c.Request = req

	// When: Calling the authenticate middleware with valid API key
	handler(c)

	// Then: Request should proceed (not aborted)
	assert.False(t, c.IsAborted())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_Authenticate_MissingAPIKey(t *testing.T) {
	// Given: No API key and auth middleware
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	authMiddleware := NewAuthMiddleware(mockLogger)
	handler := authMiddleware.Authenticate()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request without API key
	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	// When: Calling the authenticate middleware without API key
	handler(c)

	// Then: Request should be aborted with unauthorized status
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Verify response body contains error message
	assert.Contains(t, w.Body.String(), "API key required")
}

func TestAuthMiddleware_Authenticate_InvalidAPIKey(t *testing.T) {
	// Given: Invalid API key and auth middleware
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	authMiddleware := NewAuthMiddleware(mockLogger)
	handler := authMiddleware.Authenticate()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with invalid API key
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("api_key", "invalid_key")
	c.Request = req

	// When: Calling the authenticate middleware with invalid API key
	handler(c)

	// Then: Request should be aborted with unauthorized status
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Verify response body contains error message
	assert.Contains(t, w.Body.String(), "Invalid API key")
}

func TestAuthMiddleware_Authenticate_EmptyAPIKey(t *testing.T) {
	// Given: Empty API key and auth middleware
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	authMiddleware := NewAuthMiddleware(mockLogger)
	handler := authMiddleware.Authenticate()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with empty API key
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("api_key", "")
	c.Request = req

	// When: Calling the authenticate middleware with empty API key
	handler(c)

	// Then: Request should be aborted with unauthorized status
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Verify response body contains error message
	assert.Contains(t, w.Body.String(), "API key required")
}

func TestAuthMiddleware_Authenticate_WhitespaceAPIKey(t *testing.T) {
	// Given: Whitespace-only API key and auth middleware
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	authMiddleware := NewAuthMiddleware(mockLogger)
	handler := authMiddleware.Authenticate()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with whitespace-only API key
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("api_key", "   ")
	c.Request = req

	// When: Calling the authenticate middleware with whitespace-only API key
	handler(c)

	// Then: Request should be aborted with unauthorized status
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Verify response body contains error message
	assert.Contains(t, w.Body.String(), "API key required")
}

func TestAuthMiddleware_Authorize(t *testing.T) {
	// Given: Auth middleware
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	authMiddleware := NewAuthMiddleware(mockLogger)
	handler := authMiddleware.Authorize()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	// When: Calling the authorize middleware
	handler(c)

	// Then: Request should proceed (not aborted) since authorize is currently a no-op
	assert.False(t, c.IsAborted())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_Authenticate_WithNextHandler(t *testing.T) {
	// Given: Auth middleware and a next handler
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	authMiddleware := NewAuthMiddleware(mockLogger)
	authHandler := authMiddleware.Authenticate()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with valid API key
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("api_key", "apitest")
	c.Request = req

	// Add a next handler to verify it gets called
	nextCalled := false
	nextHandler := func(c *gin.Context) {
		nextCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the middleware chain
	authHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Next handler should be called and response should be successful
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestAuthMiddleware_Authenticate_AbortedWithNextHandler(t *testing.T) {
	// Given: Auth middleware and a next handler
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	authMiddleware := NewAuthMiddleware(mockLogger)
	authHandler := authMiddleware.Authenticate()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request without API key
	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	// Add a next handler to verify it does NOT get called
	nextCalled := false
	nextHandler := func(c *gin.Context) {
		nextCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the middleware chain
	authHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Next handler should NOT be called due to authentication failure
	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "API key required")
}
