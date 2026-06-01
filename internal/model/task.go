package model

import (
	"errors"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// Priority represents the urgency level of a task.
type Priority int

const (
	PriorityLow      Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

// Status represents the current state of a task.
type Status int

const (
	StatusOpen       Status = iota
	StatusInProgress
	StatusDone
	StatusCancelled
)

// Task is a todo item owned by a User.
type Task struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Status      Status
	Priority    Priority
	DueDate     time.Time
	StartDate   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTask creates a validated Task with Status=StatusOpen and a generated ID.
func NewTask(title, userID string, priority Priority, dueDate time.Time) (*Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title must not be empty")
	}
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("userID must not be empty")
	}
	if dueDate.IsZero() {
		return nil, errors.New("dueDate must not be zero")
	}
	if !dueDate.After(time.Now().UTC()) {
		return nil, errors.New("dueDate must be in the future")
	}
	now := time.Now().UTC()
	return &Task{
		ID:        ulid.Make().String(),
		UserID:    userID,
		Title:     title,
		Status:    StatusOpen,
		Priority:  priority,
		DueDate:   dueDate,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
