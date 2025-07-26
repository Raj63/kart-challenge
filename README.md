# Kart Challenge – Food Ordering Platform

A robust, modular food ordering backend and supporting library, designed for extensibility, performance, and developer productivity. This project features a RESTful API for product listing, cart management, and order processing, as well as a shared Go library for logging, configuration, and integrations. The platform includes an optimized coupon processing system capable of handling large files (1-2 GB) with parallel processing and resume functionality.


## 🚀 What's Inside

- ⚙️ Microservices with clear domain boundaries
- 🌐 REST APIs with versioning
- 🧪 Unit, integration, and contract tests
- 🐳 Docker-based local dev environment
- 🗃️ Mongodb integrations
- 📊 Observability: Logging, Tracing, Metrics
- 🔐 Secure config, secrets management, and auth hooks
- 📦 DevTools: Task runners, hot reloads, linters, CI templates
- 📁 Clean project layout (Go idiomatic)

## 🎯 Purpose

This repository is created as a reference implementation to:
- Demonstrate **best practices** for building Go microservices
- Highlight **developer productivity tools** for local and team workflows
- Provide a **foundation** for production-grade Go backend systems

---

## 🚀 Latest Features & Improvements

### **Enhanced Unit Testing & Mocking**
- **Generated Mocks**: Automated mock generation using GoMock for all interfaces
- **Interface-Based Testing**: Clean separation of concerns with interface-based repository design
- **Comprehensive Test Coverage**: Unit tests for all services with proper mocking patterns
- **Human-Readable Tests**: Well-structured test cases with clear Given/When/Then patterns
- **Mock Validation**: Proper expectation verification and error scenario testing

### **Improved Repository Architecture**
- **Interface-Driven Design**: MongoDB collections abstracted through interfaces for better testability
- **Dependency Injection**: Clean dependency management with interface-based mocking
- **Type Safety**: Strong typing with proper error handling and validation
- **Modular Structure**: Clear separation between data access and business logic

### **Performance Optimizations**
- **High-Performance File Processing**: Optimized coupon file processor capable of handling 1-2 GB files with:
  - Parallel processing with worker pools (4 concurrent workers)
  - Increased batch sizes (5000 items)
  - Optimized I/O operations with larger buffers (1MB scanner, 64KB hash buffer)
  - Memory-efficient processing with pre-allocated slices
  - **3-5x faster processing** for large files

### **Developer Experience Enhancements**
- **Comprehensive Documentation**: All exposed types and functions now have detailed GoDoc comments
- **Enhanced Logger Interface**: Complete ILogger interface with all public methods documented
- **Automated Code Quality**: Pre-commit hooks with linting, formatting, and security scanning
- **Mock Generation**: Automated mock generation for testing with GoMock
- **API Documentation**: Auto-generated Swagger/OpenAPI documentation

### **Robust Error Handling & Reliability**
- **Resume Functionality**: Coupon processing can resume from where it left off after failures
- **File Deduplication**: MD5 hash-based file tracking prevents duplicate processing
- **Graceful Error Recovery**: Comprehensive error handling with proper status updates
- **Context Cancellation**: Support for graceful shutdown and timeout handling

---

## Project Structure

```bash
kart-challenge/
├── backend-challenge/
│   ├── library/                # Shared Go library (logger, config, etc.)
│   │   ├── logger/            # Advanced logging with file rotation and async support
│   │   │   ├── mocks/         # Generated mocks for testing
│   │   │   └── interface.go   # Complete ILogger interface definition
│   │   └── config/            # Configuration management with validation
│   ├── services/
│   │   ├── coupons/            # Coupons microservice with optimized file processing
│   │   │   ├── cmd/processor/  # Coupons processor entrypoint
│   │   │   ├── internal/       # Service internals (processor, repository, config)
│   │   │   │   ├── repository/ # Interface-based repository with generated mocks
│   │   │   │   │   ├── mocks/  # Generated mocks for testing
│   │   │   │   │   └── interface.go # Collection and CouponRepository interfaces
│   │   │   │   └── processor/  # Optimized file processor with resume capability
│   │   │   └── data/          # Sample data and test files
│   │   └── orderfoodonline/    # Main food ordering service
│   │       ├── cmd/rest/       # REST API entrypoint and docs
│   │       ├── internal/       # Service internals (handlers, middlewares, routes, etc.)
│   │       │   ├── service/    # Business logic with comprehensive unit tests
│   │       │   ├── repository/ # Data access layer with generated mocks
│   │       │   └── http/       # HTTP handlers with proper error handling
│   │       ├── migrations/     # Database migrations and seeding
│   │       └── Dockerfile      # Multi-stage build for the service
│   └── Makefile                # Root Makefile for orchestration
├── api/
│   └── openapi.yaml            # OpenAPI spec for the API
├── docker-compose.local.yml    # Local dev orchestration
├── .pre-commit-config.yaml     # Pre-commit hooks for Go quality and CI
├── shell.nix                   # Nix shell for reproducible dev env
└── README.md                   # This file
```

---

## 🛠️ Developer-Friendly Features

### **Quick Setup & Development**
- **One-Command Setup**: `make start` launches the entire development environment
- **Hot Reload-TODO**: Automatic code reloading with `air` for instant feedback
- **Docker Compose**: Complete environment with MongoDB, services, and networking
- **Nix Shell**: Reproducible development environment (optional)

### **Code Quality & Testing**
- **Pre-commit Hooks**: Automated quality gates before commits
- **Comprehensive Testing**: Unit tests with coverage reporting and generated mocks
- **Mock Generation**: Automated mock creation for isolated testing with GoMock
- **Static Analysis**: Security scanning, linting, and code formatting
- **Interface Compliance**: Automated verification of interface implementations

### **Documentation & API**
- **Auto-generated Docs**: Swagger documentation from code comments
- **GoDoc Comments**: Comprehensive documentation for all public APIs
- **OpenAPI Spec**: Machine-readable API specification
- **Interactive API**: Swagger UI for testing endpoints

### **Monitoring & Debugging**
- **Structured Logging**: JSON-formatted logs with configurable levels
- **File Rotation**: Automatic log file management
- **Performance Metrics**: Built-in timing and monitoring
- **Error Tracking**: Detailed error context and stack traces
- **Prometheus Metrics**: Comprehensive metrics collection with `/metrics` endpoint
  - HTTP request counts, durations, and status codes
  - Database query latency and operation counts
  - Order processing metrics and business KPIs
  - Active database connections monitoring

---

## Features

### **Core Platform**
- **Product Listing & Cart API**: RESTful endpoints for products, cart, and order management
- **Coupon Processing**: High-performance file processing with resume capability
- **Database Migrations**: Automated schema management and data seeding
- **Authentication**: API key-based authentication middleware

### **Shared Infrastructure**
- **Advanced Logging**: Configurable logging with file rotation, colors, and async support
- **Configuration Management**: Environment-aware configuration with validation
- **Error Handling**: Centralized error management and reporting
- **Health Checks**: Built-in health and version endpoints

### **Development Tools**
- **Pre-commit Quality Gates**: Lint, format, staticcheck, security scan, and tests
- **API Documentation**: Swagger/OpenAPI 3.1 docs, auto-generated with `swag`
- **Dockerized**: Multi-stage Dockerfiles for efficient builds and minimal runtime images
- **Makefile Automation**: Common tasks for build, test, docs, and pre-commit
- **Mock Generation**: Automated test double creation for isolated testing

---

## Development Setup

### Prerequisites

- [Go 1.23+](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Make](https://www.gnu.org/software/make/)
- [Nix (optional)](https://nixos.org/) for reproducible environments
- [pre-commit](https://pre-commit.com/) (install via pip or your package manager)

### Quick Start (Recommended)

```sh
# Clone and setup
git clone <your-repo-url>
cd kart-challenge
pre-commit install

# Start everything with one command
make start

# Access the API
curl http://localhost:8080/api/health
```

### Manual Setup

```sh
# Install pre-commit hooks
pre-commit install

# Start services
make start

# Or run individual services
docker compose -f docker-compose.local.yml up --build
```

---

## Makefile Targets

### Root Level (`Makefile`)
- `start` – Start all services with Docker Compose
- `stop` – Stop all services
- `precommit-orderfoodonline` – Generate docs/mocks, build, and test orderfoodonline
- `precommit-coupons` – Generate docs/mocks, build, and test coupons service
- `test-api` – Run Postman collection tests with default settings
- `test-api-env` – Run Postman collection tests using environment file
- `test-api-npx` – Run Postman collection tests using npx (no global install required)
- `test-api-with-key` – Run Postman collection tests with custom API key
- `install-newman` – Install Newman CLI tool for Postman testing

### Service Level (e.g., `backend-challenge/services/orderfoodonline/Makefile`)
- `dep` – Run `go mod tidy`
- `build` – Build Docker image
- `test` – Run all Go tests with coverage
- `generate-mocks` – Generate GoMock mocks
- `generate-docs` – Generate Swagger docs
- `precommit` – Run all of the above for CI

---

## Pre-commit Hooks

Configured in `.pre-commit-config.yaml`:

- **Go Build/Format/Lint/Staticcheck** – Code quality and style
- **Go Test** – Automated testing with coverage
- **Go Vet** – Static analysis for common mistakes
- **GoSec** – Security vulnerability scanning
- **Gitlint** – Commit message style enforcement
- **Docs/Mocks Generation** – Automated documentation and mock creation

---

## How to Run

### **Quick Development (Recommended)**
```sh
make start
# API: http://localhost:8080
# Swagger UI: http://localhost:8080/swagger/index.html
# MongoDB: mongodb://localhost:27017
```

### **Individual Services**
```sh
# Start specific service
docker compose -f docker-compose.local.yml up orderfoodonline
docker compose -f docker-compose.local.yml up coupons

# Run tests
make -C backend-challenge/services/orderfoodonline test
make -C backend-challenge/library test
make -C backend-challenge/services/coupons test
```

### **Development with Hot Reload-TODO**
```sh
# The services are configured with air for hot reloading
# Changes to Go files will automatically restart the services
```

---

## API Testing

### **Postman Collection Testing**

The project includes a comprehensive Postman collection for testing all API endpoints. You can run the tests using Newman (Postman's CLI tool).

#### **Quick API Testing**
```sh
# Run tests with npx (recommended - no global install required)
make test-api-npx

# Run tests using environment file
make test-api-env

# Run tests with custom API key
make test-api-with-key
```

#### **Manual Setup**
```sh
# Option 1: Use npx (no installation required)
npx newman run "Order Food Online.postman_collection.json" \
  --environment "Order Food Online.postman_environment.json"

# Option 2: Install Newman globally
make install-newman
# or manually:
npm install -g newman --unsafe-perm=true
# or with sudo:
sudo npm install -g newman

# Option 3: Install via Homebrew (macOS)
brew install newman

# Run collection with environment file
newman run "Order Food Online.postman_collection.json" \
  --environment "Order Food Online.postman_environment.json" \
  --reporters cli,json \
  --reporter-json-export postman-results.json
```

#### **Test Results**
- **CLI Output**: Real-time test results in the terminal
- **JSON Report**: Detailed results saved to `postman-results.json`
- **Coverage**: Tests all endpoints including authentication, products, and orders
- **Validation**: Response time, status codes, and data structure validation

#### **Environment Variables**
Update `Order Food Online.postman_environment.json` to customize:
- `host`: API host (default: localhost)
- `port`: API port (default: 8080)
- `api_key`: Your API key for authentication
- `firstProductId`: Auto-populated by tests for dependent requests

---

## API Documentation

- **OpenAPI Spec:** [`api/openapi.yaml`](api/openapi.yaml)
- **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) (when running locally)
- **Auto-generated docs:** `backend-challenge/services/orderfoodonline/cmd/rest/docs/`
- **Health Check:** `GET /api/health`
- **Version Info:** `GET /api/version`
- **Prometheus Metrics:** `GET /metrics` - Comprehensive application metrics in Prometheus format

---

## Performance Features

### **Coupon Processing**
- **Parallel Processing**: 4 concurrent workers for database operations
- **Optimized Batching**: 5000 items per batch (5x improvement)
- **Memory Efficiency**: Pre-allocated slices and buffer reuse
- **Resume Capability**: Continue processing from failure point
- **File Deduplication**: MD5 hash-based duplicate detection

### **API Performance**
- **Rate Limiting**: Built-in request throttling
- **CORS Support**: Cross-origin resource sharing
- **Graceful Shutdown**: Proper cleanup and timeout handling
- **Connection Pooling**: Efficient database connection management

---

## Metrics & Observability

### **Prometheus Metrics Endpoint**
The application exposes comprehensive metrics at `/metrics` endpoint in Prometheus format for monitoring and alerting.

#### **Available Metrics**
- **HTTP Metrics**
  - `http_requests_total` - Total request count by method, endpoint, and status code
  - `http_request_duration_seconds` - Request duration histogram by method and endpoint

- **Database Metrics**
  - `database_queries_total` - Database operation counts by operation, collection, and status
  - `database_query_duration_seconds` - Query duration histogram by operation and collection
  - `database_active_connections` - Number of active database connections

- **Business Metrics**
  - `order_processing_duration_seconds` - Order processing time by status
  - `orders_total` - Order counts by status (success, validation_error, etc.)

#### **Usage Examples**
```bash
# View metrics in browser
curl http://localhost:8080/metrics

# Scrape with Prometheus
# Add to prometheus.yml:
scrape_configs:
  - job_name: 'orderfoodonline'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'

# Monitor specific metrics
curl http://localhost:8080/metrics | grep http_requests_total
curl http://localhost:8080/metrics | grep database_query_duration_seconds
```

#### **Integration**
- **Grafana Dashboards**: Create custom dashboards for business and technical metrics
- **Alerting**: Set up alerts for high error rates, slow queries, or business anomalies
- **SLA Monitoring**: Track API response times and availability
- **Capacity Planning**: Monitor database connection usage and query performance

---

## Testing Strategy

### **Unit Testing Approach**
- **Interface-Based Design**: All dependencies are abstracted through interfaces
- **Generated Mocks**: Automated mock generation using GoMock for consistent testing
- **Comprehensive Coverage**: Tests for all public methods and error scenarios
- **Human-Readable Tests**: Clear Given/When/Then structure with descriptive test names
- **Mock Validation**: Proper expectation verification and cleanup

### **Test Structure**
```go
func TestService_Method_Success(t *testing.T) {
    // Given: Setup with mocked dependencies
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockRepository(ctrl)
    
    // When: Execute the method under test
    result, err := service.Method(ctx, input)
    
    // Then: Verify expectations and assertions
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### **Mock Generation**
```sh
# Generate mocks for interfaces
make generate-mocks

# Mocks are automatically generated in mocks/ directories
# and used in unit tests for isolated testing
```

---

## Extending the Project

### **Adding New Microservices**
1. Create new service under `backend-challenge/services/`
2. Follow the established structure (cmd/, internal/, Dockerfile)
3. Define interfaces for dependencies and generate mocks
4. Update root Makefile with new targets
5. Add to docker-compose.local.yml

### **Adding Shared Utilities**
1. Add to `backend-challenge/library/`
2. Include comprehensive tests and documentation
3. Update go.mod dependencies
4. Generate mocks for any new interfaces

### **Database Migrations**
1. Add migration files to `migrations/` directory
2. Follow the existing naming convention (0001_, 0002_, etc.)
3. Include both schema changes and seed data

### **Adding New Repository Methods**
1. Define the method in the interface
2. Implement in the concrete repository
3. Generate mocks: `make generate-mocks`
4. Write comprehensive unit tests with proper mocking

---

## Troubleshooting

### **Common Issues**
- **Port conflicts**: Ensure ports 8080 and 27017 are available
- **Permission issues**: Run `chmod +x` on shell scripts if needed
- **Docker issues**: Ensure Docker daemon is running
- **Mock generation**: Run `make generate-mocks` if mocks are outdated

### **Development Tips**
- Use `make stop && make start` to restart all services
- Check logs with `docker compose logs -f service-name`
- Run tests with coverage: `make -C backend-challenge/services/orderfoodonline test`
- Regenerate mocks after interface changes: `make generate-mocks`

---

## Resources

- [API documentation (live)](https://orderfoodonline.deno.dev/public/openapi.html)
- [API specification (yaml)](https://orderfoodonline.deno.dev/public/openapi.yaml)
- [Figma design file](./design.fig)
- [Red Hat Text font](https://fonts.google.com/specimen/Red+Hat+Text)

---

## License

MIT or as specified in this repository.

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `pre-commit run --all-files` to ensure quality
5. Submit a pull request

If you have any questions or want to contribute, please open an issue or pull request!

