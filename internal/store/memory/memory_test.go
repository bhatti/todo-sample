package memory_test

import (
	"testing"

	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
	"github.com/user/todo/internal/store/memory"
)

func TestCreateAndGetUser(t *testing.T) {
	s := memory.New()
	u := model.User{ID: "01", Username: "alice", Email: "alice@example.com"}
	if err := s.CreateUser(u); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	got, err := s.GetUser("01")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.ID != u.ID || got.Username != u.Username || got.Email != u.Email {
		t.Errorf("got %+v, want %+v", got, u)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	s := memory.New()
	_, err := s.GetUser("missing")
	if err != store.ErrNotFound {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestListUsers_Empty(t *testing.T) {
	s := memory.New()
	users, err := s.ListUsers()
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("want empty list, got %d items", len(users))
	}
}

func TestDeleteUser(t *testing.T) {
	s := memory.New()
	u := model.User{ID: "01", Username: "alice", Email: "alice@example.com"}
	_ = s.CreateUser(u)
	if err := s.DeleteUser("01"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	_, err := s.GetUser("01")
	if err != store.ErrNotFound {
		t.Errorf("want ErrNotFound after delete, got %v", err)
	}
}

func TestCreateTodo(t *testing.T) {
	s := memory.New()
	todo := model.Todo{ID: "t1", UserID: "u1", Title: "Test"}
	if err := s.CreateTodo(todo); err != nil {
		t.Fatalf("CreateTodo: %v", err)
	}
	got, err := s.GetTodo("u1", "t1")
	if err != nil {
		t.Fatalf("GetTodo: %v", err)
	}
	if got.ID != todo.ID || got.UserID != todo.UserID {
		t.Errorf("got %+v, want %+v", got, todo)
	}
}

func TestListTodos_ForUser(t *testing.T) {
	s := memory.New()
	_ = s.CreateTodo(model.Todo{ID: "t1", UserID: "u1", Title: "A"})
	_ = s.CreateTodo(model.Todo{ID: "t2", UserID: "u1", Title: "B"})
	_ = s.CreateTodo(model.Todo{ID: "t3", UserID: "u2", Title: "C"})

	todos, err := s.ListTodos("u1")
	if err != nil {
		t.Fatalf("ListTodos: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("want 2 todos for u1, got %d", len(todos))
	}
}

func TestUpdateTodo(t *testing.T) {
	s := memory.New()
	todo := model.Todo{ID: "t1", UserID: "u1", Title: "Test", Status: model.StatusPending}
	_ = s.CreateTodo(todo)

	todo.Status = model.StatusDone
	if err := s.UpdateTodo(todo); err != nil {
		t.Fatalf("UpdateTodo: %v", err)
	}
	got, _ := s.GetTodo("u1", "t1")
	if got.Status != model.StatusDone {
		t.Errorf("want StatusDone, got %v", got.Status)
	}
}

func TestDeleteTodo(t *testing.T) {
	s := memory.New()
	todo := model.Todo{ID: "t1", UserID: "u1", Title: "Test"}
	_ = s.CreateTodo(todo)
	if err := s.DeleteTodo("u1", "t1"); err != nil {
		t.Fatalf("DeleteTodo: %v", err)
	}
	_, err := s.GetTodo("u1", "t1")
	if err != store.ErrNotFound {
		t.Errorf("want ErrNotFound after delete, got %v", err)
	}
}
