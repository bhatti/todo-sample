package repository

import (
	"errors"
	"sync"
	"testing"

	"github.com/todo-app/todo/internal/model"
)

func newTask(id, userID string) *model.TodoTask {
	return model.NewTodoTask(id, userID, "Task "+id)
}

func TestTaskRepository_CreateAndGet(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	task := newTask("t1", "u1")
	if err := r.Create(task); err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := r.GetByID("t1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.ID != "t1" {
		t.Errorf("expected ID t1, got %s", got.ID)
	}
}

func TestTaskRepository_Create_AlreadyExists(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	task := newTask("t1", "u1")
	_ = r.Create(task)
	if err := r.Create(task); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	if _, err := r.GetByID("missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestTaskRepository_List(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	_ = r.Create(newTask("t1", "u1"))
	_ = r.Create(newTask("t2", "u1"))
	list, err := r.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(list))
	}
}

func TestTaskRepository_ListByUserID(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	_ = r.Create(newTask("t1", "u1"))
	_ = r.Create(newTask("t2", "u1"))
	_ = r.Create(newTask("t3", "u2"))

	u1Tasks, err := r.ListByUserID("u1")
	if err != nil {
		t.Fatalf("ListByUserID: %v", err)
	}
	if len(u1Tasks) != 2 {
		t.Errorf("expected 2 tasks for u1, got %d", len(u1Tasks))
	}

	u2Tasks, _ := r.ListByUserID("u2")
	if len(u2Tasks) != 1 {
		t.Errorf("expected 1 task for u2, got %d", len(u2Tasks))
	}

	noTasks, _ := r.ListByUserID("u3")
	if len(noTasks) != 0 {
		t.Errorf("expected 0 tasks for u3, got %d", len(noTasks))
	}
}

func TestTaskRepository_Update(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	task := newTask("t1", "u1")
	_ = r.Create(task)
	task.Title = "Updated"
	task.Priority = model.PriorityHigh
	if err := r.Update(task); err != nil {
		t.Fatalf("Update: %v", err)
	}
	got, _ := r.GetByID("t1")
	if got.Title != "Updated" {
		t.Errorf("expected Updated, got %s", got.Title)
	}
	if got.Priority != model.PriorityHigh {
		t.Errorf("expected PriorityHigh, got %v", got.Priority)
	}
}

func TestTaskRepository_Update_NotFound(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	task := newTask("t1", "u1")
	if err := r.Update(task); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestTaskRepository_Delete(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	_ = r.Create(newTask("t1", "u1"))
	if err := r.Delete("t1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := r.GetByID("t1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestTaskRepository_Delete_NotFound(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	if err := r.Delete("missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestTaskRepository_ConcurrentCreateAndList(t *testing.T) {
	r := NewInMemoryTodoTaskRepository()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		id := "t" + string(rune('0'+i%10)) + string(rune('0'+i/10))
		go func(id string) {
			defer wg.Done()
			_ = r.Create(newTask(id, "u1"))
			_, _ = r.List()
		}(id)
	}
	wg.Wait()
}
