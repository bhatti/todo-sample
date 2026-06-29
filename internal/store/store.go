package store

import (
	"context"
	"errors"

	"github.com/user/todo/internal/model"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type UserStore interface {
	Create(ctx context.Context, u *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id string) error
}

type TodoStore interface {
	Create(ctx context.Context, t *model.Todo) error
	GetByID(ctx context.Context, id string) (*model.Todo, error)
	ListByUser(ctx context.Context, userID string) ([]*model.Todo, error)
	Update(ctx context.Context, t *model.Todo) error
	Delete(ctx context.Context, id string) error
}
