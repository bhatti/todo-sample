package store_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
)

func newUser(id, username, email string) *model.User {
	now := time.Now().UTC()
	return &model.User{ID: id, Username: username, Email: email, CreatedAt: now, UpdatedAt: now}
}

func newTodo(id, userID, title string) *model.Todo {
	now := time.Now().UTC()
	return &model.Todo{ID: id, UserID: userID, Title: title, CreatedAt: now, UpdatedAt: now}
}

func TestUserStore_CreateAndGet(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	u := newUser("u1", "alice", "alice@example.com")
	if err := s.Create(ctx, u); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetByID(ctx, "u1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Username != "alice" {
		t.Fatalf("expected alice, got %s", got.Username)
	}
}

func TestUserStore_GetNotFound(t *testing.T) {
	s := store.NewMemoryUserStore()
	_, err := s.GetByID(context.Background(), "missing")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserStore_List(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	_ = s.Create(ctx, newUser("u1", "alice", "a@e.com"))
	_ = s.Create(ctx, newUser("u2", "bob", "b@e.com"))
	list, err := s.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 users, got %d", len(list))
	}
}

func TestUserStore_Update(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	_ = s.Create(ctx, newUser("u1", "alice", "a@e.com"))
	u, _ := s.GetByID(ctx, "u1")
	u.Username = "alice2"
	if err := s.Update(ctx, u); err != nil {
		t.Fatal(err)
	}
	got, _ := s.GetByID(ctx, "u1")
	if got.Username != "alice2" {
		t.Fatalf("expected alice2, got %s", got.Username)
	}
}

func TestUserStore_UpdateNotFound(t *testing.T) {
	s := store.NewMemoryUserStore()
	err := s.Update(context.Background(), newUser("missing", "x", "x@e.com"))
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserStore_Delete(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	_ = s.Create(ctx, newUser("u1", "alice", "a@e.com"))
	if err := s.Delete(ctx, "u1"); err != nil {
		t.Fatal(err)
	}
	_, err := s.GetByID(ctx, "u1")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestUserStore_DeleteNotFound(t *testing.T) {
	s := store.NewMemoryUserStore()
	err := s.Delete(context.Background(), "missing")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserStore_Conflict(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	_ = s.Create(ctx, newUser("u1", "alice", "a@e.com"))
	err := s.Create(ctx, newUser("u2", "alice", "b@e.com"))
	if !errors.Is(err, store.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestUserStore_ConcurrentReadWrite(t *testing.T) {
	s := store.NewMemoryUserStore()
	ctx := context.Background()
	_ = s.Create(ctx, newUser("u1", "alice", "a@e.com"))
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.GetByID(ctx, "u1")
			_, _ = s.List(ctx)
		}()
	}
	wg.Wait()
}

func TestTodoStore_CreateAndGet(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	td := newTodo("t1", "u1", "buy milk")
	if err := s.Create(ctx, td); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetByID(ctx, "t1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "buy milk" {
		t.Fatalf("expected buy milk, got %s", got.Title)
	}
}

func TestTodoStore_GetNotFound(t *testing.T) {
	s := store.NewMemoryTodoStore()
	_, err := s.GetByID(context.Background(), "missing")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTodoStore_ListByUser(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	_ = s.Create(ctx, newTodo("t1", "u1", "a"))
	_ = s.Create(ctx, newTodo("t2", "u1", "b"))
	_ = s.Create(ctx, newTodo("t3", "u2", "c"))
	list, err := s.ListByUser(ctx, "u1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 todos for u1, got %d", len(list))
	}
}

func TestTodoStore_Update(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	_ = s.Create(ctx, newTodo("t1", "u1", "old title"))
	td, _ := s.GetByID(ctx, "t1")
	td.Title = "new title"
	if err := s.Update(ctx, td); err != nil {
		t.Fatal(err)
	}
	got, _ := s.GetByID(ctx, "t1")
	if got.Title != "new title" {
		t.Fatalf("expected new title, got %s", got.Title)
	}
}

func TestTodoStore_UpdateNotFound(t *testing.T) {
	s := store.NewMemoryTodoStore()
	err := s.Update(context.Background(), newTodo("missing", "u1", "x"))
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTodoStore_Delete(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	_ = s.Create(ctx, newTodo("t1", "u1", "task"))
	if err := s.Delete(ctx, "t1"); err != nil {
		t.Fatal(err)
	}
	_, err := s.GetByID(ctx, "t1")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestTodoStore_DeleteNotFound(t *testing.T) {
	s := store.NewMemoryTodoStore()
	err := s.Delete(context.Background(), "missing")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTodoStore_ConcurrentReadWrite(t *testing.T) {
	s := store.NewMemoryTodoStore()
	ctx := context.Background()
	_ = s.Create(ctx, newTodo("t1", "u1", "task"))
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.GetByID(ctx, "t1")
			_, _ = s.ListByUser(ctx, "u1")
		}()
	}
	wg.Wait()
}
