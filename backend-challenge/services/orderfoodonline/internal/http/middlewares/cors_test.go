// Package middlewares provides HTTP middleware components for the Order Food Online service.
package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCorsHandler(t *testing.T) {
	// Given: CORS middleware
	corsHandler := CorsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/test", nil)
	c.Request = req

	// Add a next handler
	nextCalled := false
	nextHandler := func(c *gin.Context) {
		nextCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the CORS middleware
	corsHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Next handler should be called and response should be successful
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestCorsHandler_DifferentHTTPMethods(t *testing.T) {
	// Given: CORS middleware
	corsHandler := CorsHandler()

	testCases := []struct {
		method     string
		path       string
		statusCode int
	}{
		{"GET", "/api/health", http.StatusOK},
		{"POST", "/api/order", http.StatusOK},
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
			req.Header.Set("Origin", "http://localhost:3000")
			c.Request = req

			// Add a next handler
			nextHandler := func(c *gin.Context) {
				c.Status(tc.statusCode)
			}

			// When: Processing the request through the CORS middleware
			corsHandler(c)
			if !c.IsAborted() {
				nextHandler(c)
			}

			// Then: Request should be processed with correct status code and CORS headers
			assert.Equal(t, tc.statusCode, w.Code)
			assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

func TestCorsHandler_CustomHeaders(t *testing.T) {
	// Given: CORS middleware
	corsHandler := CorsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request with custom headers
	req, _ := http.NewRequest("POST", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("X-Geo-Location", "US")
	req.Header.Set("X-Language", "en")
	req.Header.Set("X-Timezone", "UTC")
	c.Request = req

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the CORS middleware
	corsHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be processed successfully with CORS headers
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCorsHandler_ExposeHeaders(t *testing.T) {
	// Given: CORS middleware
	corsHandler := CorsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	c.Request = req

	// Add a next handler that sets custom headers
	nextHandler := func(c *gin.Context) {
		c.Header("Content-Length", "123")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the CORS middleware
	corsHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be processed successfully with exposed headers
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Expose-Headers"), "Content-Length")
}

func TestCorsHandler_NoOriginHeader(t *testing.T) {
	// Given: CORS middleware
	corsHandler := CorsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create test request without Origin header
	req, _ := http.NewRequest("GET", "/api/test", nil)
	c.Request = req

	// Add a next handler
	nextHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the CORS middleware
	corsHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Request should be processed successfully
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}
