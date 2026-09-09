# Go Enum Zero-Value Pitfall

## Problem
When using `iota` in Go, the first constant equals `0`, which is also the zero value for `int`. This means an uninitialized struct field silently takes on the first enum value.

## Example (bad)
```go
type Status int
const (
    StatusPending Status = iota // = 0, same as zero-value
    StatusInProgress
)
// Todo{} silently has Status=StatusPending — may be wrong
```

## Fix options
1. Start iota at 1 and define an explicit unset sentinel:
```go
type Status int
const (
    StatusUnset   Status = iota // = 0 — explicit sentinel
    StatusPending               // = 1
    StatusInProgress            // = 2
)
```
2. Validate non-zero in Validate():
```go
if t.Status == 0 { return errors.New("status is required") }
```

## When to flag in review
Flag any `iota` enum where the first constant is semantically meaningful (not an explicit Unset/Unknown) and is used as a struct field without a Validate() guard.