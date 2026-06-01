# Implementation Plan: Todo App Data Model (Go)

**Issue:** #1 — Create data model for user and todo tasks with priority and timeline  
**Language:** Go  
**Date:** 2026-06-01

---

## Overview

Build a Go project from scratch implementing data models for `User` and `Task` (todo items). Tasks have priority levels and a timeline (due date, optional start date). Follow standard Go project layout (`cmd/`, `internal/`, `pkg/`) with tests for every model.

---

## Task Summary

| ID | Title | Priority | Estimate | Wave | Mode | Depends On |
|----|-------|----------|----------|------|------|------------|
| task-001 | Go project setup (module, layout, Makefile) | P0 | S | 1 | AFK | none |
| task-002 | User data model | P0 | S | 1 | AFK | task-001 |
| task-003 | Todo Task data model (priority + timeline) | P0 | M | 2 | AFK | task-002 |

---

## Dependency Waves

```
Wave 1 (parallel):
  task-001  Go project scaffold
  task-002  User model          ← depends on task-001

Wave 2:
  task-003  Task model (priority + timeline)  ← depends on task-002
```

---

## Task Details

### task-001 — Go Project Setup
**Estimate:** S (< half day) | **Wave:** 1 | **AFK**

Set up the Go module and project skeleton.

**Files:**
- `go.mod` — module `github.com/user/todoapp`, Go 1.21+
- `.gitignore` — Go standard ignores
- `Makefile` — `build`, `test`, `lint` targets
- `cmd/todoapp/main.go` — minimal `main()`

**Acceptance criteria:**
- `go build ./...` succeeds with zero errors
- `make test` runs `go test ./...` and exits 0

---

### task-002 — User Data Model
**Estimate:** S (< half day) | **Wave:** 1 | **AFK**

Define the `User` struct as the owner entity for tasks.

**Files:**
- `internal/model/user.go`
- `internal/model/user_test.go`

**Key fields:**
```go
type User struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Constructor:** `NewUser(name, email string) (*User, error)`  
**Validation:** name non-empty, email non-empty and contains `@`

**Acceptance criteria:**
- `NewUser("Alice", "alice@example.com")` returns a user with non-empty ID
- `NewUser("", "bad")` returns an error
- `go test ./internal/model/...` passes

---

### task-003 — Todo Task Data Model (Priority + Timeline)
**Estimate:** M (~1 day) | **Wave:** 2 | **AFK**

Define the `Task` struct with status, priority enum, and timeline fields.

**Files:**
- `internal/model/task.go`
- `internal/model/task_test.go`

**Key types:**
```go
type Priority int
const (
    PriorityLow Priority = iota
    PriorityMedium
    PriorityHigh
    PriorityCritical
)

type Status int
const (
    StatusOpen Status = iota
    StatusInProgress
    StatusDone
    StatusCancelled
)

type Task struct {
    ID          string
    UserID      string
    Title       string
    Description string
    Status      Status
    Priority    Priority
    DueDate     time.Time
    StartDate   *time.Time  // optional
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Constructor:** `NewTask(title, userID string, priority Priority, dueDate time.Time) (*Task, error)`  
**Validation:** title non-empty, userID non-empty, dueDate not zero and not in the past

**Acceptance criteria:**
- All four Priority constants exist and are distinct
- All four Status constants exist and are distinct
- `NewTask` with past due date returns an error
- `NewTask` with valid args returns Task with Status=StatusOpen
- `go test ./internal/model/...` passes

---

## Project Layout (Target)

```
.
├── cmd/
│   └── todoapp/
│       └── main.go
├── internal/
│   └── model/
│       ├── user.go
│       ├── user_test.go
│       ├── task.go
│       └── task_test.go
├── go.mod
├── .gitignore
├── Makefile
└── PLAN.md
```

---

## Verification

After all tasks are complete, run:

```bash
make build   # must succeed
make test    # must pass: go test ./...
```
