# Implementation Plan

## Issue
Create a simple todo data model in Go with fields: id, title, description, completed, created_at.

## Task Breakdown

### Task 1: Define Todo struct with JSON tags
**Priority:** P0  
**Estimate:** S  
**Wave:** 1 (foundation)  
**Mode:** AFK

**Files to create:**
- `todo.go` — package `todo`, defines the `Todo` struct
- `todo_test.go` — tests for JSON marshaling and field presence

**Description:**
Create a Go package with a `Todo` struct containing all required fields with proper JSON tags and types.

**Acceptance Criteria:**
- `Todo` struct has fields: `ID string`, `Title string`, `Description string`, `Completed bool`, `CreatedAt time.Time`
- All fields have JSON tags matching snake_case: `id`, `title`, `description`, `completed`, `created_at`
- Package compiles with `go build ./...`
- Tests verify JSON round-trip (marshal → unmarshal produces identical struct)
- `go test ./...` passes

---

## Summary

| # | Title | Priority | Estimate | Wave | Mode | Deps |
|---|-------|----------|----------|------|------|------|
| 1 | Define Todo struct with JSON tags | P0 | S | 1 | AFK | — |

**Total tasks:** 1  
**Critical path:** Task 1  
**Total complexity:** S
