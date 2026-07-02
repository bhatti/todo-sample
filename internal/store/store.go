// Package store defines the storage interfaces for the todo application.
package store

import "github.com/user/todo/internal/model"

// UserStore defines CRUD operations for users.
type UserStore interface {
	Create(u *model.User) (*model.User, error)
	GetByID(id string) (*model.User, error)
	List() ([]*model.User, error)
	Update(u *model.User) (*model.User, error)
	Delete(id string) error
}

// TodoStore defines CRUD operations for todos.
type TodoStore interface {
	Create(t *model.Todo) (*model.Todo, error)
	GetByID(id string) (*model.Todo, error)
	ListByUser(userID string) ([]*model.Todo, error)
	Update(t *model.Todo) (*model.Todo, error)
	Delete(id string) error
}
