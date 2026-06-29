package store

import (
	"context"
	"sync"

	"github.com/user/todo/internal/model"
)

type MemoryUserStore struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{users: make(map[string]*model.User)}
}

func (s *MemoryUserStore) Create(_ context.Context, u *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.users {
		if existing.Username == u.Username || existing.Email == u.Email {
			return ErrConflict
		}
	}
	copy := *u
	s.users[u.ID] = &copy
	return nil
}

func (s *MemoryUserStore) GetByID(_ context.Context, id string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	copy := *u
	return &copy, nil
}

func (s *MemoryUserStore) List(_ context.Context) ([]*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.User, 0, len(s.users))
	for _, u := range s.users {
		copy := *u
		result = append(result, &copy)
	}
	return result, nil
}

func (s *MemoryUserStore) Update(_ context.Context, u *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[u.ID]; !ok {
		return ErrNotFound
	}
	copy := *u
	s.users[u.ID] = &copy
	return nil
}

func (s *MemoryUserStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return ErrNotFound
	}
	delete(s.users, id)
	return nil
}

type MemoryTodoStore struct {
	mu    sync.RWMutex
	todos map[string]*model.Todo
}

func NewMemoryTodoStore() *MemoryTodoStore {
	return &MemoryTodoStore{todos: make(map[string]*model.Todo)}
}

func (s *MemoryTodoStore) Create(_ context.Context, t *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := *t
	s.todos[t.ID] = &copy
	return nil
}

func (s *MemoryTodoStore) GetByID(_ context.Context, id string) (*model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.todos[id]
	if !ok {
		return nil, ErrNotFound
	}
	copy := *t
	return &copy, nil
}

func (s *MemoryTodoStore) ListByUser(_ context.Context, userID string) ([]*model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*model.Todo
	for _, t := range s.todos {
		if t.UserID == userID {
			copy := *t
			result = append(result, &copy)
		}
	}
	return result, nil
}

func (s *MemoryTodoStore) Update(_ context.Context, t *model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[t.ID]; !ok {
		return ErrNotFound
	}
	copy := *t
	s.todos[t.ID] = &copy
	return nil
}

func (s *MemoryTodoStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return nil
}
