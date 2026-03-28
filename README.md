# TrellGo API

A Trello-like project management REST API built with Go.

## Tech Stack

| Layer | Tool |
|---|---|
| Framework | [Fiber v3](https://github.com/gofiber/fiber) |
| Database | PostgreSQL via [pgx v5](https://github.com/jackc/pgx) |
| Queries | [sqlc](https://sqlc.dev) |
| Migrations | [Goose](https://github.com/pressly/goose) |
| Validation | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Hot reload | [Air](https://github.com/air-verse/air) |

## Getting Started

The project uses a **Dev Container** (Docker Compose). Open in VS Code with the Dev Containers extension — everything is provisioned automatically.

### Environment

Copy and adjust the default values in `.env`:

```env
DB_URL=postgres://postgres:postgres@postgres:5432/trellgo_development?sslmode=disable
PORT=3000
```

### Run migrations

```sh
make migrate-up
```

### Start the server

```sh
# Development (hot reload)
make dev

# Production build
make build && ./tmp/main
```

## Make Targets

| Command | Description |
|---|---|
| `make dev` | Start server with hot reload (Air) |
| `make build` | Compile binary to `./tmp/main` |
| `make run` | Run without building a binary |
| `make migrate-up` | Apply pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-status` | Show migration status |
| `make migrate-create` | Create a new migration file |
| `make sqlc` | Regenerate database code from SQL |

## Project Structure

```
cmd/api/          # Entrypoint
internal/
  app/            # Server bootstrap, router, middleware, container
  module/         # Feature modules (each owns handler/service/repository)
    auth/
  platform/       # Shared infrastructure
    apierrors/    # Typed API errors
    config/       # Environment config
    db/           # pgxpool setup
    validator/    # Request validation → i18n-friendly error keys
db/
  migrations/     # Goose SQL migrations
  queries/        # sqlc SQL queries
  sqlc/           # Generated Go code (do not edit)
```

## API

### Health

```
GET /health
```

### Auth

```
POST /api/v1/auth/signup
```

**Request**
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

**Response `201`**
```json
{
  "email": "user@example.com"
}
```

**Validation errors `422`**
```json
{
  "error": {
    "code": "unprocessable_entity",
    "message": "validation_failed",
    "details": {
      "email": ["required"],
      "password": ["min:8"]
    }
  }
}
```

Error keys in `details` are i18n-friendly (e.g. `required`, `email`, `min:8`) and intended to be translated on the frontend.
