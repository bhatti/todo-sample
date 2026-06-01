package model

import (
	"errors"
	"strings"
	"time"
)

// User represents an account that owns todo tasks.
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a User with the given identity fields and timestamps set to now.
func NewUser(id, name, email string) *User {
	now := time.Now().UTC()
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate returns an error if required fields are missing or malformed.
func (u *User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("user: ID is required")
	}
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("user: Name is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("user: Email is required")
	}
	return nil
}
