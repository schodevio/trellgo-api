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
SECRET_KEY=<hex-encoded 32-byte key>
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
| `make test` | Run all tests |
| `make migrate-up` | Apply pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-status` | Show migration status |
| `make migrate-create` | Create a new migration file |
| `make sqlc` | Regenerate database code from SQL |
| `make docs` | Regenerate OpenAPI spec (`docs/`) |

## API Documentation

Interactive documentation is available when the server is running:

| URL | Description |
|---|---|
| `/swagger` | Swagger UI — interactive, supports "Try it out" |
| `/docs` | ReDoc — clean three-panel layout |

## Project Structure

```
cmd/api/          # Entrypoint
internal/
  app/            # Server bootstrap, router, middleware, DI container
  module/         # Feature modules (each owns handler/service/repository)
    auth/         # Sign up, sign in, sign out, token refresh
    boards/       # Boards CRUD (scoped to a user)
    lists/        # Lists CRUD (scoped to a board)
    cards/        # Cards CRUD (scoped to a list)
    health/       # Health check
  platform/       # Shared infrastructure
    apierrors/    # Typed API errors
    config/       # Environment config
    db/           # pgxpool setup
    docs/         # Swagger UI + ReDoc route registration
    middleware/   # Auth, CORS, logger
    validator/    # Request validation → i18n-friendly error keys
db/
  migrations/     # Goose SQL migrations
  queries/        # sqlc SQL queries
  sqlc/           # Generated Go code (do not edit)
docs/             # Generated OpenAPI spec (do not edit manually)
```

## API

> Endpoints marked with a lock require a Bearer token: `Authorization: Bearer <access_token>`

### Health

<details>
<summary><code>GET /api/v1/health</code></summary>

Response `200`:
```json
{ "status": "ok" }
```

</details>

### Auth

<details>
<summary><code>POST /api/v1/auth/signup</code></summary>

Request:
```json
{ "email": "user@example.com", "password": "secret123" }
```
Response `201`:
```json
{ "email": "user@example.com" }
```

</details>

<details>
<summary><code>POST /api/v1/auth/signin</code></summary>

Request:
```json
{ "email": "user@example.com", "password": "secret123" }
```
Response `200`:
```json
{ "access_token": "<token>" }
```
The refresh token is set as an `HttpOnly` cookie.

</details>

<details>
<summary><code>POST /api/v1/auth/refresh</code></summary>

Reads the refresh token from the `HttpOnly` cookie and issues a new access token.

Response `200`:
```json
{ "access_token": "<token>" }
```

</details>

<details>
<summary><code>DELETE /api/v1/auth/signout</code></summary>

Revokes the refresh token cookie.

Response `204`.

</details>

### Boards

All boards endpoints require a Bearer token:
```
Authorization: Bearer <access_token>
```

<details>
<summary><code>POST /api/v1/boards</code></summary>

Request:
```json
{ "name": "My Board" }
```
Response `201`:
```json
{ "board": { "id": "...", "name": "My Board", "user_id": "...", "created_at": "...", "updated_at": "..." } }
```

</details>

<details>
<summary><code>GET /api/v1/boards</code></summary>

Response `200`:
```json
{ "boards": [ { "id": "...", "name": "My Board", "user_id": "...", "created_at": "...", "updated_at": "..." } ] }
```

</details>

<details>
<summary><code>GET /api/v1/boards/:id</code></summary>

Response `200`:
```json
{ "board": { "id": "...", "name": "My Board", "user_id": "...", "created_at": "...", "updated_at": "..." } }
```
Returns `404` if the board does not belong to the authenticated user.

</details>

<details>
<summary><code>PATCH /api/v1/boards/:id</code></summary>

Request:
```json
{ "name": "Renamed Board" }
```
Response `200`:
```json
{ "board": { "id": "...", "name": "Renamed Board", "user_id": "...", "created_at": "...", "updated_at": "..." } }
```
Returns `404` if the board does not belong to the authenticated user.

</details>

<details>
<summary><code>DELETE /api/v1/boards/:id</code></summary>

Response `204`. Returns `404` if the board does not belong to the authenticated user.

</details>

### Lists

All lists endpoints require a Bearer token:
```
Authorization: Bearer <access_token>
```

<details>
<summary><code>POST /api/v1/boards/:board_id/lists</code></summary>

Request:
```json
{ "name": "To Do", "position": 0 }
```
Response `201`:
```json
{ "list": { "id": "...", "name": "To Do", "board_id": "...", "position": 0, "created_at": "...", "updated_at": "..." } }
```

</details>

<details>
<summary><code>GET /api/v1/boards/:board_id/lists</code></summary>

Response `200`:
```json
{ "lists": [ { "id": "...", "name": "To Do", "board_id": "...", "position": 0, "created_at": "...", "updated_at": "..." } ] }
```
Returns `404` if the board does not belong to the authenticated user.

</details>

<details>
<summary><code>PATCH /api/v1/lists/:id</code></summary>

Request:
```json
{ "name": "In Progress" }
```
Response `200`:
```json
{ "list": { "id": "...", "name": "In Progress", "board_id": "...", "position": 1, "created_at": "...", "updated_at": "..." } }
```
Returns `404` if the list does not belong to the authenticated user.

</details>

<details>
<summary><code>PATCH /api/v1/lists/:id/move</code></summary>

Request:
```json
{ "position": 2 }
```
Response `200`:
```json
{ "list": { "id": "...", "name": "...", "board_id": "...", "position": 2, "created_at": "...", "updated_at": "..." } }
```
Returns `404` if the list does not belong to the authenticated user.

</details>

<details>
<summary><code>DELETE /api/v1/lists/:id</code></summary>

Response `204`. Returns `404` if the list does not belong to the authenticated user.

</details>

### Cards

All cards endpoints require a Bearer token:
```
Authorization: Bearer <access_token>
```

<details>
<summary><code>POST /api/v1/lists/:list_id/cards</code></summary>

Request:
```json
{ "title": "My Card", "description": "Optional description", "position": 0 }
```
Response `201`:
```json
{ "card": { "id": "...", "list_id": "...", "title": "My Card", "description": "Optional description", "position": 0, "created_at": "...", "updated_at": "..." } }
```

</details>

<details>
<summary><code>GET /api/v1/lists/:list_id/cards</code></summary>

Response `200`:
```json
{ "cards": [ { "id": "...", "list_id": "...", "title": "My Card", "description": "...", "position": 0, "created_at": "...", "updated_at": "..." } ] }
```
Returns `404` if the list does not belong to the authenticated user.

</details>

<details>
<summary><code>PATCH /api/v1/cards/:id</code></summary>

Request:
```json
{ "title": "Updated Card", "description": "Updated description" }
```
Response `200`:
```json
{ "card": { "id": "...", "list_id": "...", "title": "Updated Card", "description": "Updated description", "position": 0, "created_at": "...", "updated_at": "..." } }
```
Returns `404` if the card does not belong to the authenticated user.

</details>

<details>
<summary><code>PATCH /api/v1/cards/:id/move</code></summary>

Request:
```json
{ "list_id": "...", "position": 1 }
```
Response `200`:
```json
{ "card": { "id": "...", "list_id": "...", "title": "...", "description": "...", "position": 1, "created_at": "...", "updated_at": "..." } }
```
Returns `404` if the card does not belong to the authenticated user.

</details>

<details>
<summary><code>DELETE /api/v1/cards/:id</code></summary>

Response `204`. Returns `404` if the card does not belong to the authenticated user.

</details>

### Error format

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
