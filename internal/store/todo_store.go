package store

import (
	"context"
	"sync"

	"github.com/user/todo/internal/model"
)

// TodoStore defines CRUD operations for todos.
type TodoStore interface {
	Create(ctx context.Context, t *model.Todo) error
	Get(ctx context.Context, id string) (*model.Todo, error)
	ListByUser(ctx context.Context, userID string) ([]*model.Todo, error)
	Update(ctx context.Context, t *model.Todo) error
	Delete(ctx context.Context, id string) error
}

// MemoryTodoStore is an in-memory implementation of TodoStore.
type MemoryTodoStore struct {
	mu    sync.RWMutex
	todos map[string]*model.Todo
}

// NewMemoryTodoStore returns a new MemoryTodoStore.
func NewMemoryTodoStore() *MemoryTodoStore {
	return &MemoryTodoStore{
		todos: make(map[string]*model.Todo),
	}
}

// Create adds a new todo to the store.
func (s *MemoryTodoStore) Create(ctx context.Context, t *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := *t
	s.todos[t.ID] = &cp
	return nil
}

// Get retrieves a todo by ID.
func (s *MemoryTodoStore) Get(ctx context.Context, id string) (*model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.todos[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *t
	return &cp, nil
}

// ListByUser returns all todos belonging to a given user.
func (s *MemoryTodoStore) ListByUser(ctx context.Context, userID string) ([]*model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*model.Todo
	for _, t := range s.todos {
		if t.UserID == userID {
			cp := *t
			result = append(result, &cp)
		}
	}
	if result == nil {
		result = []*model.Todo{}
	}
	return result, nil
}

// Update replaces an existing todo.
func (s *MemoryTodoStore) Update(ctx context.Context, t *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[t.ID]; !ok {
		return ErrNotFound
	}
	cp := *t
	s.todos[t.ID] = &cp
	return nil
}

// Delete removes a todo by ID.
func (s *MemoryTodoStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return nil
}
