// Package store provides in-memory storage implementations for the todo application.
package store

import (
	"context"
	"errors"
	"sync"

	"github.com/user/todo/internal/model"
)

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("not found")

// UserStore defines CRUD operations for users.
type UserStore interface {
	Create(ctx context.Context, u *model.User) error
	Get(ctx context.Context, id string) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id string) error
}

// MemoryUserStore is an in-memory implementation of UserStore.
type MemoryUserStore struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

// NewMemoryUserStore returns a new MemoryUserStore.
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users: make(map[string]*model.User),
	}
}

// Create adds a new user to the store.
func (s *MemoryUserStore) Create(ctx context.Context, u *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := *u
	s.users[u.ID] = &cp
	return nil
}

// Get retrieves a user by ID.
func (s *MemoryUserStore) Get(ctx context.Context, id string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *u
	return &cp, nil
}

// List returns all users.
func (s *MemoryUserStore) List(ctx context.Context) ([]*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.User, 0, len(s.users))
	for _, u := range s.users {
		cp := *u
		result = append(result, &cp)
	}
	return result, nil
}

// Update replaces an existing user.
func (s *MemoryUserStore) Update(ctx context.Context, u *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[u.ID]; !ok {
		return ErrNotFound
	}
	cp := *u
	s.users[u.ID] = &cp
	return nil
}

// Delete removes a user by ID.
func (s *MemoryUserStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return ErrNotFound
	}
	delete(s.users, id)
	return nil
}
