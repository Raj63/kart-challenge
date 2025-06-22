// Package middlewares provides HTTP middleware components for the Order Food Online service.
package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricsMiddleware(t *testing.T) {
	// When: Creating a new metrics middleware
	metricsMiddleware := NewMetricsMiddleware()

	// Then: It should not be nil and implement the interface
	assert.NotNil(t, metricsMiddleware)
}

func TestMetricsMiddleware_RecordMetrics_Success(t *testing.T) {
	// Given: Metrics middleware
	metricsMiddleware := NewMetricsMiddleware()
	handler := metricsMiddleware.RecordMetrics()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/test", nil)
	c.Request = req

	// Add a next handler to simulate processing
	nextCalled := false
	nextHandler := func(c *gin.Context) {
		nextCalled = true
		time.Sleep(10 * time.Millisecond) // Simulate some processing time
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the metrics middleware
	handler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Next handler should be called and response should be successful
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestMetricsMiddleware_RecordMetrics_WithFullPath(t *testing.T) {
	// Given: Metrics middleware with a route that has a full path
	metricsMiddleware := NewMetricsMiddleware()
	handler := metricsMiddleware.RecordMetrics()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("POST", "/api/product/123", nil)
	c.Request = req
	// Note: FullPath is set by Gin router, not directly in tests

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"id": "123"})
	}

	// When: Processing the request through the metrics middleware
	handler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be processed successfully
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}

func TestMetricsMiddleware_RecordMetrics_WithQueryParameters(t *testing.T) {
	// Given: Metrics middleware with query parameters
	metricsMiddleware := NewMetricsMiddleware()
	handler := metricsMiddleware.RecordMetrics()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with query parameters
	req, _ := http.NewRequest("GET", "/api/product?category=food&limit=10", nil)
	c.Request = req
	// Note: FullPath is set by Gin router, not directly in tests

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"products": []string{}})
	}

	// When: Processing the request through the metrics middleware
	handler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be processed successfully
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "products")
}

func TestMetricsMiddleware_RecordMetrics_ErrorResponse(t *testing.T) {
	// Given: Metrics middleware that will result in an error
	metricsMiddleware := NewMetricsMiddleware()
	handler := metricsMiddleware.RecordMetrics()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/nonexistent", nil)
	c.Request = req

	// Add a next handler that returns an error
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	}

	// When: Processing the request through the metrics middleware
	handler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Error response should be returned
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Not found")
}

func TestMetricsMiddleware_RecordMetrics_DifferentHTTPMethods(t *testing.T) {
	// Given: Metrics middleware
	metricsMiddleware := NewMetricsMiddleware()
	handler := metricsMiddleware.RecordMetrics()

	testCases := []struct {
		method     string
		path       string
		statusCode int
	}{
		{"GET", "/api/health", http.StatusOK},
		{"POST", "/api/order", http.StatusOK},
		{"PUT", "/api/product/123", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.method+"_"+tc.path, func(t *testing.T) {
			// Set up Gin context
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create test request
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			c.Request = req

			// Add a next handler
			nextHandler := func(c *gin.Context) {
				c.Status(tc.statusCode)
			}

			// When: Processing the request through the metrics middleware
			handler(c)
			if !c.IsAborted() {
				nextHandler(c)
			}

			// Then: Request should be processed with correct status code
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestMetricsMiddleware_RecordMetrics_AbortedRequest(t *testing.T) {
	// Given: Metrics middleware
	metricsMiddleware := NewMetricsMiddleware()
	handler := metricsMiddleware.RecordMetrics()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/test", nil)
	c.Request = req

	// Add a next handler that aborts
	nextHandler := func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	}

	// When: Processing the request through the metrics middleware
	handler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be aborted with forbidden status
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMetricsMiddleware_RecordMetrics_EmptyPath(t *testing.T) {
	// Given: Metrics middleware with empty path
	metricsMiddleware := NewMetricsMiddleware()
	handler := metricsMiddleware.RecordMetrics()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with empty path
	req, _ := http.NewRequest("GET", "/", nil)
	c.Request = req
	// Note: FullPath is set by Gin router, not directly in tests

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "root"})
	}

	// When: Processing the request through the metrics middleware
	handler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be processed successfully using URL path
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "root")
}
