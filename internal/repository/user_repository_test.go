package repository

import (
	"errors"
	"sync"
	"testing"

	"github.com/todo-app/todo/internal/model"
)

func newUser(id string) *model.User {
	return model.NewUser(id, "Name "+id, id+"@example.com")
}

func TestUserRepository_CreateAndGet(t *testing.T) {
	r := NewInMemoryUserRepository()
	u := newUser("u1")
	if err := r.Create(u); err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := r.GetByID("u1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.ID != "u1" {
		t.Errorf("expected ID u1, got %s", got.ID)
	}
}

func TestUserRepository_Create_AlreadyExists(t *testing.T) {
	r := NewInMemoryUserRepository()
	u := newUser("u1")
	_ = r.Create(u)
	if err := r.Create(u); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	r := NewInMemoryUserRepository()
	if _, err := r.GetByID("missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserRepository_List(t *testing.T) {
	r := NewInMemoryUserRepository()
	_ = r.Create(newUser("u1"))
	_ = r.Create(newUser("u2"))
	list, err := r.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 users, got %d", len(list))
	}
}

func TestUserRepository_Update(t *testing.T) {
	r := NewInMemoryUserRepository()
	u := newUser("u1")
	_ = r.Create(u)
	u.Name = "Updated"
	if err := r.Update(u); err != nil {
		t.Fatalf("Update: %v", err)
	}
	got, _ := r.GetByID("u1")
	if got.Name != "Updated" {
		t.Errorf("expected Updated, got %s", got.Name)
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	r := NewInMemoryUserRepository()
	u := newUser("u1")
	if err := r.Update(u); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	r := NewInMemoryUserRepository()
	_ = r.Create(newUser("u1"))
	if err := r.Delete("u1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := r.GetByID("u1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	r := NewInMemoryUserRepository()
	if err := r.Delete("missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserRepository_ConcurrentCreateAndList(t *testing.T) {
	r := NewInMemoryUserRepository()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		id := "u" + string(rune('0'+i%10)) + string(rune('0'+i/10))
		go func(id string) {
			defer wg.Done()
			_ = r.Create(newUser(id))
			_, _ = r.List()
		}(id)
	}
	wg.Wait()
}
