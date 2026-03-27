# 🧠 TrellGo — AI & Engineering Guidelines

## 🎯 Goal

Build a Trello-like application (**TrellGo**) using:

* Go (Golang)
* Modular Monolith architecture
* Clean Architecture (lightweight)
* PostgreSQL + sqlc
* Fiber (HTTP API)

The system should be:

* scalable
* easy to maintain
* easy to evolve into microservices

---

## 🏗️ Architecture

### Core Principle

👉 **Feature-Based Modular Monolith**

Each feature (module) is self-contained:

* `domain` → business logic & interfaces
* `application` → use cases
* `infrastructure` → DB, external systems
* `delivery` → HTTP handlers

---

## 📁 Project Structure

```
/cmd/api/main.go

/internal/
  /app/
    server.go
    router.go
    container.go

  /platform/
    /db/
    /config/
    /logger/

  /modules/
    /boards/
    /lists/
    /cards/
    /users/

/db/
  /migrations/
  /queries/
```

---

## 🧩 Domain Model (Trello-like)

### Entities

* **Board**

  * id
  * name
  * owner_id
  * created_at
  * updated_at

* **List**

  * id
  * board_id
  * position
  * name
  * created_at
  * updated_at

* **Card**

  * id
  * list_id
  * title
  * description
  * position
  * created_at
  * updated_at

* **User**

  * id
  * email
  * password_hash
  * created_at
  * updated_at

---

## 🔥 Core Rules

### 1. Feature-Based Structure ONLY

✅ Correct:

```
modules/boards
modules/cards
```

❌ Forbidden:

```
/handlers
/services
/repositories
```

---

### 2. Domain Layer

* Pure Go (no frameworks)
* No DB access
* Defines:

  * entities
  * interfaces

```go
type Board struct {
    ID      string
    Name    string
    OwnerID string
}
```

---

### 3. Repository Pattern

* Interface in `domain`
* Implementation in `infrastructure`

```go
type Repository interface {
    Create(board *Board) error
}
```

---

### 4. Use Case Pattern

Each action = one use case.

Examples:

* `CreateBoardUseCase`
* `CreateListUseCase`
* `CreateCardUseCase`
* `MoveCardUseCase`

❌ Avoid:

* generic "BoardService"

---

### 5. Delivery Layer (Fiber)

Responsibilities:

* parse request
* call use case
* return response

❌ No business logic allowed

---

### 6. Infrastructure Layer

* PostgreSQL via sqlc
* external integrations

---

### 7. Dependency Injection

* Manual wiring in `container.go`
* No DI frameworks

---

## 🗄️ Database Design Rules

### Use PostgreSQL

### Positioning (CRUCIAL)

Approach: Integer + GAP-based ordering

All ordered entities (Lists, Cards) must use:

```sql
position INTEGER NOT NULL
```

---

#### 📏 Positioning Rules
Default spacing

Positions must be spaced using gaps:

1000, 2000, 3000, ...

👉 This allows inserting items between without reordering all records.

---

#### ➕ Insert at End

When adding a new item at the end:

```sql
SELECT COALESCE(MAX(position), 0) + 1000
FROM cards
WHERE list_id = $1;
```

---

#### ↕️ Insert Between Items

To insert between two items:

```go
newPosition := (prev.Position + next.Position) / 2
```

---

#### ❗ Edge Case: No Space Between Positions

If:

next.Position - prev.Position <= 1

Then:

👉 Rebalancing is required

---

#### 🔄 Rebalancing Strategy

Reassign positions with fresh gaps:

```sql
WITH ordered AS (
  SELECT id, ROW_NUMBER() OVER (ORDER BY position) as rn
  FROM cards
  WHERE list_id = $1
)
UPDATE cards
SET position = ordered.rn * 1000
FROM ordered
WHERE cards.id = ordered.id;
```

---

#### ⚙️ Move Operation Algorithm

All move operations MUST follow this logic:

```go
if next.Position - prev.Position > 1 {
    newPosition = (prev.Position + next.Position) / 2
} else {
    rebalance(listID)
    retryMove()
}
```

---

### Example Tables

```sql
CREATE TABLE boards (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  owner_id TEXT NOT NULL
  created_at TIMESTAMPZ NOT NULL
  updated_at TIMESTAMPZ NOT NULL
);

CREATE TABLE lists (
  id TEXT PRIMARY KEY,
  board_id TEXT NOT NULL,
  position INT NOT NULL
  name TEXT NOT NULL,
  created_at TIMESTAMPZ NOT NULL
  updated_at TIMESTAMPZ NOT NULL
);

CREATE TABLE cards (
  id TEXT PRIMARY KEY,
  list_id TEXT NOT NULL,
  title TEXT NOT NULL,
  description TEXT,
  position INT NOT NULL
  created_at TIMESTAMPZ NOT NULL
  updated_at TIMESTAMPZ NOT NULL
);
```

---

## ⚙️ sqlc Rules

* ALL queries must be written in SQL
* NO ORM allowed

Example:

```sql
-- name: CreateCard :exec
INSERT INTO cards (id, list_id, title, position)
VALUES ($1, $2, $3, $4);
```

---

## 📡 API Design

### REST Endpoints

* `POST /users/signup`

* `POST /users/signin`

* `POST /boards`

* `GET /boards/:id`

* `POST /lists`

* `POST /cards`

* `POST /cards/:id/move`

---

## 🔄 Example Flow

```
HTTP Request
  → Handler (Fiber)
    → Use Case
      → Repository
        → DB (sqlc)
```

---

## 🧪 Testing Strategy

* Unit test use cases
* Mock repositories
* No DB in domain tests

---

## 🚀 Development Workflow

When implementing a feature:

1. Define domain entity
2. Define repository interface
3. Write SQL (sqlc)
4. Implement repository
5. Implement use case
6. Add HTTP handler
7. Wire in container

---

## ⚠️ Anti-Patterns (STRICTLY FORBIDDEN)

* ❌ Business logic in handlers
* ❌ Shared "utils" folder chaos
* ❌ ORM usage
* ❌ Global state
* ❌ Circular dependencies
* ❌ Microservices on day 1

---

## 🧠 Scaling Strategy

Start as monolith.

Later:

* Extract modules into services
* Keep domain boundaries intact

---

## 💡 Philosophy

* Prefer simplicity over cleverness
* Explicit > implicit
* Small focused use cases
* Design for change

---

## 🏁 Final Rule

👉 If you're unsure — keep it simple and modular.

---
