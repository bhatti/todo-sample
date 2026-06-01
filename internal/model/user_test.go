package model

import (
	"testing"
)

func TestNewUser_FieldsPopulated(t *testing.T) {
	u := NewUser("u1", "Alice", "alice@example.com")
	if u.ID != "u1" {
		t.Errorf("expected ID u1, got %s", u.ID)
	}
	if u.Name != "Alice" {
		t.Errorf("expected Name Alice, got %s", u.Name)
	}
	if u.Email != "alice@example.com" {
		t.Errorf("expected Email alice@example.com, got %s", u.Email)
	}
	if u.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if u.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestUser_Validate_Valid(t *testing.T) {
	u := NewUser("u1", "Alice", "alice@example.com")
	if err := u.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUser_Validate_MissingID(t *testing.T) {
	u := NewUser("", "Alice", "alice@example.com")
	if err := u.Validate(); err == nil {
		t.Error("expected error for empty ID")
	}
}

func TestUser_Validate_MissingName(t *testing.T) {
	u := NewUser("u1", "", "alice@example.com")
	if err := u.Validate(); err == nil {
		t.Error("expected error for empty Name")
	}
}

func TestUser_Validate_MissingEmail(t *testing.T) {
	u := NewUser("u1", "Alice", "")
	if err := u.Validate(); err == nil {
		t.Error("expected error for empty Email")
	}
}
