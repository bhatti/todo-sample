# Implementation Plan: Create Data Model for Todo App (Issue #1)

## Summary
Set up a Go project with a data model for users and todo tasks, including priority and timeline support.

---

## Task 1: Initialize Go Project Structure
**Complexity**: S

### Description
Bootstrap the Go module and project layout following Go best practices (flat package layout for a small service, with `cmd/`, `internal/`, and test files alongside source).

### Files to Create
- `go.mod` — module definition (`module github.com/user/todo`)
- `go.sum` — dependency lockfile (auto-generated)
- `cmd/todo/main.go` — entrypoint (minimal, wires everything together)
- `internal/model/` — package directory for data models

### Steps
1. `go mod init github.com/user/todo`
2. Create directory structure: `cmd/todo/`, `internal/model/`
3. Add a minimal `main.go` that imports the model package and compiles

### Acceptance Criteria
- `go build ./...` succeeds with zero errors
- `go vet ./...` reports no issues
- Project layout matches Go standard layout conventions

---

## Task 2: Define User and Todo Data Models
**Complexity**: S

### Description
Define the core data structures in `internal/model/` — `User` and `Todo` with all required fields: priority, timeline (due date), status, and associations.

### Files to Create
- `internal/model/user.go` — `User` struct
- `internal/model/todo.go` — `Todo` struct, `Priority` enum, `Status` enum

### Data Model Design

```go
// user.go
type User struct {
    ID        string    // ULID
    Username  string
    Email     string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// todo.go
type Priority int
const (
    PriorityLow Priority = iota
    PriorityMedium
    PriorityHigh
    PriorityUrgent
)

type Status int
const (
    StatusPending Status = iota
    StatusInProgress
    StatusDone
    StatusCancelled
)

type Todo struct {
    ID          string    // ULID
    UserID      string    // FK → User.ID
    Title       string
    Description string
    Priority    Priority
    Status      Status
    DueAt       *time.Time  // optional timeline/deadline
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Acceptance Criteria
- `go build ./...` succeeds
- All exported types have godoc comments
- Priority and Status have `String()` methods
- `go vet ./...` reports no issues

---

## Task 3: Unit Tests for Data Models
**Complexity**: S

### Description
Write unit tests covering model construction, validation helpers, and the `String()` methods for enums. Tests live alongside the source in `internal/model/`.

### Files to Create
- `internal/model/user_test.go`
- `internal/model/todo_test.go`

### Test Coverage
- Construct a `User` and verify field assignment
- Construct a `Todo` with each `Priority` / `Status` value
- Verify `Priority.String()` and `Status.String()` return expected labels
- Verify zero-value `DueAt` is nil (optional timeline)
- Verify `UserID` association is stored correctly on `Todo`

### Acceptance Criteria
- `go test ./...` passes with zero failures
- Coverage ≥ 90% on `internal/model` package (`go test -cover ./internal/model/`)
- No data races (`go test -race ./...`)

---

## Dependency Order
```
Task 1 (project setup) → Task 2 (models) → Task 3 (tests)
```

Each task produces a independently compilable/testable increment.
