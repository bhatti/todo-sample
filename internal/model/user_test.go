package model

import (
	"testing"
	"time"
)

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
