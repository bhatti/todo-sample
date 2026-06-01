package model

import (
	"errors"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// User is the owner entity for todo tasks.
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a validated User with a generated ID.
func NewUser(name, email string) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name must not be empty")
	}
	if strings.TrimSpace(email) == "" || !strings.Contains(email, "@") {
		return nil, errors.New("email must be non-empty and contain '@'")
	}
	now := time.Now().UTC()
	return &User{
		ID:        ulid.Make().String(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
