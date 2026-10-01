# micro-blog

REST API for a micro-blogging platform. Built with Go, Chi, PostgreSQL, and Wire DI.

## Prerequisites

- Go 1.26+
- Docker & Docker Compose
- [goose](https://github.com/pressly/goose) (migrations)
- [golangci-lint](https://golangci-lint.run) (linting)

## Quick start

```bash
cp .env.example .env
# Set JWT_SECRET in .env

# Start Postgres
make dc.up

# Run migrations
make mg.up

# Run the app
go run ./api/v0
```

## API

All endpoints return JSON. Authenticated routes require `Authorization: Bearer <token>`.

| Method | Path                       | Auth | Description    |
|--------|----------------------------|------|----------------|
| POST   | `/api/users`               | No   | Register user  |
| POST   | `/api/login`               | No   | Login (JWT)    |
| POST   | `/api/posts`               | Yes  | Create post    |
| GET    | `/api/posts`               | Yes  | List all posts |
| POST   | `/api/posts/{postId}/like` | Yes  | Like a post    |

## Makefile targets

Run `make help` for a full list.

| Target         | Description                          |
|----------------|--------------------------------------|
| `make test`    | Run all tests                        |
| `make race`    | Run tests with race detector         |
| `make bench`   | Run benchmarks                       |
| `make check`   | Lint + race + bench                  |
| `make lint`    | Run golangci-lint                    |
| `make wire`    | Regenerate Wire DI                   |
| `make dc.up`   | Start Docker containers              |
| `make dc.down` | Stop and remove containers           |
| `make mg.up`   | Apply migrations                     |
| `make mg.down` | Rollback last migration              |

## Docker

Build and run with Docker Compose:

```bash
# Build the app image
docker build -f .ops/docker/Dockerfile -t micro-blog .

# Or use compose (uncomment app service in docker-compose.yaml)
make dc.rebuild
```

The Dockerfile uses a multi-stage build: Go builder + minimal Alpine runtime. Migrations are included in the image at `./migrations`.

## Tests

Unit tests (mocked, no dependencies):

```bash
make test
```

Integration tests (requires Docker daemon for testcontainers):

```bash
go test ./internal/service/test/... -v
```

## Project structure

```
api/v0/             Entrypoint
internal/
  app/               Wire DI setup
  config/            Environment config
  database/          DB pool, DBTX interface
  domain/            Domain models
  dto/               Data transfer objects
  exception/         Error types
  handler/           HTTP handlers
  middleware/        Auth, error handling
  repository/        Data access layer
  request/           Request structs + validation
  response/          Response structs
  service/           Business logic
    auth/            JWT, password hashing
    post/            Post CRUD, likes
    user/            User registration
migrations/          Goose SQL migrations
route/               Router registration
test_help/           Shared test utilities
```
