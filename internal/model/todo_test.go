package model

import (
	"testing"
	"time"
)

func TestPriorityString(t *testing.T) {
	cases := []struct {
		p    Priority
		want string
	}{
		{PriorityLow, "Low"},
		{PriorityMedium, "Medium"},
		{PriorityHigh, "High"},
		{PriorityUrgent, "Urgent"},
		{Priority(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.p.String(); got != c.want {
			t.Errorf("Priority(%d).String() = %q, want %q", c.p, got, c.want)
		}
	}
}

func TestStatusString(t *testing.T) {
	cases := []struct {
		s    Status
		want string
	}{
		{StatusPending, "Pending"},
		{StatusInProgress, "In Progress"},
		{StatusDone, "Done"},
		{StatusCancelled, "Cancelled"},
		{Status(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}

func TestTodoFields(t *testing.T) {
	now := time.Now()
	due := now.Add(48 * time.Hour)

	todo := Todo{
		ID:          "01J1Z0000000000000000001",
		UserID:      "01J1Z0000000000000000000",
		Title:       "Write tests",
		Description: "Cover all model types",
		Priority:    PriorityHigh,
		Status:      StatusInProgress,
		DueAt:       &due,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if todo.UserID != "01J1Z0000000000000000000" {
		t.Errorf("unexpected UserID: %s", todo.UserID)
	}
	if todo.Priority != PriorityHigh {
		t.Errorf("unexpected Priority: %v", todo.Priority)
	}
	if todo.Status != StatusInProgress {
		t.Errorf("unexpected Status: %v", todo.Status)
	}
	if todo.DueAt == nil || !todo.DueAt.Equal(due) {
		t.Errorf("unexpected DueAt")
	}
}

func TestTodoDueAtNilByDefault(t *testing.T) {
	todo := Todo{
		ID:     "01J1Z0000000000000000002",
		UserID: "01J1Z0000000000000000000",
		Title:  "No deadline",
	}
	if todo.DueAt != nil {
		t.Errorf("expected nil DueAt for zero-value Todo, got %v", todo.DueAt)
	}
}

func TestTodoAllPriorities(t *testing.T) {
	priorities := []Priority{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}
	for _, p := range priorities {
		todo := Todo{Priority: p}
		if todo.Priority != p {
			t.Errorf("expected priority %v, got %v", p, todo.Priority)
		}
	}
}

func TestTodoAllStatuses(t *testing.T) {
	statuses := []Status{StatusPending, StatusInProgress, StatusDone, StatusCancelled}
	for _, s := range statuses {
		todo := Todo{Status: s}
		if todo.Status != s {
			t.Errorf("expected status %v, got %v", s, todo.Status)
		}
	}
}
