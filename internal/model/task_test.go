package model

import (
	"testing"
	"time"
)

func TestPriorityConstants(t *testing.T) {
	priorities := []Priority{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical}
	seen := map[Priority]bool{}
	for _, p := range priorities {
		if seen[p] {
			t.Errorf("duplicate Priority value: %d", p)
		}
		seen[p] = true
	}
}

func TestStatusConstants(t *testing.T) {
	statuses := []Status{StatusOpen, StatusInProgress, StatusDone, StatusCancelled}
	seen := map[Status]bool{}
	for _, s := range statuses {
		if seen[s] {
			t.Errorf("duplicate Status value: %d", s)
		}
		seen[s] = true
	}
}

func TestNewTask_Valid(t *testing.T) {
	due := time.Now().UTC().Add(24 * time.Hour)
	task, err := NewTask("Buy milk", "user-123", PriorityHigh, due)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.ID == "" {
		t.Error("ID must not be empty")
	}
	if task.Status != StatusOpen {
		t.Errorf("Status = %d, want StatusOpen", task.Status)
	}
	if task.Priority != PriorityHigh {
		t.Errorf("Priority = %d, want PriorityHigh", task.Priority)
	}
}

func TestNewTask_EmptyTitle(t *testing.T) {
	due := time.Now().UTC().Add(24 * time.Hour)
	_, err := NewTask("", "user-123", PriorityLow, due)
	if err == nil {
		t.Error("expected error for empty title")
	}
}

func TestNewTask_EmptyUserID(t *testing.T) {
	due := time.Now().UTC().Add(24 * time.Hour)
	_, err := NewTask("Buy milk", "", PriorityLow, due)
	if err == nil {
		t.Error("expected error for empty userID")
	}
}

func TestNewTask_PastDueDate(t *testing.T) {
	past := time.Now().UTC().Add(-24 * time.Hour)
	_, err := NewTask("Buy milk", "user-123", PriorityLow, past)
	if err == nil {
		t.Error("expected error for past due date")
	}
}

func TestNewTask_ZeroDueDate(t *testing.T) {
	_, err := NewTask("Buy milk", "user-123", PriorityLow, time.Time{})
	if err == nil {
		t.Error("expected error for zero due date")
	}
}
