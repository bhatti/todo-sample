package api_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func createUser(t *testing.T, base string) string {
	t.Helper()
	resp, err := http.Post(base+"/users", "application/json", mustJSON(t, map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var u map[string]any
	json.NewDecoder(resp.Body).Decode(&u)
	id, _ := u["ID"].(string)
	return id
}

func TestTodos_CRUD(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	base := srv.URL

	userID := createUser(t, base)

	// Create todo
	resp, err := http.Post(base+"/users/"+userID+"/todos", "application/json", mustJSON(t, map[string]any{
		"title":       "buy milk",
		"description": "full fat",
		"priority":    1,
		"status":      0,
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create todo: expected 201, got %d", resp.StatusCode)
	}
	var created map[string]any
	json.NewDecoder(resp.Body).Decode(&created)
	todoID, _ := created["ID"].(string)
	if todoID == "" {
		t.Fatal("expected non-empty todo ID")
	}

	// List todos
	resp2, _ := http.Get(base + "/users/" + userID + "/todos")
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("list todos: expected 200, got %d", resp2.StatusCode)
	}
	var list []any
	json.NewDecoder(resp2.Body).Decode(&list)
	if len(list) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(list))
	}

	// Get by ID
	resp3, _ := http.Get(base + "/users/" + userID + "/todos/" + todoID)
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("get todo: expected 200, got %d", resp3.StatusCode)
	}

	// Update todo
	req4, _ := http.NewRequest(http.MethodPut, base+"/users/"+userID+"/todos/"+todoID, mustJSON(t, map[string]any{
		"title":    "updated title",
		"priority": 2,
		"status":   1,
	}))
	req4.Header.Set("Content-Type", "application/json")
	resp4, _ := http.DefaultClient.Do(req4)
	defer resp4.Body.Close()
	if resp4.StatusCode != http.StatusOK {
		t.Fatalf("update todo: expected 200, got %d", resp4.StatusCode)
	}

	// Delete todo
	req5, _ := http.NewRequest(http.MethodDelete, base+"/users/"+userID+"/todos/"+todoID, nil)
	resp5, _ := http.DefaultClient.Do(req5)
	defer resp5.Body.Close()
	if resp5.StatusCode != http.StatusNoContent {
		t.Fatalf("delete todo: expected 204, got %d", resp5.StatusCode)
	}

	// Get after delete -> 404
	resp6, _ := http.Get(base + "/users/" + userID + "/todos/" + todoID)
	defer resp6.Body.Close()
	if resp6.StatusCode != http.StatusNotFound {
		t.Fatalf("get after delete: expected 404, got %d", resp6.StatusCode)
	}
}

func TestTodos_UserNotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	base := srv.URL

	// Create todo for non-existent user
	resp, _ := http.Post(base+"/users/nonexistent/todos", "application/json", mustJSON(t, map[string]any{
		"title": "task",
	}))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for missing user, got %d", resp.StatusCode)
	}

	// List todos for non-existent user
	resp2, _ := http.Get(base + "/users/nonexistent/todos")
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for missing user in list, got %d", resp2.StatusCode)
	}
}

func TestTodos_NotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	base := srv.URL
	userID := createUser(t, base)

	resp, _ := http.Get(base + "/users/" + userID + "/todos/nonexistent")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get nonexistent todo: expected 404, got %d", resp.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodPut, base+"/users/"+userID+"/todos/nonexistent", mustJSON(t, map[string]any{
		"title": "x",
	}))
	req.Header.Set("Content-Type", "application/json")
	resp2, _ := http.DefaultClient.Do(req)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusNotFound {
		t.Fatalf("update nonexistent todo: expected 404, got %d", resp2.StatusCode)
	}

	req3, _ := http.NewRequest(http.MethodDelete, base+"/users/"+userID+"/todos/nonexistent", nil)
	resp3, _ := http.DefaultClient.Do(req3)
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusNotFound {
		t.Fatalf("delete nonexistent todo: expected 404, got %d", resp3.StatusCode)
	}
}

func TestTodos_ValidationErrors(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	base := srv.URL
	userID := createUser(t, base)

	// Missing title
	resp, _ := http.Post(base+"/users/"+userID+"/todos", "application/json", mustJSON(t, map[string]any{
		"title": "",
	}))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("empty title: expected 400, got %d", resp.StatusCode)
	}
}

func TestTodos_EmptyList(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	base := srv.URL
	userID := createUser(t, base)

	resp, _ := http.Get(base + "/users/" + userID + "/todos")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var list []any
	json.NewDecoder(resp.Body).Decode(&list)
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d", len(list))
	}
}
