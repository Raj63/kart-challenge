// Package middlewares provides HTTP middleware components for the Order Food Online service.
package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterHandler(t *testing.T) {
	// Given: Rate limiter middleware
	rateLimiterHandler := RateLimiterHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345" // Set client IP
	c.Request = req

	// Add a next handler
	nextCalled := false
	nextHandler := func(c *gin.Context) {
		nextCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the rate limiter middleware
	rateLimiterHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Next handler should be called and response should be successful
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestRateLimiterHandler_DifferentClientIPs(t *testing.T) {
	// Given: Rate limiter middleware
	rateLimiterHandler := RateLimiterHandler()

	testCases := []struct {
		clientIP   string
		statusCode int
	}{
		{"127.0.0.1:12345", http.StatusOK},
		{"192.168.1.1:54321", http.StatusOK},
		{"10.0.0.1:8080", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run("ClientIP_"+tc.clientIP, func(t *testing.T) {
			// Set up Gin context
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create test request
			req, _ := http.NewRequest("GET", "/api/test", nil)
			req.RemoteAddr = tc.clientIP
			c.Request = req

			// Add a next handler
			nextHandler := func(c *gin.Context) {
				c.Status(tc.statusCode)
			}

			// When: Processing the request through the rate limiter middleware
			rateLimiterHandler(c)
			if !c.IsAborted() {
				nextHandler(c)
			}

			// Then: Request should be processed with correct status code
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestRateLimiterHandler_DifferentHTTPMethods(t *testing.T) {
	// Given: Rate limiter middleware
	rateLimiterHandler := RateLimiterHandler()

	testCases := []struct {
		method     string
		path       string
		statusCode int
	}{
		{"GET", "/api/health", http.StatusOK},
		{"PUT", "/api/product/123", http.StatusOK},
		{"PATCH", "/api/product/123", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.method+"_"+tc.path, func(t *testing.T) {
			// Set up Gin context
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create test request
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			req.RemoteAddr = "127.0.0.1:12345"
			c.Request = req

			// Add a next handler
			nextHandler := func(c *gin.Context) {
				c.Status(tc.statusCode)
			}

			// When: Processing the request through the rate limiter middleware
			rateLimiterHandler(c)
			if !c.IsAborted() {
				nextHandler(c)
			}

			// Then: Request should be processed with correct status code
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestRateLimiterHandler_WithXForwardedFor(t *testing.T) {
	// Given: Rate limiter middleware
	rateLimiterHandler := RateLimiterHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with X-Forwarded-For header
	req, _ := http.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.1")
	c.Request = req

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the rate limiter middleware
	rateLimiterHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be processed successfully
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestRateLimiterHandler_WithXRealIP(t *testing.T) {
	// Given: Rate limiter middleware
	rateLimiterHandler := RateLimiterHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with X-Real-IP header
	req, _ := http.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Real-IP", "198.51.100.1")
	c.Request = req

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the rate limiter middleware
	rateLimiterHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be processed successfully
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestRateLimiterHandler_NoRemoteAddr(t *testing.T) {
	// Given: Rate limiter middleware
	rateLimiterHandler := RateLimiterHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request without RemoteAddr
	req, _ := http.NewRequest("GET", "/api/test", nil)
	// Note: RemoteAddr is empty
	c.Request = req

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the rate limiter middleware
	rateLimiterHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should still be processed successfully
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestRateLimiterHandler_ConcurrentRequests(t *testing.T) {
	// Given: Rate limiter middleware
	rateLimiterHandler := RateLimiterHandler()

	// Test concurrent requests from the same IP
	// Note: This is a basic test - in a real scenario, you'd want to test actual rate limiting
	// but that would require more complex setup with time-based testing

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	c.Request = req

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing multiple requests through the rate limiter middleware
	for i := 0; i < 5; i++ {
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		rateLimiterHandler(c)
		if !c.IsAborted() {
			nextHandler(c)
		}

		// Then: Each request should be processed successfully
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	}
}
