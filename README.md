# Kart Challenge – Food Ordering Platform

A robust, modular food ordering backend and supporting library, designed for extensibility and correctness. This project features a RESTful API for product listing, cart management, and order processing, as well as a shared Go library for logging, configuration, and integrations.

---

## Project Structure

```bash
kart-challenge/
├── backend-challenge/
│   ├── library/                # Shared Go library (logger, config, slack, etc.)
│   ├── services/
│   │   ├── coupons/            # Coupons microservice (skeleton)
│   │   │   └── cmd/processor/  # Coupons processor entrypoint
│   │   └── orderfoodonline/    # Main food ordering service
│   │       ├── cmd/rest/       # REST API entrypoint and docs
│   │       ├── internal/       # Service internals (handlers, middlewares, routes, etc.)
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

## Features

- **Product Listing & Cart API**: RESTful endpoints for products, cart, and order management.
- **Shared Go Library**: Centralized logging, configuration, and Slack integration.
- **Pre-commit Quality Gates**: Lint, format, staticcheck, security scan, and tests.
- **API Documentation**: Swagger/OpenAPI 3.1 docs, auto-generated with `swag`.
- **Dockerized**: Multi-stage Dockerfiles for efficient builds and minimal runtime images.
- **Extensible Microservices**: Coupons service skeleton for future discount/offer logic.
- **Makefile Automation**: Common tasks for build, test, docs, and pre-commit.

---

## Development Setup

### Prerequisites

- [Go 1.23+](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Make](https://www.gnu.org/software/make/)
- [Nix (optional)](https://nixos.org/) for reproducible environments
- [pre-commit](https://pre-commit.com/) (install via pip or your package manager)

### Clone and Prepare

```sh
git clone <your-repo-url>
cd kart-challenge
pre-commit install
```

### Local Development

- **Start all services (with MongoDB):**
  
  ```sh
  make start
  ```

- **Stop all services:**
  
  ```sh
  make stop
  ```

- **Run all pre-commit checks:**
  
  ```sh
  pre-commit run --all-files
  ```

---

## Makefile Targets

From the root `Makefile`:

- `start` – Start all services with Docker Compose
- `stop` – Stop all services
- `precommit-orderfoodonline` – Generate docs/mocks, build, and test orderfoodonline
- `precommit-coupons` – (Reserved for coupons service)

From the service `Makefile` (e.g., `backend-challenge/services/orderfoodonline/Makefile`):

- `dep` – Run `go mod tidy`
- `build` – Build Docker image
- `test` – Run all Go tests with coverage
- `generate-mocks` – Generate GoMock mocks
- `generate-docs` – Generate Swagger docs
- `precommit` – Run all of the above for CI

---

## Pre-commit Hooks

Configured in `.pre-commit-config.yaml`:

- **Go Build/Format/Lint/Staticcheck**
- **Go Test** (library, orderfoodonline, coupons)
- **Go Vet** (library, orderfoodonline)
- **GoSec** (security scan)
- **Gitlint** (commit message style)
- **Docs/Mocks Generation** (via Makefile)

---

## How to Run

### With Docker Compose

```sh
make start
# or directly:
docker compose -f docker-compose.local.yml up --build
```

- API will be available at [http://localhost:8080](http://localhost:8080)
- MongoDB at [mongodb://localhost:27017](mongodb://localhost:27017)

### Run Tests

```sh
make -C backend-challenge/services/orderfoodonline test
make -C backend-challenge/library test
make -C backend-challenge/services/coupons test
```

### Generate API Docs

```sh
make -C backend-challenge/services/orderfoodonline generate-docs
```

---

## API Documentation

- **OpenAPI Spec:** [`api/openapi.yaml`](api/openapi.yaml)
- **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) (when running locally)
- **Auto-generated docs:** `backend-challenge/services/orderfoodonline/cmd/rest/docs/`

---

## Extending the Project

- Add new microservices under `backend-challenge/services/`
- Add shared utilities to `backend-challenge/library/`
- Update pre-commit hooks and Makefiles as needed

## TODOs

- **Hot Reload**: Local development with live reload using `air`.
  
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

If you have any questions or want to contribute, please open an issue or pull request!

