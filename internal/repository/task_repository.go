package repository

import (
	"sync"

	"github.com/todo-app/todo/internal/model"
)

// TodoTaskRepository defines CRUD operations for TodoTask entities.
type TodoTaskRepository interface {
	Create(t *model.TodoTask) error
	GetByID(id string) (*model.TodoTask, error)
	List() ([]*model.TodoTask, error)
	ListByUserID(userID string) ([]*model.TodoTask, error)
	Update(t *model.TodoTask) error
	Delete(id string) error
}

// InMemoryTodoTaskRepository is a thread-safe in-memory implementation of TodoTaskRepository.
type InMemoryTodoTaskRepository struct {
	mu    sync.RWMutex
	store map[string]*model.TodoTask
}

// NewInMemoryTodoTaskRepository returns an empty InMemoryTodoTaskRepository.
func NewInMemoryTodoTaskRepository() *InMemoryTodoTaskRepository {
	return &InMemoryTodoTaskRepository{store: make(map[string]*model.TodoTask)}
}

func (r *InMemoryTodoTaskRepository) Create(t *model.TodoTask) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[t.ID]; ok {
		return ErrAlreadyExists
	}
	copied := *t
	r.store[t.ID] = &copied
	return nil
}

func (r *InMemoryTodoTaskRepository) GetByID(id string) (*model.TodoTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	copied := *t
	return &copied, nil
}

func (r *InMemoryTodoTaskRepository) List() ([]*model.TodoTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*model.TodoTask, 0, len(r.store))
	for _, t := range r.store {
		copied := *t
		out = append(out, &copied)
	}
	return out, nil
}

func (r *InMemoryTodoTaskRepository) ListByUserID(userID string) ([]*model.TodoTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*model.TodoTask
	for _, t := range r.store {
		if t.UserID == userID {
			copied := *t
			out = append(out, &copied)
		}
	}
	return out, nil
}

func (r *InMemoryTodoTaskRepository) Update(t *model.TodoTask) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[t.ID]; !ok {
		return ErrNotFound
	}
	copied := *t
	r.store[t.ID] = &copied
	return nil
}

func (r *InMemoryTodoTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}
