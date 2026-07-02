// Package memory provides an in-memory implementation of the store interfaces.
package memory

import (
	"sync"

	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
)

// Store is a thread-safe in-memory implementation of UserStore and TodoStore.
type Store struct {
	mu    sync.RWMutex
	users map[string]model.User
	todos map[string]map[string]model.Todo // userID -> todoID -> Todo
}

// New returns a new empty Store.
func New() *Store {
	return &Store{
		users: make(map[string]model.User),
		todos: make(map[string]map[string]model.Todo),
	}
}

// CreateUser stores a new user. Returns ErrAlreadyExists if the ID is taken.
func (s *Store) CreateUser(u model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[u.ID]; ok {
		return store.ErrAlreadyExists
	}
	s.users[u.ID] = u
	return nil
}

// GetUser retrieves a user by ID. Returns ErrNotFound if missing.
func (s *Store) GetUser(id string) (model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return model.User{}, store.ErrNotFound
	}
	return u, nil
}

// ListUsers returns all stored users.
func (s *Store) ListUsers() ([]model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]model.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users, nil
}

// DeleteUser removes a user by ID. Returns ErrNotFound if missing.
func (s *Store) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return store.ErrNotFound
	}
	delete(s.users, id)
	delete(s.todos, id)
	return nil
}

// CreateTodo stores a new todo. Returns ErrAlreadyExists if the ID is taken.
func (s *Store) CreateTodo(t model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[t.UserID]; !ok {
		s.todos[t.UserID] = make(map[string]model.Todo)
	}
	if _, ok := s.todos[t.UserID][t.ID]; ok {
		return store.ErrAlreadyExists
	}
	s.todos[t.UserID][t.ID] = t
	return nil
}

// GetTodo retrieves a todo by userID and todoID. Returns ErrNotFound if missing.
func (s *Store) GetTodo(userID, id string) (model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userTodos, ok := s.todos[userID]
	if !ok {
		return model.Todo{}, store.ErrNotFound
	}
	t, ok := userTodos[id]
	if !ok {
		return model.Todo{}, store.ErrNotFound
	}
	return t, nil
}

// ListTodos returns all todos belonging to the given user.
func (s *Store) ListTodos(userID string) ([]model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userTodos, ok := s.todos[userID]
	if !ok {
		return []model.Todo{}, nil
	}
	todos := make([]model.Todo, 0, len(userTodos))
	for _, t := range userTodos {
		todos = append(todos, t)
	}
	return todos, nil
}

// UpdateTodo replaces an existing todo. Returns ErrNotFound if missing.
func (s *Store) UpdateTodo(t model.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	userTodos, ok := s.todos[t.UserID]
	if !ok {
		return store.ErrNotFound
	}
	if _, ok := userTodos[t.ID]; !ok {
		return store.ErrNotFound
	}
	s.todos[t.UserID][t.ID] = t
	return nil
}

// DeleteTodo removes a todo by userID and todoID. Returns ErrNotFound if missing.
func (s *Store) DeleteTodo(userID, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	userTodos, ok := s.todos[userID]
	if !ok {
		return store.ErrNotFound
	}
	if _, ok := userTodos[id]; !ok {
		return store.ErrNotFound
	}
	delete(s.todos[userID], id)
	return nil
}
