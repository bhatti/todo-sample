# Go Project Setup with Module and Structure

**ID:** task-001
**Source:** ISSUE_CONTEXT.md
**Priority:** P0 (MUST)
**Estimate:** S
**Wave:** 1
**Execution:** AFK
**Depends on:** none

## Description
Initialize the Go module, establish the standard project layout (cmd/, internal/, pkg/), and configure basic tooling (go.mod, .gitignore, Makefile with build/test targets). This is the foundational scaffold all other tasks depend on.

## Files to Change
- `go.mod` — new
- `.gitignore` — new
- `Makefile` — new
- `cmd/todoapp/main.go` — new (minimal entrypoint)

## Acceptance Criteria

### Scenario: Module initializes and builds
- **GIVEN** an empty repo
- **WHEN** `go build ./...` is run
- **THEN** it succeeds with zero errors

### Scenario: Test target works
- **GIVEN** the project scaffold is in place
- **WHEN** `make test` is run
- **THEN** `go test ./...` executes with exit code 0

## Checklist
- [ ] Implementation complete
- [ ] Tests passing
- [ ] Acceptance scenarios verified

## Notes
Use module name `github.com/user/todoapp` (or project-appropriate name). Go version >= 1.21. Standard layout: `cmd/` for executables, `internal/` for private packages, `pkg/` for reusable public packages.
