# Golang Gin Template

A modular, testable, and containerized backend for a simple bookstore, built with Go, Gin, GORM, and Docker. The project follows best practices for configuration, authentication, and integration testing.

## Table of Contents

- [Simple Bookstore Backend](#simple-bookstore-backend)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Project Structure](#project-structure)
  - [Setup \& Installation](#setup--installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
  - [Environment Configuration](#environment-configuration)
  - [Docker Compose Services](#docker-compose-services)
  - [Testing](#testing)
  - [API Documentation](#api-documentation)
  - [Key Components](#key-components)
  - [Contributing](#contributing)
  - [License](#license)

---

## Features

- RESTful API for book and user management
- JWT-based authentication and authorization
- Integration tests with isolated Dockerized PostgreSQL
- Configurable via environment variables
- Modular package structure

---

## Project Structure

```
backend/
├── api/               # API route definitions and Swagger docs
├── benchmark/         # (Reserved for benchmarks)
├── cmd/               # Application entrypoints (main.go, etc.)
├── config/            # Configuration loading and management
├── db/                # Database migrations
├── docs/              # Generated API documentation (Swagger)
├── infrastructure/    # Auth, DB, Docker, Logger, etc.
├── internal/          # Business logic (books(with unittest), users(with unittest), middleware, utils)
├── pkg/               # Shared packages (crypto, migration, validator)
├── test/              # Integration tests
├── .env               # Environment variables (local)
├── .example.env       # Example environment variables
├── go.mod, go.sum     # Go dependencies
└── README.md          # Project documentation
```

---

## Setup & Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/chai-rs/go-gin-template.git
   cd go-gin-template
   ```

2. **Copy and configure environment variables:**
   ```sh
   cp .example.env .env
   # Edit .env as needed
   ```

3. **Install Go dependencies:**
   ```sh
   go mod download
   ```

4. **(Optional) Install Docker**  
   Required for running integration tests.

---

## Configuration

- All configuration is managed via the `config` package and `.env` file.
- Secrets (e.g., `ACCESS_SECRET`, `REFRESH_SECRET`) are set in the environment and can be overridden in test setup.

---

## Running the Application

```sh
go run ./cmd/server
```

Or use Docker Compose:

```sh
docker-compose --env-file .docker.env up --build
```

- The application will be available at http://localhost:8000
- Only the server port is exposed; Redis and Postgres are accessible only within the Docker network.
- Environment variables for Compose are managed in `.docker.env` (see below).

---

## Environment Configuration

- `.env` is for local development (e.g., running with `go run` on your host).
- `.docker.env` is for Docker Compose. Service hostnames must match the Compose service names (`backend-postgres`, `backend-redis`).
- See `.example.env` for all available variables and their expected values.

---

## Docker Compose Services

- **backend-server**: Your Go app, built from `/cmd/server`.
- **backend-postgres**: PostgreSQL 17, data stored in `pgdata` volume.
- **backend-redis**: Redis 7, data stored in `redisdata` volume.

All services are connected via the `bookstore-network` for secure internal communication.

---

## Testing

- Integration tests use `testify/suite` and `ory/dockertest` to spin up a PostgreSQL container, run migrations, and clean up automatically.
- Example test suite: `test/base_test.go`, `test/user_test.go`

Run all tests:

```sh
go test ./...
```

---

## API Documentation

- Swagger/OpenAPI docs are available in `docs/` and generated from code comments.
- To view locally, run the server and visit `/swagger/index.html` (if route is enabled).

---

## Key Components

- **Authentication:** JWT-based, with token generation and verification in `infrastructure/auth`.
- **Authorization:** Enforced via Casbin with Gorm adapter, configured in `auth_model.conf`.
- **Database:** GORM ORM with migrations in `db/migrations`.
- **Testing:** Uses Dockerized PostgreSQL for isolation, see `BaseSuite` in `test/base_test.go`.
- **Handlers & Services:** Business logic in `internal/`, separated by domain (books, users).

---

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes
4. Push to the branch (`git push origin feature/fooBar`)
5. Open a pull request

---

## License

MIT License
