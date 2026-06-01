// Package model defines the core data structures for the todo application.
package model

import (
	"errors"
	"time"
)

// User represents an account that can manage todo tasks.
type User struct {
	// ID is a ULID string uniquely identifying the user.
	ID string
	// Username is the user's display name.
	Username string
	// Email is the user's email address.
	Email string
	// CreatedAt is when the user was created.
	CreatedAt time.Time
	// UpdatedAt is when the user was last modified.
	UpdatedAt time.Time
}

// Validate returns an error if the User is missing required fields.
func (u *User) Validate() error {
	if u.ID == "" {
		return errors.New("user ID is required")
	}
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	return nil
}
