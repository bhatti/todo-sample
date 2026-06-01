# User Data Model

**ID:** task-002
**Source:** ISSUE_CONTEXT.md
**Priority:** P0 (MUST)
**Estimate:** S
**Wave:** 1
**Execution:** AFK
**Depends on:** task-001

## Description
Define the `User` struct in `internal/model/user.go` with fields for identity, profile, and timestamps. Include constructor, validation, and unit tests. This is the owner entity for all todo tasks.

## Files to Change
- `internal/model/user.go` — new
- `internal/model/user_test.go` — new

## Acceptance Criteria

### Scenario: Valid user creation
- **GIVEN** a valid name and email
- **WHEN** `NewUser(name, email)` is called
- **THEN** a `User` is returned with a non-empty ID, CreatedAt set, and no error

### Scenario: Invalid email rejected
- **GIVEN** an empty or malformed email string
- **WHEN** `NewUser(name, email)` is called
- **THEN** an error is returned and no user is created

### Scenario: Unit tests pass
- **GIVEN** the user model package
- **WHEN** `go test ./internal/model/...` is run
- **THEN** all tests pass with 0 failures

## Checklist
- [ ] Implementation complete
- [ ] Tests passing
- [ ] Acceptance scenarios verified

## Notes
Fields: ID (string, UUID/ULID), Name (string), Email (string), CreatedAt (time.Time), UpdatedAt (time.Time). Validation: name non-empty, email non-empty and contains `@`. No external ORM — plain structs.
