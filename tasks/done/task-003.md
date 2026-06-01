# In-Memory Repository Layer

**ID:** task-003
**Source:** ISSUE_CONTEXT.md
**Priority:** P0 (MUST)
**Estimate:** S
**Wave:** 2
**Mode:** AFK
**Depends on:** task-002

## Description
Define repository interfaces for `User` and `TodoTask` in `internal/repository/` and provide in-memory implementations. This makes the data models usable end-to-end (create, read, update, delete) and validates that the model structs integrate correctly with a storage contract. Add unit tests for CRUD operations.

## Acceptance Criteria

### Scenario: UserRepository interface is satisfied by in-memory impl
- **GIVEN** `InMemoryUserRepository` is instantiated
- **WHEN** `Create`, `GetByID`, `List`, `Update`, `Delete` are called in sequence
- **THEN** each operation returns the correct result with no errors for valid input, and `GetByID` returns `ErrNotFound` for missing IDs

### Scenario: TodoTaskRepository filters by user
- **GIVEN** tasks owned by two different users are stored
- **WHEN** `ListByUserID(userID)` is called
- **THEN** only tasks belonging to that user are returned

### Scenario: Concurrent safety
- **GIVEN** the in-memory repository
- **WHEN** multiple goroutines call `Create` and `List` concurrently
- **THEN** no data race is detected (`go test -race ./internal/repository/...`)

### Scenario: Unit tests pass
- **GIVEN** `internal/repository/` contains interface and impl files with `_test.go`
- **WHEN** `go test -race ./internal/repository/...` is run
- **THEN** all tests pass

## Checklist
- [ ] `internal/repository/user_repository.go` — `UserRepository` interface + `InMemoryUserRepository`
- [ ] `internal/repository/task_repository.go` — `TodoTaskRepository` interface + `InMemoryTodoTaskRepository`
- [ ] `internal/repository/errors.go` — `ErrNotFound`, `ErrAlreadyExists` sentinel errors
- [ ] `internal/repository/user_repository_test.go` — CRUD + error path tests
- [ ] `internal/repository/task_repository_test.go` — CRUD + `ListByUserID` + filter tests
- [ ] `go test -race ./...` passes

## Notes
- Use `sync.RWMutex` for thread safety in in-memory impl
- Interface methods: `Create`, `GetByID`, `List`/`ListByUserID`, `Update`, `Delete`
- Return `ErrNotFound` (not nil pointer) when ID is missing — callers use `errors.Is`
- No external dependencies — stdlib only
