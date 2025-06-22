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
	mocksMetricsHandler := middlewaresMock.NewMockMetricsMiddleware(ctrl)

	tests := []struct {
		name        string
		args        Dependencies
		wantErr     bool
		expectedErr string
	}{
		{
			name: "All dependencies are provided",
			args: Dependencies{
				AuthMiddleware:    mockAuthMiddleware,
				ProductHandler:    mockProductHandler,
				SwaggerHandler:    mocksSwaggerHandler,
				MetricsMiddleware: mocksMetricsHandler,
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "AuthMiddleware is nil",
			args: Dependencies{
				AuthMiddleware:    nil,
				ProductHandler:    mockProductHandler,
				SwaggerHandler:    mocksSwaggerHandler,
				MetricsMiddleware: mocksMetricsHandler,
			},
			wantErr:     true,
			expectedErr: "authMiddleware cannot be nil",
		},
		{
			name: "ProductMiddleware is nil",
			args: Dependencies{
				AuthMiddleware:    mockAuthMiddleware,
				SwaggerHandler:    mocksSwaggerHandler,
				MetricsMiddleware: mocksMetricsHandler,
			},
			wantErr:     true,
			expectedErr: "productHandler cannot be nil",
		},
		{
			name: "SwaggerHandler is nil",
			args: Dependencies{
				AuthMiddleware:    mockAuthMiddleware,
				ProductHandler:    mockProductHandler,
				MetricsMiddleware: mocksMetricsHandler,
			},
			wantErr:     true,
			expectedErr: "swaggerHandler cannot be nil",
		},
		{
			name: "MetricsMiddleware is nil",
			args: Dependencies{
				AuthMiddleware: mockAuthMiddleware,
				ProductHandler: mockProductHandler,
				SwaggerHandler: mocksSwaggerHandler,
			},
			wantErr:     true,
			expectedErr: "metricsMiddleware cannot be nil",
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
