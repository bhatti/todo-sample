package repository

import (
	"sync"

	"github.com/todo-app/todo/internal/model"
)

// UserRepository defines CRUD operations for User entities.
type UserRepository interface {
	Create(u *model.User) error
	GetByID(id string) (*model.User, error)
	List() ([]*model.User, error)
	Update(u *model.User) error
	Delete(id string) error
}

// InMemoryUserRepository is a thread-safe in-memory implementation of UserRepository.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	store map[string]*model.User
}

// NewInMemoryUserRepository returns an empty InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{store: make(map[string]*model.User)}
}

func (r *InMemoryUserRepository) Create(u *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[u.ID]; ok {
		return ErrAlreadyExists
	}
	copied := *u
	r.store[u.ID] = &copied
	return nil
}

func (r *InMemoryUserRepository) GetByID(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	copied := *u
	return &copied, nil
}

func (r *InMemoryUserRepository) List() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*model.User, 0, len(r.store))
	for _, u := range r.store {
		copied := *u
		out = append(out, &copied)
	}
	return out, nil
}

func (r *InMemoryUserRepository) Update(u *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[u.ID]; !ok {
		return ErrNotFound
	}
	copied := *u
	r.store[u.ID] = &copied
	return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}
