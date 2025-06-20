package routes

import (
	handlersMock "orderfoodonline/internal/http/handlers/mocks"
	middlewaresMock "orderfoodonline/internal/http/middlewares/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_validateDependencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockAuthMiddleware := middlewaresMock.NewMockAuthMiddleware(ctrl)
	mockProductHandler := handlersMock.NewMockProductHandler(ctrl)
	mocksSwaggerHandler := handlersMock.NewMockSwaggerHandler(ctrl)

	tests := []struct {
		name        string
		args        Dependencies
		wantErr     bool
		expectedErr string
	}{
		{
			name: "All dependencies are provided",
			args: Dependencies{
				AuthMiddleware: mockAuthMiddleware,
				ProductHandler: mockProductHandler,
				SwaggerHandler: mocksSwaggerHandler,
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "AuthMiddleware is nil",
			args: Dependencies{
				AuthMiddleware: nil,
				ProductHandler: mockProductHandler,
				SwaggerHandler: mocksSwaggerHandler,
			},
			wantErr:     true,
			expectedErr: "authMiddleware cannot be nil",
		},
		{
			name: "ProductMiddleware is nil",
			args: Dependencies{
				AuthMiddleware: mockAuthMiddleware,
				SwaggerHandler: mocksSwaggerHandler,
			},
			wantErr:     true,
			expectedErr: "productHandler cannot be nil",
		},
		{
			name: "SwaggerHandler is nil",
			args: Dependencies{
				AuthMiddleware: mockAuthMiddleware,
				ProductHandler: mockProductHandler,
			},
			wantErr:     true,
			expectedErr: "swaggerHandler cannot be nil",
		},
		// Add more test cases for each nil dependency as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDependencies(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// func TestApiRoutes(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// Create mocks for the handlers and middleware
// 	mockAuthMiddleware := mock.NewMockAuthMiddleware(ctrl)
// 	mockAuthHandler := mocks.NewMockAuthHandler(ctrl)
// 	mockLevelHandler := mocks.NewMockLevelHandler(ctrl)
// 	mockOrganizationHandler := mocks.NewMockOrganizationHandler(ctrl)
// 	mockProductHandler := mocks.NewMockProductHandler(ctrl)
// 	mockUserHandler := mocks.NewMockUserHandler(ctrl)
// 	mockRoleHandler := mocks.NewMockRoleHandler(ctrl)
// 	mockPackagingHandular := mocks.NewMockPackagingHandler(ctrl)

// 	// Create a dependencies object with the mocks
// 	deps := Dependencies{
// 		AuthMiddleware:      mockAuthMiddleware,
// 		AuthHandler:         mockAuthHandler,
// 		LevelHandler:        mockLevelHandler,
// 		OrganizationHandler: mockOrganizationHandler,
// 		ProductHandler:      mockProductHandler,
// 		UserHandler:         mockUserHandler,
// 		RoleHandler:         mockRoleHandler,
// 		PackagingHandler:    mockPackagingHandular,
// 	}

// 	router := gin.Default()

// 	// Mock the authentication middleware
// 	mockAuthMiddleware.EXPECT().Authenticate().Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()

// 	// Mock the authorization middleware with expected roles
// 	mockAuthMiddleware.EXPECT().Authorize(models.SUPERADMIN, models.ADMIN, models.MAINTAINENCE).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()

// 	mockAuthMiddleware.EXPECT().Authorize(models.ADMIN, models.MAINTAINENCE, models.PRINTERS).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()

// 	mockAuthMiddleware.EXPECT().Authorize(models.SUPERADMIN, models.ADMIN, models.MAINTAINENCE, models.BOOKEEPER).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()
// 	mockAuthMiddleware.EXPECT().Authorize(models.SUPERADMIN, models.ADMIN, models.MAINTAINENCE, models.CUSTOM).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()
// 	mockAuthMiddleware.EXPECT().Authorize(models.SUPERADMIN, models.ADMIN).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()
// 	mockAuthMiddleware.EXPECT().Authorize(models.SUPERADMIN, models.MAINTAINENCE).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()

// 	mockAuthMiddleware.EXPECT().Authorize(models.SUPERADMIN, models.ADMIN, models.PRINTERS).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()
// 	mockAuthMiddleware.EXPECT().Authorize(models.ADMIN, models.MAINTAINENCE).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()
// 	mockAuthMiddleware.EXPECT().Authorize(models.ADMIN, models.MAINTAINENCE, models.PRINTERS).Return(func(c *gin.Context) {
// 		c.Next()
// 	}).AnyTimes()
// 	err := ApiRoutes(router, deps)
// 	assert.NoError(t, err)

// 	t.Run("Test SignIn Endpoint", func(t *testing.T) {
// 		// Set up expected behavior for the mock
// 		mockAuthHandler.EXPECT().SignIn(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"message": "sign in successful"})
// 		})

// 		// Create a request to the sign-in endpoint
// 		req, _ := http.NewRequest("POST", "/api/signin", nil)
// 		resp := httptest.NewRecorder()

// 		// Perform the request
// 		router.ServeHTTP(resp, req)

// 		// Assert the response
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "sign in successful")
// 	})
// 	t.Run("Test SignIn Endpoint when invalid input is passed", func(t *testing.T) {
// 		mockAuthHandler.EXPECT().SignIn(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		})
// 		reqBody := strings.NewReader(`{"email": "user@example.com", "password": "password"}`)
// 		req, _ := http.NewRequest("POST", "/api/signin", reqBody)
// 		req.Header.Set("Content-Type", "application/json")

// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "Invalid request")
// 	})
// 	t.Run("Test SignIn Endpoint for Internalserver error", func(t *testing.T) {
// 		mockAuthHandler.EXPECT().SignIn(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
// 		})
// 		reqBody := strings.NewReader(`{"email": "user@example.com", "password": "password"}`)
// 		req, _ := http.NewRequest("POST", "/api/signin", reqBody)
// 		req.Header.Set("Content-Type", "application/json")

// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	// Test cases for public endpoints
// 	t.Run("Test SignOut Endpoint", func(t *testing.T) {
// 		mockAuthHandler.EXPECT().SignOut(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"message": "sign out successful"})
// 		})

// 		req, _ := http.NewRequest("POST", "/api/signout", nil)
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "sign out successful")
// 	})
// 	t.Run("Test SignUp Endpoint", func(t *testing.T) {
// 		mockUserHandler.EXPECT().CreateUser(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"user_id": primitive.NewObjectID().Hex()})
// 		})

// 		createUserReq := `{
// 			"name":"john",
// 			"email":"john@gmail.com",
// 			"password":"John@123"
// 			}`

// 		req, _ := http.NewRequest("POST", "/api/signup", strings.NewReader(createUserReq))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "user_id")
// 	})
// 	t.Run("Test SignUp Endpoint for Badrequest", func(t *testing.T) {
// 		mockUserHandler.EXPECT().CreateUser(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
// 		})

// 		createUserReq := `{
// 			"name":"",
// 			"email":"john@gmail.com",
// 			"password":"John@123"
// 			}`

// 		req, _ := http.NewRequest("POST", "/api/signup", strings.NewReader(createUserReq))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "Bad Request")
// 	})

// 	t.Run("Test SignUp Endpoint for Internal server error", func(t *testing.T) {
// 		mockUserHandler.EXPECT().CreateUser(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 		})

// 		createUserReq := `{
// 			"name":"",
// 			"email":"john@gmail.com",
// 			"password":"John@123"
// 			}`

// 		req, _ := http.NewRequest("POST", "/api/signup", strings.NewReader(createUserReq))
// 		req.Header.Set("Content-Type", "application/json")

// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	// Test cases for Organization endpoints
// 	t.Run("Test CreateOrganization Endpoint", func(t *testing.T) {
// 		mockOrganizationHandler.EXPECT().CreateOrganization(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"org_id": primitive.NewObjectID().Hex()})
// 		})
// 		createOrgReq := `{
// 			"name":"purplease",
// 			"admin":{
// 			"account_id":"1",
// 				  "role":1
// 				}
// 		}`
// 		req, _ := http.NewRequest("POST", "/api/organization", strings.NewReader(createOrgReq))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "org_id")
// 	})
// 	t.Run("Test CreateOrganization Endpoint for Bad Request", func(t *testing.T) {
// 		mockOrganizationHandler.EXPECT().CreateOrganization(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
// 		})
// 		createOrgReq := `{
// 			"name":"purplease",
// 			"admin":{
// 			"account_id":"1",
// 				  "role":1
// 				}
// 		}`
// 		req, _ := http.NewRequest("POST", "/api/organization", strings.NewReader(createOrgReq))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test CreateOrganization Endpoint forInternal Server Error", func(t *testing.T) {
// 		mockOrganizationHandler.EXPECT().CreateOrganization(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		})
// 		createOrgReq := `{
// 			"name":"purplease",
// 			"admin":{
// 			"account_id":"1",
// 				  "role":1
// 				}
// 		}`
// 		req, _ := http.NewRequest("POST", "/api/organization", strings.NewReader(createOrgReq))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	// Test cases for Products endpoints
// 	t.Run("Test CreateProduct Endpoint for Bad Request", func(t *testing.T) {
// 		mockProductHandler.EXPECT().CreateProduct(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
// 		})

// 		productReq := `{"name":"Product-1"}`
// 		req, _ := http.NewRequest("POST", "/api/invalid_id/products", strings.NewReader(productReq))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test CreateProduct Endpoint for nil body", func(t *testing.T) {
// 		mockProductHandler.EXPECT().CreateProduct(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
// 		})
// 		req, _ := http.NewRequest("POST", "/api/1/products", nil)
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test CreateProduct Endpoint for Internal server error", func(t *testing.T) {
// 		mockProductHandler.EXPECT().CreateProduct(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
// 		})
// 		req, _ := http.NewRequest("POST", "/api/1/products", nil)
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test CreateProduct Endpoint", func(t *testing.T) {
// 		mockProductHandler.EXPECT().CreateProduct(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"product_id": primitive.NewObjectID().Hex()})
// 		})

// 		productReq := `{"name":"Product-1"}`
// 		req, _ := http.NewRequest("POST", "/api/1/products", strings.NewReader(productReq))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "product_id")
// 	})
// 	t.Run("Test CreateProductQR Endpoint", func(t *testing.T) {
// 		mockProductHandler.EXPECT().CreateProductQR(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"product_ids": "product QR created"})
// 		})
// 		productQr := `{
// 			"name":"product Qr code",
// 			"parent_id":"67066bb1f53f98eb8ca171c1",
// 			"count":5
// 		}`
// 		req, _ := http.NewRequest("POST", "/api/1/product-qr", strings.NewReader(productQr))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "product_ids")
// 	})
// 	t.Run("Test CreateProductQR Endpoint for Bad request", func(t *testing.T) {
// 		mockProductHandler.EXPECT().CreateProductQR(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
// 		})
// 		productQr := `{
// 			"name":"product Qr code",
// 			"parent_id":"",
// 			"count":5
// 		}`
// 		req, _ := http.NewRequest("POST", "/api/1/product-qr", strings.NewReader(productQr))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test CreateProductQR Endpoint for nil body", func(t *testing.T) {
// 		mockProductHandler.EXPECT().CreateProductQR(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
// 		})
// 		req, _ := http.NewRequest("POST", "/api/1/product-qr", nil)
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	t.Run("Test CreateProductQR Endpoint for internal error", func(t *testing.T) {
// 		mockProductHandler.EXPECT().CreateProductQR(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
// 		})
// 		productQr := `{
// 			"name":"product Qr code",
// 			"parent_id":"",
// 			"count":5
// 		}`
// 		req, _ := http.NewRequest("POST", "/api/1/product-qr", strings.NewReader(productQr))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	t.Run("Test GetProductsByOrgID Endpoint", func(t *testing.T) {
// 		mockProductHandler.EXPECT().GetProductsByOrgID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"products": "products retrieved"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/1/products", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "products retrieved")
// 	})
// 	t.Run("Test GetProductsByOrgID Endpoint for Bad request", func(t *testing.T) {
// 		mockProductHandler.EXPECT().GetProductsByOrgID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid IDs"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/invalid/products", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test GetProductsByOrgID Endpoint for internal error", func(t *testing.T) {
// 		mockProductHandler.EXPECT().GetProductsByOrgID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/invalid/products", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	t.Run("Test GetProductByID Endpoint", func(t *testing.T) {
// 		mockProductHandler.EXPECT().GetProductByID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"products": "product retrieved"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/1/products/1", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "product retrieved")
// 	})
// 	t.Run("Test GetProductByID Endpoint for Bad Request", func(t *testing.T) {
// 		mockProductHandler.EXPECT().GetProductByID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/1/products/invalidID", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	t.Run("Test GetProductByID Endpoint for Internal error ", func(t *testing.T) {
// 		mockProductHandler.EXPECT().GetProductByID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/1/products/1", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	// Test cases for Levels endpoints

// 	t.Run("Test CreateLevel Endpoint", func(t *testing.T) {
// 		mockLevelHandler.EXPECT().CreateLevel(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"level_id": primitive.NewObjectID().Hex()})
// 		})
// 		levelReq := `{

//         "name":"Level-1",
//         "parent_id":"1",
// 		"organization_id":"2"
// 			}`
// 		req, _ := http.NewRequest("POST", "/api/1/level", strings.NewReader(levelReq))
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "level_id")
// 	})
// 	t.Run("Test CreateLevel Endpoint for Bad Request", func(t *testing.T) {
// 		mockLevelHandler.EXPECT().CreateLevel(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
// 		})
// 		levelReq := `{

//         "name":"Level-1",
//         "parent_id":"1",
// 		"organization_id":"2"
// 			}`
// 		req, _ := http.NewRequest("POST", "/api/1/level", strings.NewReader(levelReq))
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test CreateLevel Endpoint for nil body", func(t *testing.T) {
// 		mockLevelHandler.EXPECT().CreateLevel(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
// 		})

// 		req, _ := http.NewRequest("POST", "/api/1/level", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test CreateLevel Endpoint for Internal error", func(t *testing.T) {
// 		mockLevelHandler.EXPECT().CreateLevel(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
// 		})
// 		levelReq := `{

//         "name":"Level-1",
//         "parent_id":"1",
// 		"organization_id":"2"
// 			}`
// 		req, _ := http.NewRequest("POST", "/api/1/level", strings.NewReader(levelReq))
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test GetLevelsByParentID Endpoint", func(t *testing.T) {
// 		mockLevelHandler.EXPECT().GetLevelsByParentID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"levels": "levels retrived"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/1/levels/parent/2", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "levels retrived")
// 	})
// 	t.Run("Test GetLevelsByParentID Endpoint for Bad Request", func(t *testing.T) {
// 		mockLevelHandler.EXPECT().GetLevelsByParentID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent ID"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/1/levels/parent/InvalidID", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test GetLevelsByParentID Endpoint for Intrnal error", func(t *testing.T) {
// 		mockLevelHandler.EXPECT().GetLevelsByParentID(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
// 		})

// 		req, _ := http.NewRequest("GET", "/api/1/levels/parent/InvalidID", nil)
// 		resp := httptest.NewRecorder()
// 		req.Header.Set("Content-Type", "application/json")
// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	// Test cases for roles EndPoint

// 	t.Run("Test AssignRoleToUser Endpoint", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().AssignRoleToUser(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"status": "Role assigned successfully"})
// 		})
// 		assignRoleReq := `{
//          "user_id": "1",
//          "role_id": "2",
//          "role": "0"
// 			}`
// 		req, _ := http.NewRequest("POST", "/api/1/roles/assign", strings.NewReader(assignRoleReq))
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "status")
// 	})
// 	t.Run("Test AssignRoleToUser Endpoint for Bad Request", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().AssignRoleToUser(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
// 		})
// 		assignRoleReq := `{
//          "organization_id": "1",
//          "level_id": "2",
//          "parent_id": "3"
// 			}`
// 		req, _ := http.NewRequest("POST", "/api/1/roles/assign", strings.NewReader(assignRoleReq))
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test AssignRoleToUser Endpoint for nil body  Request", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().AssignRoleToUser(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
// 		})

// 		req, _ := http.NewRequest("POST", "/api/1/roles/assign", nil)
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test AssignRoleToUser Endpoint for Internal error", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().AssignRoleToUser(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "INternal error"})
// 		})
// 		assignRoleReq := `{
//          "organization_id": "1",
//          "level_id": "2",
//          "parent_id": "3"
// 			}`
// 		req, _ := http.NewRequest("POST", "/api/1/roles/assign", strings.NewReader(assignRoleReq))
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test GetUsersByOrganization Endpoint", func(t *testing.T) {
// 		mockUserHandler.EXPECT().GetUsersByOrganization(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"users": " users retrived by orgID"})
// 		})
// 		req, _ := http.NewRequest("GET", "/api/users/organization/1", nil)
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "users")
// 	})

// 	t.Run("Test GetUsersByOrganization Endpoint for Bad Request", func(t *testing.T) {
// 		mockUserHandler.EXPECT().GetUsersByOrganization(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
// 		})
// 		req, _ := http.NewRequest("GET", "/api/users/organization/InvalidID", nil)
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	t.Run("Test GetUsersByOrganization Endpoint for Internal error", func(t *testing.T) {
// 		mockUserHandler.EXPECT().GetUsersByOrganization(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
// 		})
// 		req, _ := http.NewRequest("GET", "/api/users/organization/InvalidID", nil)
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	t.Run("Test RemoveUserRole Endpoint", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().RemoveUserRole(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"message": "User role removed successfully"})
// 		})
// 		req, err := http.NewRequest("DELETE", "/api/1/roles/2/admin", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "message")
// 	})
// 	t.Run("Test RemoveUserRole Endpoint for Bad request", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().RemoveUserRole(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
// 		})
// 		req, err := http.NewRequest("DELETE", "/api/1/roles/InvalidID/adnin", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test RemoveUserRole Endpoint for internal error", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().RemoveUserRole(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		})
// 		req, err := http.NewRequest("DELETE", "/api/1/roles/1/admin", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test GetUsersWithRole Endpoint", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().GetUsersWithRole(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"user_ids": "userids"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/1/users/superadmin/1", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "user_ids")
// 	})
// 	t.Run("Test GetUsersWithRole Endpoint for internal error", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().GetUsersWithRole(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/1/users/superadmin/1", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test GetUsersWithRole Endpoint for Bad Request", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().GetUsersWithRole(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/1/users/superadmin/''", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test RequestRoleAccess Endpoint", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().RequestRoleAccess(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"message": "Role access request submitted successfully"})
// 		})
// 		request := `{"role":"some_role"}`
// 		req, err := http.NewRequest("POST", "/api/1/roles/request", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "message")
// 	})
// 	t.Run("Test RequestRoleAccess Endpoint for Bad request", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().RequestRoleAccess(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
// 		})
// 		request := `{"role":"InvalidReq"}`
// 		req, err := http.NewRequest("POST", "/api/1/roles/request", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})
// 	t.Run("Test RequestRoleAccess Endpoint for Internal error", func(t *testing.T) {
// 		mockRoleHandler.EXPECT().RequestRoleAccess(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		})
// 		request := `{"role":"some_role"}`
// 		req, err := http.NewRequest("POST", "/api/1/roles/request", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "error")
// 	})

// 	t.Run("Test CreatePackagingInfo Endpoint", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().CreatePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusCreated, gin.H{"packaging_info": "createdInfo"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("POST", "/api/1/packaging", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusCreated, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "createdInfo")
// 	})
// 	t.Run("Test CreatePackagingInfo Endpoint for BadRequest", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().CreatePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "badrequest"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("POST", "/api/''/packaging", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "badrequest")
// 	})
// 	t.Run("Test CreatePackagingInfo Endpoint for internal error", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().CreatePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("POST", "/api/''/packaging", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "internal")
// 	})
// 	t.Run("Test GetPackagingInfo Endpoint ", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().GetPackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"packaging_info": "createdInfo"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/1/packaging/1", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "createdInfo")
// 	})
// 	t.Run("Test GetPackagingInfo Endpoint for badrequest ", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().GetPackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "badrequest"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/''/packaging/1", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "badrequest")
// 	})
// 	t.Run("Test GetPackagingInfo Endpoint for badrequest ", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().GetPackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/1/packaging/1", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "internal error")
// 	})

// 	t.Run("Test UpdatePackagingInfo Endpoint", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().UpdatePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"message": "Packaging info updated successfully"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("PUT", "/api/1/packaging", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "Packaging info updated successfully")
// 	})
// 	t.Run("Test UpdatePackagingInfo Endpoint for bad request", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().UpdatePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "badrequest"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("PUT", "/api/1/packaging", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "badrequest")
// 	})
// 	t.Run("Test UpdatePackagingInfo Endpoint for internal error", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().UpdatePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("PUT", "/api/1/packaging", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "internal error")
// 	})
// 	t.Run("Test DeletePackagingInfo Endpoint ", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().DeletePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"message": "Packaging info deleted successfully"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("DELETE", "/api/1/packaging/1", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "Packaging info deleted successfully")
// 	})
// 	t.Run("Test DeletePackagingInfo Endpoint for bad request ", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().DeletePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("DELETE", "/api/''/packaging/1", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "Invalid ID")
// 	})

// 	t.Run("Test DeletePackagingInfo Endpoint for internal error", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().DeletePackagingInfo(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "internal error"})
// 		})
// 		request := `{"_id":"1","name":"demo","count":"1"}`
// 		req, err := http.NewRequest("DELETE", "/api/1/packaging/1", strings.NewReader(request))
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "internal error")
// 	})
// 	t.Run("Test GetPackagingForProduct Endpoint ", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().GetPackagingForProduct(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"packaging_info": "packagingInfo"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/1/packaging/product/1", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "packagingInfo")
// 	})
// 	t.Run("Test GetPackagingForProduct Endpoint for badrequest", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().GetPackagingForProduct(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusBadRequest, gin.H{"errro": "badrequest"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/1/packaging/product/1", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "badrequest")
// 	})
// 	t.Run("Test GetPackagingForProduct Endpoint for internal error", func(t *testing.T) {
// 		mockPackagingHandular.EXPECT().GetPackagingForProduct(gomock.Any()).Times(1).Do(func(c *gin.Context) {
// 			c.JSON(http.StatusInternalServerError, gin.H{"errro": "internal error"})
// 		})
// 		req, err := http.NewRequest("GET", "/api/1/packaging/product/1", nil)
// 		if err != nil {
// 			t.Fatal(err.Error())
// 			return
// 		}
// 		resp := httptest.NewRecorder()
// 		router.ServeHTTP(resp, req)
// 		assert.Equal(t, http.StatusInternalServerError, resp.Code)
// 		assert.Contains(t, resp.Body.String(), "internal error")
// 	})

// }

// func TestHandleOptionsRequests(t *testing.T) {
// 	router := gin.Default()
// 	router.Use(handleOptionsRequests())

// 	router.GET("/test", func(c *gin.Context) {
// 		c.String(http.StatusOK, "test")
// 	})

// 	req, _ := http.NewRequest(http.MethodOptions, "/test", nil)

// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNoContent, w.Code)

// 	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
// 	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
// 	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
// }

// func TestHandleOptionsRequestsNonOptions(t *testing.T) {
// 	router := gin.Default()
// 	router.Use(handleOptionsRequests())

// 	router.GET("/test", func(c *gin.Context) {
// 		c.String(http.StatusOK, "test")
// 	})

// 	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	assert.Equal(t, "test", w.Body.String())
// }
