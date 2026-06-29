package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/todo/internal/api"
	"github.com/user/todo/internal/store"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	us := store.NewMemoryUserStore()
	ts := store.NewMemoryTodoStore()
	return httptest.NewServer(api.NewRouter(us, ts))
}

func mustJSON(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewBuffer(b)
}

func TestUsers_CRUD(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	base := srv.URL

	// Create
	resp, err := http.Post(base+"/users", "application/json", mustJSON(t, map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", resp.StatusCode)
	}
	var created map[string]any
	json.NewDecoder(resp.Body).Decode(&created)
	id, _ := created["ID"].(string)
	if id == "" {
		t.Fatal("expected non-empty ID")
	}

	// Get by ID
	resp2, _ := http.Get(base + "/users/" + id)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", resp2.StatusCode)
	}

	// List
	resp3, _ := http.Get(base + "/users")
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", resp3.StatusCode)
	}
	var list []any
	json.NewDecoder(resp3.Body).Decode(&list)
	if len(list) != 1 {
		t.Fatalf("list: expected 1 user, got %d", len(list))
	}

	// Update
	req, _ := http.NewRequest(http.MethodPut, base+"/users/"+id, mustJSON(t, map[string]string{
		"username": "alice2",
		"email":    "alice2@example.com",
	}))
	req.Header.Set("Content-Type", "application/json")
	resp4, _ := http.DefaultClient.Do(req)
	defer resp4.Body.Close()
	if resp4.StatusCode != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", resp4.StatusCode)
	}

	// Delete
	req5, _ := http.NewRequest(http.MethodDelete, base+"/users/"+id, nil)
	resp5, _ := http.DefaultClient.Do(req5)
	defer resp5.Body.Close()
	if resp5.StatusCode != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", resp5.StatusCode)
	}

	// Get after delete -> 404
	resp6, _ := http.Get(base + "/users/" + id)
	defer resp6.Body.Close()
	if resp6.StatusCode != http.StatusNotFound {
		t.Fatalf("get after delete: expected 404, got %d", resp6.StatusCode)
	}
}

func TestUsers_NotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	base := srv.URL

	resp, _ := http.Get(base + "/users/nonexistent")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodPut, base+"/users/nonexistent", mustJSON(t, map[string]string{
		"username": "x", "email": "x@e.com",
	}))
	req.Header.Set("Content-Type", "application/json")
	resp2, _ := http.DefaultClient.Do(req)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusNotFound {
		t.Fatalf("put nonexistent: expected 404, got %d", resp2.StatusCode)
	}

	req3, _ := http.NewRequest(http.MethodDelete, base+"/users/nonexistent", nil)
	resp3, _ := http.DefaultClient.Do(req3)
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusNotFound {
		t.Fatalf("delete nonexistent: expected 404, got %d", resp3.StatusCode)
	}
}

func TestUsers_ValidationErrors(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	base := srv.URL

	// Missing username
	resp, _ := http.Post(base+"/users", "application/json", mustJSON(t, map[string]string{
		"username": "", "email": "a@e.com",
	}))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("empty username: expected 400, got %d", resp.StatusCode)
	}

	// Missing email
	resp2, _ := http.Post(base+"/users", "application/json", mustJSON(t, map[string]string{
		"username": "alice", "email": "",
	}))
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusBadRequest {
		t.Fatalf("empty email: expected 400, got %d", resp2.StatusCode)
	}
}

func TestUsers_EmptyList(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()
	resp, _ := http.Get(srv.URL + "/users")
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
