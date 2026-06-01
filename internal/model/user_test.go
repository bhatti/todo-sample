package model

import (
	"testing"
)

func TestNewUser_Valid(t *testing.T) {
	u, err := NewUser("Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == "" {
		t.Error("ID must not be empty")
	}
	if u.Name != "Alice" {
		t.Errorf("Name = %q, want %q", u.Name, "Alice")
	}
	if u.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", u.Email, "alice@example.com")
	}
	if u.CreatedAt.IsZero() {
		t.Error("CreatedAt must not be zero")
	}
	if u.UpdatedAt.IsZero() {
		t.Error("UpdatedAt must not be zero")
	}
}

func TestNewUser_EmptyName(t *testing.T) {
	_, err := NewUser("", "alice@example.com")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewUser_InvalidEmail(t *testing.T) {
	cases := []string{"", "bademail", "   "}
	for _, email := range cases {
		_, err := NewUser("Alice", email)
		if err == nil {
			t.Errorf("expected error for email %q", email)
		}
	}
}
