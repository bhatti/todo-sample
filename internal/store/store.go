// Package store defines interfaces and sentinel errors for data persistence.
package store

import (
	"errors"

	"github.com/user/todo/internal/model"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// ErrAlreadyExists is returned when an entity with the same ID already exists.
var ErrAlreadyExists = errors.New("already exists")

// UserStore defines operations for managing users.
type UserStore interface {
	CreateUser(u model.User) error
	GetUser(id string) (model.User, error)
	ListUsers() ([]model.User, error)
	DeleteUser(id string) error
}

// TodoStore defines operations for managing todos.
type TodoStore interface {
	CreateTodo(t model.Todo) error
	GetTodo(userID, id string) (model.Todo, error)
	ListTodos(userID string) ([]model.Todo, error)
	UpdateTodo(t model.Todo) error
	DeleteTodo(userID, id string) error
}
