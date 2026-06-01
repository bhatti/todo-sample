package todo

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTodoJSONRoundTrip(t *testing.T) {
	original := Todo{
		ID:          "abc123",
		Title:       "Buy groceries",
		Description: "Milk, eggs, bread",
		Completed:   false,
		CreatedAt:   time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Todo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID: got %q, want %q", decoded.ID, original.ID)
	}
	if decoded.Title != original.Title {
		t.Errorf("Title: got %q, want %q", decoded.Title, original.Title)
	}
	if decoded.Description != original.Description {
		t.Errorf("Description: got %q, want %q", decoded.Description, original.Description)
	}
	if decoded.Completed != original.Completed {
		t.Errorf("Completed: got %v, want %v", decoded.Completed, original.Completed)
	}
	if !decoded.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("CreatedAt: got %v, want %v", decoded.CreatedAt, original.CreatedAt)
	}
}

func TestTodoJSONKeys(t *testing.T) {
	todo := Todo{ID: "1", Title: "T", Description: "D", Completed: true, CreatedAt: time.Now()}

	data, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	for _, key := range []string{"id", "title", "description", "completed", "created_at"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("missing expected JSON key %q", key)
		}
	}
}
