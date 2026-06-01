package model

import "time"

// Priority represents the urgency level of a todo task.
type Priority int

const (
	// PriorityLow is for tasks with no time pressure.
	PriorityLow Priority = iota
	// PriorityMedium is for tasks with moderate urgency.
	PriorityMedium
	// PriorityHigh is for tasks that should be done soon.
	PriorityHigh
	// PriorityUrgent is for tasks requiring immediate attention.
	PriorityUrgent
)

// String returns a human-readable label for the Priority.
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "Low"
	case PriorityMedium:
		return "Medium"
	case PriorityHigh:
		return "High"
	case PriorityUrgent:
		return "Urgent"
	default:
		return "Unknown"
	}
}

// Status represents the lifecycle state of a todo task.
type Status int

const (
	// StatusPending is the initial state for a new task.
	StatusPending Status = iota
	// StatusInProgress means the task is actively being worked on.
	StatusInProgress
	// StatusDone means the task has been completed.
	StatusDone
	// StatusCancelled means the task was abandoned.
	StatusCancelled
)

// String returns a human-readable label for the Status.
func (s Status) String() string {
	switch s {
	case StatusPending:
		return "Pending"
	case StatusInProgress:
		return "In Progress"
	case StatusDone:
		return "Done"
	case StatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

// Todo represents a task belonging to a user.
type Todo struct {
	// ID is a ULID string uniquely identifying the task.
	ID string
	// UserID references the owning User.ID.
	UserID string
	// Title is a short summary of the task.
	Title string
	// Description provides additional detail about the task.
	Description string
	// Priority indicates the urgency of the task.
	Priority Priority
	// Status tracks the lifecycle state of the task.
	Status Status
	// DueAt is the optional deadline for the task.
	DueAt *time.Time
	// CreatedAt is when the task was created.
	CreatedAt time.Time
	// UpdatedAt is when the task was last modified.
	UpdatedAt time.Time
}
