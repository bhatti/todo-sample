package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
)

// ---- UserStore tests ----

func TestUserStore_Create(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	u := &model.User{
		ID:        "01J1Z000000000000000000A",
		Username:  "alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.Create(ctx, u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, err := s.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 user, got %d", len(list))
	}
}

func TestUserStore_Get_NotFound(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	_, err := s.Get(ctx, "nonexistent")
	if err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserStore_List(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	for i, u := range []model.User{
		{ID: "id1", Username: "alice", Email: "alice@example.com"},
		{ID: "id2", Username: "bob", Email: "bob@example.com"},
	} {
		uu := u
		if err := s.Create(ctx, &uu); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}
	list, err := s.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 users, got %d", len(list))
	}
}

func TestUserStore_Update(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	u := &model.User{ID: "id1", Username: "alice", Email: "alice@example.com"}
	if err := s.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}
	u.Username = "alice2"
	if err := s.Update(ctx, u); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err := s.Get(ctx, "id1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Username != "alice2" {
		t.Errorf("expected username alice2, got %s", got.Username)
	}
}

func TestUserStore_Delete(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	u := &model.User{ID: "id1", Username: "alice", Email: "alice@example.com"}
	if err := s.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.Delete(ctx, "id1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := s.Get(ctx, "id1")
	if err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

// ---- TodoStore tests ----

func TestTodoStore_Create(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	todo := &model.Todo{
		ID:     "t1",
		UserID: "u1",
		Title:  "Buy milk",
	}
	if err := s.Create(ctx, todo); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.Get(ctx, "t1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %q", got.Title)
	}
}

func TestTodoStore_ListByUser(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	todos := []model.Todo{
		{ID: "t1", UserID: "userA", Title: "A1"},
		{ID: "t2", UserID: "userA", Title: "A2"},
		{ID: "t3", UserID: "userB", Title: "B1"},
	}
	for i := range todos {
		if err := s.Create(ctx, &todos[i]); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}
	list, err := s.ListByUser(ctx, "userA")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 todos for userA, got %d", len(list))
	}
}

func TestTodoStore_Get_NotFound(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	_, err := s.Get(ctx, "unknown")
	if err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestTodoStore_Update(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	todo := &model.Todo{ID: "t1", UserID: "u1", Title: "Task", Status: model.StatusPending}
	if err := s.Create(ctx, todo); err != nil {
		t.Fatalf("create: %v", err)
	}
	todo.Status = model.StatusDone
	if err := s.Update(ctx, todo); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err := s.Get(ctx, "t1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Status != model.StatusDone {
		t.Errorf("expected StatusDone, got %v", got.Status)
	}
}

func TestTodoStore_Delete(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	todo := &model.Todo{ID: "t1", UserID: "u1", Title: "Task"}
	if err := s.Create(ctx, todo); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.Delete(ctx, "t1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := s.Get(ctx, "t1")
	if err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}
