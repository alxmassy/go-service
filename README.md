# Task Service

A production-style Go backend service demonstrating clean architecture,
request lifecycle management, and PostgreSQL-backed persistence.

## Overview

This project is not a todo app or a frontend product.

It is a reference implementation of a production-grade Go backend service,
focusing on:
- explicit ownership and boundaries
- clean separation between HTTP handling and storage
- graceful shutdown and request cancellation
- observability (logging + request IDs)
- swappable storage implementations

## Architecture

The service is structured around three core layers:

### Entry Point (`cmd/server`)

- Responsible for application lifecycle
- Loads configuration from environment variables
- Initializes shared dependencies (database, store)
- Wires HTTP routes and middleware
- Handles graceful shutdown

### HTTP Layer (`internal/http`)

- Contains request handlers and middleware
- Depends only on interfaces, not concrete implementations
- Handles:
  - request parsing and validation
  - response formatting
  - error handling

### Storage Layer (`internal/store`)

- Owns the domain model (`Task`)
- Defines the `TaskStore` interface
- Provides multiple implementations:
  - `MemoryTaskStore` (for local/testing)
  - `PostgresTaskStore` (for persistence)

Handlers depend on the `TaskStore` interface, allowing the storage
implementation to be swapped without changing HTTP logic.

## Request Lifecycle

1. Incoming request is assigned a unique request ID
2. Request ID is injected into context
3. Logging middleware records method, path, duration, request ID
4. Handler parses and validates input
5. Handler calls store methods using request context
6. Store respects context cancellation
7. Response is written or error is returned

## Persistence

The service uses PostgreSQL for durable storage.

Schema:
CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  done BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);


- SQL is written explicitly using `database/sql`
- No ORM is used to keep behavior explicit and debuggable
- Context-aware queries are used to support cancellation and shutdown

## Error Handling

- All storage operations return explicit errors
- Handlers map internal errors to appropriate HTTP responses
- Context cancellation is respected across layers
- Fail-fast behavior is used for misconfiguration (missing env vars)

## Configuration

The service is configured via environment variables:

- PORT
- DB_HOST
- DB_PORT
- DB_USER
- DB_PASSWORD
- DB_NAME

export PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=taskuser
export DB_PASSWORD=taskpass
export DB_NAME=taskdb

## Running Locally

1. Ensure PostgreSQL is running
2. Create the database and schema
3. Set environment variables
4. Run:

go run cmd/server/main.go

## Tradeoffs and Future Improvements

- Authentication is intentionally omitted to keep focus on core backend concerns
- No migration tool is used yet; schema is managed manually
- Docker is not included to reduce cognitive overhead during development

Possible future improvements:
- Add request rate limiting
- Add DB-backed health checks
- Introduce migrations tooling
- Containerize for deployment
