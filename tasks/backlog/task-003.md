# Todo Task Data Model with Priority and Timeline

**ID:** task-003
**Source:** ISSUE_CONTEXT.md
**Priority:** P0 (MUST)
**Estimate:** M
**Wave:** 2
**Execution:** AFK
**Depends on:** task-002

## Description
Define the `Task` struct in `internal/model/task.go` with fields for title, description, status, priority (enum), and timeline (due date + optional start date). Include the `Priority` and `Status` enum types, constructor, validation, and comprehensive unit tests.

## Files to Change
- `internal/model/task.go` — new (Task struct, Priority enum, Status enum, NewTask, Validate)
- `internal/model/task_test.go` — new

## Acceptance Criteria

### Scenario: Valid task creation
- **GIVEN** a valid title, user ID, priority, and due date
- **WHEN** `NewTask(title, userID, priority, dueDate)` is called
- **THEN** a `Task` is returned with a non-empty ID, Status=Open, CreatedAt set, no error

### Scenario: Priority enum covers all levels
- **GIVEN** the Priority type
- **WHEN** all declared constants are listed
- **THEN** Low, Medium, High, and Critical are present and distinct

### Scenario: Past due date rejected
- **GIVEN** a due date in the past
- **WHEN** `NewTask(...)` is called
- **THEN** an error is returned

### Scenario: Unit tests pass
- **GIVEN** the task model package
- **WHEN** `go test ./internal/model/...` is run
- **THEN** all tests pass with 0 failures

## Checklist
- [ ] Implementation complete
- [ ] Tests passing
- [ ] Acceptance scenarios verified

## Notes
Fields: ID (string), UserID (string, FK to User), Title (string), Description (string), Status (enum: Open/InProgress/Done/Cancelled), Priority (enum: Low/Medium/High/Critical), DueDate (time.Time), StartDate (*time.Time, optional), CreatedAt (time.Time), UpdatedAt (time.Time). Use `iota` for enums with String() method. Validation: title non-empty, userID non-empty, DueDate not zero.
