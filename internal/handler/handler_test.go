package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/todo/internal/handler"
	"github.com/user/todo/internal/router"
	"github.com/user/todo/internal/store/memory"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := memory.New()
	mux := router.New(handler.NewUserHandler(s), handler.NewTodoHandler(s, s))
	return httptest.NewServer(mux)
}

func postJSON(t *testing.T, url string, body any) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	return resp
}

func TestCreateUser_OK(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := postJSON(t, srv.URL+"/users", map[string]string{
		"username": "bob",
		"email":    "bob@x.com",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, got %d", resp.StatusCode)
	}
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["id"] == "" {
		t.Error("expected non-empty id in response")
	}
}

func TestCreateUser_MissingFields(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := postJSON(t, srv.URL+"/users", map[string]string{})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if _, ok := body["error"]; !ok {
		t.Error("expected error field in response")
	}
}

func TestGetUser_OK(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	// Create user
	resp := postJSON(t, srv.URL+"/users", map[string]string{
		"username": "alice",
		"email":    "alice@x.com",
	})
	defer resp.Body.Close()
	var created map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&created)
	id := created["id"].(string)

	// Get user
	r, err := http.Get(srv.URL + "/users/" + id)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", r.StatusCode)
	}
	var got map[string]any
	_ = json.NewDecoder(r.Body).Decode(&got)
	if got["id"] != id {
		t.Errorf("want id %q, got %q", id, got["id"])
	}
}

func TestGetUser_NotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	r, _ := http.Get(srv.URL + "/users/nonexistent")
	defer r.Body.Close()
	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404, got %d", r.StatusCode)
	}
}

func TestListUsers_Empty(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	r, _ := http.Get(srv.URL + "/users")
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", r.StatusCode)
	}
	var body []any
	_ = json.NewDecoder(r.Body).Decode(&body)
	if len(body) != 0 {
		t.Errorf("want empty list, got %v", body)
	}
}

func TestDeleteUser_OK(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := postJSON(t, srv.URL+"/users", map[string]string{
		"username": "alice",
		"email":    "alice@x.com",
	})
	defer resp.Body.Close()
	var created map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&created)
	id := created["id"].(string)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/users/"+id, nil)
	r, _ := http.DefaultClient.Do(req)
	defer r.Body.Close()
	if r.StatusCode != http.StatusNoContent {
		t.Fatalf("want 204, got %d", r.StatusCode)
	}
}

func createUser(t *testing.T, srv *httptest.Server) string {
	t.Helper()
	resp := postJSON(t, srv.URL+"/users", map[string]string{
		"username": "testuser",
		"email":    "test@x.com",
	})
	defer resp.Body.Close()
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return body["id"].(string)
}

func TestCreateTodo_OK(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	userID := createUser(t, srv)
	resp := postJSON(t, srv.URL+"/users/"+userID+"/todos", map[string]any{
		"title":    "Buy milk",
		"priority": 1,
	})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, got %d", resp.StatusCode)
	}
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["id"] == "" {
		t.Error("expected non-empty id")
	}
	if body["user_id"] == "" {
		t.Error("expected non-empty user_id")
	}
}

func TestCreateTodo_UserNotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp := postJSON(t, srv.URL+"/users/bad/todos", map[string]any{"title": "x"})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404, got %d", resp.StatusCode)
	}
}

func TestListTodos_OK(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	userID := createUser(t, srv)
	postJSON(t, srv.URL+"/users/"+userID+"/todos", map[string]any{"title": "Task 1"}).Body.Close()
	postJSON(t, srv.URL+"/users/"+userID+"/todos", map[string]any{"title": "Task 2"}).Body.Close()

	r, _ := http.Get(srv.URL + "/users/" + userID + "/todos")
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", r.StatusCode)
	}
	var body []any
	_ = json.NewDecoder(r.Body).Decode(&body)
	if len(body) != 2 {
		t.Errorf("want 2 todos, got %d", len(body))
	}
}

func TestGetTodo_OK(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	userID := createUser(t, srv)
	resp := postJSON(t, srv.URL+"/users/"+userID+"/todos", map[string]any{"title": "My Task"})
	defer resp.Body.Close()
	var created map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&created)
	todoID := created["id"].(string)

	r, _ := http.Get(srv.URL + "/users/" + userID + "/todos/" + todoID)
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", r.StatusCode)
	}
	var body map[string]any
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body["id"] != todoID {
		t.Errorf("want id %q, got %q", todoID, body["id"])
	}
}

func TestGetTodo_NotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	userID := createUser(t, srv)
	r, _ := http.Get(srv.URL + "/users/" + userID + "/todos/bad")
	defer r.Body.Close()
	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404, got %d", r.StatusCode)
	}
}

func TestUpdateTodo_OK(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	userID := createUser(t, srv)
	resp := postJSON(t, srv.URL+"/users/"+userID+"/todos", map[string]any{"title": "Task"})
	defer resp.Body.Close()
	var created map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&created)
	todoID := created["id"].(string)

	body, _ := json.Marshal(map[string]any{
		"title":  "Updated Task",
		"status": 2,
	})
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/users/"+userID+"/todos/"+todoID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r, _ := http.DefaultClient.Do(req)
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", r.StatusCode)
	}
	var updated map[string]any
	_ = json.NewDecoder(r.Body).Decode(&updated)
	if updated["status"] != float64(2) {
		t.Errorf("want status 2, got %v", updated["status"])
	}
}

func TestDeleteTodo_OK(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	userID := createUser(t, srv)
	resp := postJSON(t, srv.URL+"/users/"+userID+"/todos", map[string]any{"title": "Task"})
	defer resp.Body.Close()
	var created map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&created)
	todoID := created["id"].(string)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/users/"+userID+"/todos/"+todoID, nil)
	r, _ := http.DefaultClient.Do(req)
	defer r.Body.Close()
	if r.StatusCode != http.StatusNoContent {
		t.Fatalf("want 204, got %d", r.StatusCode)
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	userID := createUser(t, srv)
	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/users/"+userID+"/todos/bad", nil)
	r, _ := http.DefaultClient.Do(req)
	defer r.Body.Close()
	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404, got %d", r.StatusCode)
	}
}
