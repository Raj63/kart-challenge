// Package middlewares provides HTTP middleware components for the Order Food Online service.
package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOptionsHandler_OptionsRequest(t *testing.T) {
	// Given: Options middleware
	optionsHandler := OptionsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create OPTIONS request
	req, _ := http.NewRequest("OPTIONS", "/api/test", nil)
	c.Request = req

	// When: Processing the OPTIONS request through the options middleware
	optionsHandler(c)

	// Then: Request should be aborted with No Content status and CORS headers
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Check CORS headers
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestOptionsHandler_NonOptionsRequest(t *testing.T) {
	// Given: Options middleware
	optionsHandler := OptionsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create non-OPTIONS request
	req, _ := http.NewRequest("GET", "/api/test", nil)
	c.Request = req

	// Add a next handler
	nextCalled := false
	nextHandler := func(c *gin.Context) {
		nextCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the non-OPTIONS request through the options middleware
	optionsHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Next handler should be called and response should be successful
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestOptionsHandler_DifferentHTTPMethods(t *testing.T) {
	// Given: Options middleware
	optionsHandler := OptionsHandler()

	testCases := []struct {
		method      string
		path        string
		shouldAbort bool
		statusCode  int
	}{
		{"OPTIONS", "/api/test", true, http.StatusNoContent},
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
			nextCalled := false
			nextHandler := func(c *gin.Context) {
				nextCalled = true
				c.Status(tc.statusCode)
			}

			// When: Processing the request through the options middleware
			optionsHandler(c)
			if !c.IsAborted() {
				nextHandler(c)
			}

			// Then: Request should be processed correctly
			if tc.shouldAbort {
				assert.True(t, c.IsAborted())
				assert.Equal(t, tc.statusCode, w.Code)
				assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.True(t, nextCalled)
				assert.Equal(t, tc.statusCode, w.Code)
			}
		})
	}
}

func TestOptionsHandler_OptionsRequestWithHeaders(t *testing.T) {
	// Given: Options middleware
	optionsHandler := OptionsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create OPTIONS request with additional headers
	req, _ := http.NewRequest("OPTIONS", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
	c.Request = req

	// When: Processing the OPTIONS request through the options middleware
	optionsHandler(c)

	// Then: Request should be aborted with No Content status and CORS headers
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Check CORS headers
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestOptionsHandler_OptionsRequestRootPath(t *testing.T) {
	// Given: Options middleware
	optionsHandler := OptionsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create OPTIONS request to root path
	req, _ := http.NewRequest("OPTIONS", "/", nil)
	c.Request = req

	// When: Processing the OPTIONS request through the options middleware
	optionsHandler(c)

	// Then: Request should be aborted with No Content status and CORS headers
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Check CORS headers
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestOptionsHandler_OptionsRequestWithQueryParams(t *testing.T) {
	// Given: Options middleware
	optionsHandler := OptionsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create OPTIONS request with query parameters
	req, _ := http.NewRequest("OPTIONS", "/api/test?param=value", nil)
	c.Request = req

	// When: Processing the OPTIONS request through the options middleware
	optionsHandler(c)

	// Then: Request should be aborted with No Content status and CORS headers
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Check CORS headers
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestOptionsHandler_NonOptionsRequestWithNextHandler(t *testing.T) {
	// Given: Options middleware
	optionsHandler := OptionsHandler()

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create GET request
	req, _ := http.NewRequest("GET", "/api/test", nil)
	c.Request = req

	// Add a next handler
	nextCalled := false
	nextHandler := func(c *gin.Context) {
		nextCalled = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// When: Processing the request through the options middleware
	optionsHandler(c)
	if !c.IsAborted() {
		nextHandler(c)
	}

	// Then: Next handler should be called and response should be successful
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")

	// CORS headers should not be set for non-OPTIONS requests
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
