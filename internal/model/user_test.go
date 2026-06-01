package model

import (
	"testing"
	"time"
)

func TestUserValidate(t *testing.T) {
	now := time.Now()
	base := User{
		ID:        "01J1Z0000000000000000000",
		Username:  "alice",
		Email:     "alice@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := base.Validate(); err != nil {
		t.Errorf("expected no error for valid user, got: %v", err)
	}

	noID := base
	noID.ID = ""
	if err := noID.Validate(); err == nil {
		t.Error("expected error when ID is empty")
	}

	noUsername := base
	noUsername.Username = ""
	if err := noUsername.Validate(); err == nil {
		t.Error("expected error when Username is empty")
	}

	noEmail := base
	noEmail.Email = ""
	if err := noEmail.Validate(); err == nil {
		t.Error("expected error when Email is empty")
	}
}

func TestUserFields(t *testing.T) {
	now := time.Now()
	u := User{
		ID:        "01J1Z0000000000000000000",
		Username:  "alice",
		Email:     "alice@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if u.ID != "01J1Z0000000000000000000" {
		t.Errorf("unexpected ID: %s", u.ID)
	}
	if u.Username != "alice" {
		t.Errorf("unexpected Username: %s", u.Username)
	}
	if u.Email != "alice@example.com" {
		t.Errorf("unexpected Email: %s", u.Email)
	}
	if !u.CreatedAt.Equal(now) {
		t.Errorf("unexpected CreatedAt")
	}
	if !u.UpdatedAt.Equal(now) {
		t.Errorf("unexpected UpdatedAt")
	}
}
