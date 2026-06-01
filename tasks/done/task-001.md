# Go Project Scaffold

**ID:** task-001
**Source:** ISSUE_CONTEXT.md
**Priority:** P0 (MUST)
**Estimate:** S
**Wave:** 1
**Mode:** AFK
**Depends on:** none

## Description
Initialize a Go project for the todo app using best practices: Go module, standard directory layout (`cmd/`, `internal/`, `pkg/`), `.gitignore`, `Makefile` with `build`/`test`/`lint` targets, and a minimal `README.md`. No business logic — just the project skeleton that all other tasks build on.

## Acceptance Criteria

### Scenario: Module initializes cleanly
- **GIVEN** the repo root has no Go files
- **WHEN** `go mod tidy` is run
- **THEN** `go.mod` exists with module path `github.com/todo-app/todo` and Go 1.22+, zero errors

### Scenario: Build target succeeds
- **GIVEN** the scaffold is in place
- **WHEN** `make build` is run
- **THEN** exits 0 with no compilation errors

### Scenario: Test target runs (empty suite is OK)
- **GIVEN** the scaffold is in place
- **WHEN** `make test` is run
- **THEN** exits 0 (no tests yet is fine, `go test ./...` must not error)

## Checklist
- [ ] `go.mod` created with correct module path and Go version
- [ ] `go.sum` generated (after `go mod tidy`)
- [ ] `cmd/todo/main.go` stub present (compiles, prints nothing or a placeholder)
- [ ] `internal/` directory created for domain packages
- [ ] `.gitignore` covers `*.exe`, `bin/`, `vendor/`, IDE files
- [ ] `Makefile` with `build`, `test`, `lint` targets
- [ ] `make build` and `make test` pass

## Notes
- Module path: `github.com/todo-app/todo`
- Use `internal/model` for domain types (added in task-002)
- No external dependencies needed for the scaffold itself
