package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/todo/internal/model"
)

// MemoryTodoStore is a thread-safe in-memory implementation of TodoStore.
type MemoryTodoStore struct {
	mu    sync.RWMutex
	todos map[string]*model.Todo
}

// NewMemoryTodoStore returns an initialised MemoryTodoStore.
func NewMemoryTodoStore() *MemoryTodoStore {
	return &MemoryTodoStore{
		todos: make(map[string]*model.Todo),
	}
}

// Create stores a new todo, assigning ID, default Status, and timestamps.
func (s *MemoryTodoStore) Create(t *model.Todo) (*model.Todo, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	created := &model.Todo{
		ID:          id,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    t.Priority,
		Status:      model.StatusPending,
		DueAt:       t.DueAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.mu.Lock()
	s.todos[id] = created
	s.mu.Unlock()
	copy := *created
	return &copy, nil
}

// GetByID retrieves a todo by its ID.
func (s *MemoryTodoStore) GetByID(id string) (*model.Todo, error) {
	s.mu.RLock()
	t, ok := s.todos[id]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("todo %q not found", id)
	}
	copy := *t
	return &copy, nil
}

// ListByUser returns all todos belonging to the given user.
func (s *MemoryTodoStore) ListByUser(userID string) ([]*model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*model.Todo
	for _, t := range s.todos {
		if t.UserID == userID {
			copy := *t
			result = append(result, &copy)
		}
	}
	if result == nil {
		result = []*model.Todo{}
	}
	return result, nil
}

// Update replaces the stored todo with the provided values.
func (s *MemoryTodoStore) Update(t *model.Todo) (*model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.todos[t.ID]
	if !ok {
		return nil, fmt.Errorf("todo %q not found", t.ID)
	}
	existing.Title = t.Title
	existing.Description = t.Description
	existing.Priority = t.Priority
	existing.Status = t.Status
	existing.DueAt = t.DueAt
	existing.UpdatedAt = time.Now().UTC()
	copy := *existing
	return &copy, nil
}

// Delete removes a todo by ID.
func (s *MemoryTodoStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return fmt.Errorf("todo %q not found", id)
	}
	delete(s.todos, id)
	return nil
}
