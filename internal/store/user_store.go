package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/user/todo/internal/model"
)

// MemoryUserStore is a thread-safe in-memory implementation of UserStore.
type MemoryUserStore struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

// NewMemoryUserStore returns an initialised MemoryUserStore.
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users: make(map[string]*model.User),
	}
}

func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// Create stores a new user, assigning ID and timestamps.
func (s *MemoryUserStore) Create(u *model.User) (*model.User, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	created := &model.User{
		ID:        id,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.mu.Lock()
	s.users[id] = created
	s.mu.Unlock()
	copy := *created
	return &copy, nil
}

// GetByID retrieves a user by its ID.
func (s *MemoryUserStore) GetByID(id string) (*model.User, error) {
	s.mu.RLock()
	u, ok := s.users[id]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("user %q not found", id)
	}
	copy := *u
	return &copy, nil
}

// List returns all stored users.
func (s *MemoryUserStore) List() ([]*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.User, 0, len(s.users))
	for _, u := range s.users {
		copy := *u
		result = append(result, &copy)
	}
	return result, nil
}

// Update replaces the stored user with the provided values.
func (s *MemoryUserStore) Update(u *model.User) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.users[u.ID]
	if !ok {
		return nil, fmt.Errorf("user %q not found", u.ID)
	}
	existing.Username = u.Username
	existing.Email = u.Email
	existing.UpdatedAt = time.Now().UTC()
	copy := *existing
	return &copy, nil
}

// Delete removes a user by ID.
func (s *MemoryUserStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return fmt.Errorf("user %q not found", id)
	}
	delete(s.users, id)
	return nil
}
