# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Project Structure

1. **cmd/**: Contains the main application and seed command entry points.

   - `server/main.go`: The main application entry point, serves the REST API.
   - `seed/main.go`: Command to seed the database with initial product data.

2. **app/**: Contains the application logic.
3. **sql/**: Contains a very simple database migration scripts setup.
4. **models/**: Contains the data models and repositories used in the application.
5. `.env`: Environment variables file for configuration (not committed; see `.env.example`).

## Setup Code Repository

1. Create a github/bitbucket/gitlab repository and push all this code as-is.
2. Create a new branch, and provide a pull-request against the main branch with your changes. Instructions to follow.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Copy `.env.example` to `.env` and adjust values as needed:
  - `HTTP_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`, `POSTGRES_SQL_DIR`.

### Running the application

- **Install dependencies:**
  - `make tidy`
- **Start infrastructure (Postgres):**
  - `make docker-up`
- **Apply migrations and seed data:**
  - `make seed` (⚠️ will destroy and re-create the database tables)
- **Run the API server:**
  - `make run`
- **Stop infrastructure:**
  - `make docker-down`

### Running tests

- **Unit tests for all packages:**
  - `make test`  
    (equivalent to `go test -v -count=1 ./... -coverprofile=coverage.out -covermode=atomic`)
- **Repository integration tests (require Postgres and seed data):**
  - `go test -tags=integration ./repositories`

## API Overview

- `GET /catalog`
  - Query parameters:
    - `offset` (optional, default: 0; must be ≥ 0).
    - `limit` (optional, default: 10; must be in \[1; 100\]).
    - `category` (optional, category code; filters products by category).
    - `priceLessThan` (optional, > 0; filters products by price).
  - Successful response: `{"data":{"items":[...],"total":N}}`.

- `GET /catalog/{id}`
  - `id` is the product code (e.g. `PROD001`).
  - Successful response: `{"data":{"code":"...","price":...,"category":{...},"variants":[...]}}`.

  # 1) GET /catalog → 200
curl "http://localhost:8484/catalog"

# 2) GET /catalog?limit=101 → 400
curl "http://localhost:8484/catalog?limit=101"

# 3) GET /catalog?priceLessThan=0 → 400
curl "http://localhost:8484/catalog?priceLessThan=0"

# 4) GET /catalog/PROD001 → 200
curl "http://localhost:8484/catalog/PROD001"

# 5) GET /catalog/UNKNOWN → 404
curl "http://localhost:8484/catalog/UNKNOWN"
