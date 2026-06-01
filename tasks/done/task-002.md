# User and TodoTask Data Models

**ID:** task-002
**Source:** ISSUE_CONTEXT.md
**Priority:** P0 (MUST)
**Estimate:** S
**Wave:** 2
**Mode:** AFK
**Depends on:** task-001

## Description
Define the core domain types in `internal/model/`: a `User` struct and a `TodoTask` struct with priority enum, status enum, and timeline fields. Include constructor functions and a `Validate()` method on each type. Add unit tests covering field validation and enum values.

## Acceptance Criteria

### Scenario: User struct has required fields
- **GIVEN** the `model` package is imported
- **WHEN** a `User` is created with `NewUser(id, name, email)`
- **THEN** the struct has `ID` (string), `Name` (string), `Email` (string), `CreatedAt` (time.Time), `UpdatedAt` (time.Time)

### Scenario: TodoTask struct has required fields
- **GIVEN** the `model` package is imported
- **WHEN** a `TodoTask` is created with `NewTodoTask(id, userID, title)`
- **THEN** the struct has `ID`, `UserID`, `Title`, `Description`, `Priority` (Priority enum), `Status` (Status enum), `DueDate` (*time.Time), `CreatedAt`, `UpdatedAt`

### Scenario: Priority enum covers all levels
- **GIVEN** the Priority type is defined
- **WHEN** all constants are enumerated
- **THEN** `PriorityLow`, `PriorityMedium`, `PriorityHigh` exist with distinct values and `String()` returns human-readable labels

### Scenario: Validation rejects bad input
- **GIVEN** a TodoTask with an empty title
- **WHEN** `task.Validate()` is called
- **THEN** an error is returned describing the missing field

### Scenario: Unit tests pass
- **GIVEN** `internal/model/` contains model definitions and `_test.go` files
- **WHEN** `go test ./internal/model/...` is run
- **THEN** all tests pass with zero failures

## Checklist
- [ ] `internal/model/user.go` — User struct, NewUser constructor, Validate method
- [ ] `internal/model/task.go` — TodoTask struct, Priority enum, Status enum, NewTodoTask constructor, Validate method
- [ ] `internal/model/user_test.go` — tests for User construction and validation
- [ ] `internal/model/task_test.go` — tests for TodoTask construction, enum values, validation
- [ ] `make test` passes

## Notes
- Status values: `StatusTodo`, `StatusInProgress`, `StatusDone`
- Priority values: `PriorityLow`, `PriorityMedium`, `PriorityHigh`
- `DueDate` is optional (pointer) — represents the timeline field
- IDs are plain strings; callers may use any ID scheme (UUIDs, etc.)
- No external dependencies — use only stdlib
