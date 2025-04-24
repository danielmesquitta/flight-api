# Flight API

A Go-based HTTP API for searching flights, authentication, health checks, and OpenAPI/Swagger documentation. Built with Fiber, Wire for DI, Go‑Playground/Validator for validation, JWT for auth, and Redis for rate‑limiting.

## Features

- Health check endpoint (`GET /api/health`)
- User login with email/password (`POST /api/v1/auth/login`)
- JWT‑based authentication middleware for protected routes
- Flight search endpoint (`GET /api/v1/flights/search`)
- OpenAPI/Swagger docs served under `/api/docs`
- Docker support & Makefile commands
- Unit & integration tests with testify, Fiber’s test harness, Dockerized Redis

## Project Structure

```sh
.
├── .air.toml # Air config for live reload
├── .editorconfig
├── .env.example
├── Dockerfile # Docker configuration to deploy the server
├── Makefile # Makefile for scripts
├── bin
│ └── install.sh # Install Go tools
├── cmd
│ └── server
│  └── main.go # Main entry point for the server
├── docs
│ ├── docs.go # Documentation handling
│ ├── openapi.json
│ ├── openapi.yaml
│ ├── swagger.json
│ └── swagger.yaml
├── embed.go # Embedding files for the server
├── generate.go # Code generation utilities
├── go.mod
├── go.sum
├── internal
│ ├── app
│ │ └── server # HTTP server, handlers, middleware, router
│ ├── config # env loading, logging, time zone, Wire setup
│ ├── domain # use‑cases, entities, error types
│ ├── pkg # utilities: jwtutil, validator, ptr, …
│ └── provider # external integrations (Amadeus, Duffel, Serp, cache)
├── test
│ ├── container # Docker containers for integration tests
│ └── integration # Integration tests
└── tmp
└── .gitkeep
```

## Prerequisites

- [Go 1.21+](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/)
- [make](https://www.gnu.org/software/make/)

## Installation

1. Clone the repo and enter directory:

```sh
git clone https://github.com/danielmesquitta/flight-api.git
cd flight-api
```

2. Install tools:

```sh
make install
```

3. Copy .env.example to .env and set required values

4. Run the server locally:

```sh
make run
```

or Run with Docker:

```sh
docker build -t flight-api .
docker run -p 8080:8080 flight-api
```

## API Documentation

Swagger UI and raw specs are available at:

/api/docs/swagger.json
/api/docs/swagger.yaml
/api/docs/openapi.json
/api/docs/openapi.yaml
/api/docs

## Testing

### Unit Tests

Run unit tests with:

```sh
make unit-test
```

### Integration Tests

Run integration tests with:

```sh
make integration-tests
```

### Run both unit and integration tests:

```sh
make tests
```
