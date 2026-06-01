package model

import (
	"testing"
	"time"
)

func TestPriority_String(t *testing.T) {
	cases := []struct {
		p    Priority
		want string
	}{
		{PriorityLow, "Low"},
		{PriorityMedium, "Medium"},
		{PriorityHigh, "High"},
		{Priority(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.p.String(); got != c.want {
			t.Errorf("Priority(%d).String() = %q, want %q", c.p, got, c.want)
		}
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    Status
		want string
	}{
		{StatusTodo, "Todo"},
		{StatusInProgress, "In Progress"},
		{StatusDone, "Done"},
		{Status(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}

func TestNewTodoTask_FieldsPopulated(t *testing.T) {
	task := NewTodoTask("t1", "u1", "Buy milk")
	if task.ID != "t1" {
		t.Errorf("expected ID t1, got %s", task.ID)
	}
	if task.UserID != "u1" {
		t.Errorf("expected UserID u1, got %s", task.UserID)
	}
	if task.Title != "Buy milk" {
		t.Errorf("expected Title 'Buy milk', got %s", task.Title)
	}
	if task.Priority != PriorityLow {
		t.Errorf("expected default Priority Low, got %v", task.Priority)
	}
	if task.Status != StatusTodo {
		t.Errorf("expected default Status Todo, got %v", task.Status)
	}
	if task.DueDate != nil {
		t.Error("expected DueDate nil by default")
	}
	if task.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNewTodoTask_WithDueDate(t *testing.T) {
	due := time.Now().Add(24 * time.Hour)
	task := NewTodoTask("t1", "u1", "Buy milk")
	task.DueDate = &due
	if task.DueDate == nil {
		t.Error("expected DueDate to be set")
	}
}

func TestTodoTask_Validate_Valid(t *testing.T) {
	task := NewTodoTask("t1", "u1", "Buy milk")
	if err := task.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTodoTask_Validate_MissingTitle(t *testing.T) {
	task := NewTodoTask("t1", "u1", "")
	if err := task.Validate(); err == nil {
		t.Error("expected error for empty Title")
	}
}

func TestTodoTask_Validate_MissingID(t *testing.T) {
	task := NewTodoTask("", "u1", "Buy milk")
	if err := task.Validate(); err == nil {
		t.Error("expected error for empty ID")
	}
}

func TestTodoTask_Validate_MissingUserID(t *testing.T) {
	task := NewTodoTask("t1", "", "Buy milk")
	if err := task.Validate(); err == nil {
		t.Error("expected error for empty UserID")
	}
}
