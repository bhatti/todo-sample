package store

import (
	"strings"
	"testing"

	"github.com/user/todo/internal/model"
)

func TestMemoryUserStore_Create(t *testing.T) {
	s := NewMemoryUserStore()
	u, err := s.Create(&model.User{Username: "alice", Email: "alice@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == "" {
		t.Error("expected non-empty ID")
	}
	if u.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if u.Username != "alice" {
		t.Errorf("expected username alice, got %s", u.Username)
	}
}

func TestMemoryUserStore_GetByID(t *testing.T) {
	s := NewMemoryUserStore()
	created, _ := s.Create(&model.User{Username: "alice", Email: "alice@example.com"})

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, got.ID)
	}
}

func TestMemoryUserStore_GetByID_NotFound(t *testing.T) {
	s := NewMemoryUserStore()
	_, err := s.GetByID("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestMemoryUserStore_Update(t *testing.T) {
	s := NewMemoryUserStore()
	created, _ := s.Create(&model.User{Username: "alice", Email: "alice@example.com"})

	updated, err := s.Update(&model.User{ID: created.ID, Username: "bob", Email: "bob@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Username != "bob" {
		t.Errorf("expected username bob, got %s", updated.Username)
	}
	if !updated.UpdatedAt.After(created.CreatedAt) {
		t.Error("expected UpdatedAt > CreatedAt")
	}
}

func TestMemoryUserStore_Delete(t *testing.T) {
	s := NewMemoryUserStore()
	created, _ := s.Create(&model.User{Username: "alice", Email: "alice@example.com"})

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

func TestMemoryUserStore_List(t *testing.T) {
	s := NewMemoryUserStore()
	s.Create(&model.User{Username: "a", Email: "a@example.com"})
	s.Create(&model.User{Username: "b", Email: "b@example.com"})
	s.Create(&model.User{Username: "c", Email: "c@example.com"})

	users, err := s.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}
