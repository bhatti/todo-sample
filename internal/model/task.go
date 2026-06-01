package model

import (
	"errors"
	"strings"
	"time"
)

// Priority represents how urgent a task is.
type Priority int

const (
	PriorityLow    Priority = iota // 0
	PriorityMedium                 // 1
	PriorityHigh                   // 2
)

// String returns the human-readable label for a Priority value.
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "Low"
	case PriorityMedium:
		return "Medium"
	case PriorityHigh:
		return "High"
	default:
		return "Unknown"
	}
}

// Status represents the lifecycle stage of a task.
type Status int

const (
	StatusTodo       Status = iota // 0
	StatusInProgress               // 1
	StatusDone                     // 2
)

// String returns the human-readable label for a Status value.
func (s Status) String() string {
	switch s {
	case StatusTodo:
		return "Todo"
	case StatusInProgress:
		return "In Progress"
	case StatusDone:
		return "Done"
	default:
		return "Unknown"
	}
}

// TodoTask is a single item in a user's task list.
type TodoTask struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Priority    Priority
	Status      Status
	DueDate     *time.Time // optional timeline field
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTodoTask creates a TodoTask with default Priority=Low and Status=Todo.
func NewTodoTask(id, userID, title string) *TodoTask {
	now := time.Now().UTC()
	return &TodoTask{
		ID:        id,
		UserID:    userID,
		Title:     title,
		Priority:  PriorityLow,
		Status:    StatusTodo,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate returns an error if required fields are missing.
func (t *TodoTask) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return errors.New("task: ID is required")
	}
	if strings.TrimSpace(t.UserID) == "" {
		return errors.New("task: UserID is required")
	}
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("task: Title is required")
	}
	return nil
}
