package store

import (
	"strings"
	"testing"

	"github.com/user/todo/internal/model"
)

func TestMemoryTodoStore_Create(t *testing.T) {
	s := NewMemoryTodoStore()
	todo, err := s.Create(&model.Todo{
		UserID:   "user1",
		Title:    "Buy milk",
		Priority: model.PriorityMedium,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.ID == "" {
		t.Error("expected non-empty ID")
	}
	if todo.Status != model.StatusPending {
		t.Errorf("expected StatusPending, got %v", todo.Status)
	}
	if todo.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestMemoryTodoStore_ListByUser(t *testing.T) {
	s := NewMemoryTodoStore()
	s.Create(&model.Todo{UserID: "user1", Title: "Task 1"})
	s.Create(&model.Todo{UserID: "user1", Title: "Task 2"})
	s.Create(&model.Todo{UserID: "user2", Title: "Task 3"})

	todos, err := s.ListByUser("user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("expected 2 todos for user1, got %d", len(todos))
	}
	for _, td := range todos {
		if td.UserID != "user1" {
			t.Errorf("expected UserID user1, got %s", td.UserID)
		}
	}
}

func TestMemoryTodoStore_GetByID_NotFound(t *testing.T) {
	s := NewMemoryTodoStore()
	_, err := s.GetByID("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestMemoryTodoStore_Update(t *testing.T) {
	s := NewMemoryTodoStore()
	created, _ := s.Create(&model.Todo{UserID: "user1", Title: "Buy milk", Priority: model.PriorityLow})

	updated, err := s.Update(&model.Todo{
		ID:       created.ID,
		UserID:   "user1",
		Title:    "Buy oat milk",
		Priority: model.PriorityHigh,
		Status:   model.StatusInProgress,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "Buy oat milk" {
		t.Errorf("expected updated title, got %s", updated.Title)
	}
	if updated.Status != model.StatusInProgress {
		t.Errorf("expected StatusInProgress, got %v", updated.Status)
	}
}

func TestMemoryTodoStore_Delete(t *testing.T) {
	s := NewMemoryTodoStore()
	created, _ := s.Create(&model.Todo{UserID: "user1", Title: "Task"})

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("unexpected error on delete: %v", err)
	}
	_, err := s.GetByID(created.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}
