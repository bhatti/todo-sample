# Todo API

A REST API built with Node.js + Express + SQLite (sql.js) to manage users and their todo task lists.

## Requirements

- Node.js ≥ 18
- npm ≥ 9

## Setup

```bash
npm install
```

## Running

```bash
npm start          # production
npm run dev        # development (auto-reload via nodemon)
```

The server starts on port `3000` by default. Override with `PORT` env var.

## Testing

```bash
npm test
```

Tests use an in-memory SQLite database — no setup required.

## API Reference

### Users

| Method | Path         | Description           | Status Codes          |
|--------|--------------|-----------------------|-----------------------|
| POST   | /users       | Create a user         | 201, 400, 409         |
| GET    | /users       | List all users        | 200                   |
| GET    | /users/:id   | Get a user            | 200, 404              |
| PUT    | /users/:id   | Update a user         | 200, 400, 404, 409    |
| DELETE | /users/:id   | Delete a user         | 204, 404              |

**POST /users** — Request body:
```json
{ "username": "alice", "email": "alice@example.com" }
```

### Todos (nested under a user)

| Method | Path                         | Description          | Status Codes       |
|--------|------------------------------|----------------------|--------------------|
| POST   | /users/:userId/todos         | Create a todo        | 201, 400, 404      |
| GET    | /users/:userId/todos         | List user todos      | 200, 404           |
| GET    | /users/:userId/todos/:id     | Get a todo           | 200, 404           |
| PUT    | /users/:userId/todos/:id     | Update a todo        | 200, 400, 404      |
| PATCH  | /users/:userId/todos/:id     | Toggle completed     | 200, 404           |
| DELETE | /users/:userId/todos/:id     | Delete a todo        | 204, 404           |

**POST /users/:userId/todos** — Request body:
```json
{
  "title": "Buy milk",
  "description": "Whole milk, 1 gallon",
  "priority": "medium",
  "status": "pending",
  "due_at": "2026-12-31T18:00:00Z"
}
```

Priority values: `low`, `medium` (default), `high`, `urgent`  
Status values: `pending` (default), `in_progress`, `done`, `cancelled`

### Error Responses

All errors return JSON in this shape:
```json
{ "error": "<message>" }
```

| HTTP Status | Meaning                       |
|-------------|-------------------------------|
| 400         | Validation error              |
| 404         | Resource not found            |
| 409         | Unique constraint violation   |
| 500         | Internal server error         |

## Project Structure

```
src/
├── app.js              Express app factory
├── server.js           HTTP server entry point
├── db.js               SQLite connection & schema migrations
├── schema.sql          Schema reference (documentation)
├── models/
│   ├── user.js         User data-access layer
│   └── todo.js         Todo data-access layer
├── routes/
│   ├── users.js        User CRUD routes
│   └── todos.js        Todo CRUD routes
└── middleware/
    ├── errorHandler.js Centralised error handler
    └── validateJson.js Malformed JSON handler
tests/
├── helpers/db.js       In-memory test DB helper
├── models/             Unit tests for model layer
└── routes/             Integration tests via supertest
```
