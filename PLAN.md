# Implementation Plan — Todo App Data Model (Go)

**Issue:** #1 — Create data model for user and todo tasks  
**Branch:** `ai/1-ceate-data-model-081bf2a2`

---

## Overview

Set up a Go project from scratch and define the core domain data models (`User`, `TodoTask`) with priority, status, and timeline fields, plus an in-memory repository layer that validates the models work end-to-end.

---

## Task Summary

| ID | Title | Priority | Size | Wave | Mode | Depends On |
|----|-------|----------|------|------|------|------------|
| [task-001](tasks/backlog/task-001.md) | Go Project Scaffold | P0 | S | 1 | AFK | — |
| [task-002](tasks/backlog/task-002.md) | User and TodoTask Data Models | P0 | S | 2 | AFK | task-001 |
| [task-003](tasks/backlog/task-003.md) | In-Memory Repository Layer | P0 | S | 2 | AFK | task-002 |

**Total:** 3 tasks · All AFK · All P0 · Total complexity: M

---

## Dependency Waves

```
Wave 1 (foundation)
└── task-001: Go Project Scaffold
      go.mod, Makefile, directory layout, .gitignore

Wave 2 (domain — sequential, task-003 uses task-002's types)
├── task-002: User and TodoTask Data Models
│     internal/model/user.go, internal/model/task.go
│     Priority enum, Status enum, DueDate, Validate()
└── task-003: In-Memory Repository Layer
      internal/repository/user_repository.go
      internal/repository/task_repository.go
      CRUD interfaces + thread-safe in-memory impls
```

---

## Directory Layout (target state)

```
.
├── cmd/
│   └── todo/
│       └── main.go          # entrypoint stub
├── internal/
│   ├── model/
│   │   ├── user.go
│   │   ├── user_test.go
│   │   ├── task.go
│   │   └── task_test.go
│   └── repository/
│       ├── errors.go
│       ├── user_repository.go
│       ├── user_repository_test.go
│       ├── task_repository.go
│       └── task_repository_test.go
├── go.mod
├── go.sum
├── Makefile
└── .gitignore
```

---

## Key Design Decisions

- **Module path:** `github.com/todo-app/todo`
- **Go version:** 1.22+
- **No external dependencies** — stdlib only for models and repositories
- **Priority:** `PriorityLow | PriorityMedium | PriorityHigh` with `String()` method
- **Status:** `StatusTodo | StatusInProgress | StatusDone`
- **Timeline:** `DueDate *time.Time` on TodoTask (optional)
- **IDs:** plain `string` — callers choose scheme
- **Thread safety:** `sync.RWMutex` in in-memory repositories
- **Error handling:** sentinel errors (`ErrNotFound`, `ErrAlreadyExists`) for use with `errors.Is`

---

## Acceptance (Definition of Done)

```bash
make build       # zero errors
make test        # go test -race ./... — all pass
```

All three tasks completed, `go test -race ./...` green, no external deps added.
