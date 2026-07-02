package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/todo/internal/store"
)

func newTestTodosHandler() (*TodosHandler, *UsersHandler, *store.MemoryUserStore, *store.MemoryTodoStore) {
	us := store.NewMemoryUserStore()
	ts := store.NewMemoryTodoStore()
	uh := NewUsersHandler(us)
	th := NewTodosHandler(us, ts)
	return th, uh, us, ts
}

func createTodoViaHandler(t *testing.T, h *TodosHandler, userID, title string, priority int) string {
	t.Helper()
	body := fmt.Sprintf(`{"title":%q,"priority":%d}`, title, priority)
	path := "/users/" + userID + "/todos"
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.URL.Path = path
	w := httptest.NewRecorder()
	h.CreateTodo(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("createTodoViaHandler: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	return resp["id"].(string)
}

func TestCreateTodo(t *testing.T) {
	th, uh, _, _ := newTestTodosHandler()
	uid := createUserViaHandler(t, uh, "alice", "alice@example.com")

	body := `{"title":"Buy milk","priority":1}`
	path := "/users/" + uid + "/todos"
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.URL.Path = path
	w := httptest.NewRecorder()
	th.CreateTodo(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("expected id field in response")
	}
	if resp["status"] != "Pending" {
		t.Errorf("expected status Pending, got %v", resp["status"])
	}
}

func TestListTodos(t *testing.T) {
	th, uh, _, _ := newTestTodosHandler()
	uid := createUserViaHandler(t, uh, "alice", "alice@example.com")
	createTodoViaHandler(t, th, uid, "Task 1", 0)
	createTodoViaHandler(t, th, uid, "Task 2", 1)

	path := "/users/" + uid + "/todos"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.URL.Path = path
	w := httptest.NewRecorder()
	th.ListTodos(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp []any
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 2 {
		t.Errorf("expected 2 todos, got %d", len(resp))
	}
}

func TestGetTodo(t *testing.T) {
	th, uh, _, _ := newTestTodosHandler()
	uid := createUserViaHandler(t, uh, "alice", "alice@example.com")
	tid := createTodoViaHandler(t, th, uid, "Buy milk", 1)

	path := "/users/" + uid + "/todos/" + tid
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.URL.Path = path
	w := httptest.NewRecorder()
	th.GetTodo(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["title"] != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %v", resp["title"])
	}
}

func TestUpdateTodo(t *testing.T) {
	th, uh, _, _ := newTestTodosHandler()
	uid := createUserViaHandler(t, uh, "alice", "alice@example.com")
	tid := createTodoViaHandler(t, th, uid, "Buy milk", 1)

	body := `{"title":"Buy milk 2","priority":2,"status":2}`
	path := "/users/" + uid + "/todos/" + tid
	req := httptest.NewRequest(http.MethodPut, path, bytes.NewBufferString(body))
	req.URL.Path = path
	w := httptest.NewRecorder()
	th.UpdateTodo(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["title"] != "Buy milk 2" {
		t.Errorf("expected title 'Buy milk 2', got %v", resp["title"])
	}
}

func TestDeleteTodo(t *testing.T) {
	th, uh, _, _ := newTestTodosHandler()
	uid := createUserViaHandler(t, uh, "alice", "alice@example.com")
	tid := createTodoViaHandler(t, th, uid, "Buy milk", 1)

	// Delete.
	path := "/users/" + uid + "/todos/" + tid
	dReq := httptest.NewRequest(http.MethodDelete, path, nil)
	dReq.URL.Path = path
	dW := httptest.NewRecorder()
	th.DeleteTodo(dW, dReq)

	if dW.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", dW.Code, dW.Body.String())
	}

	// GET after delete should 404.
	gReq := httptest.NewRequest(http.MethodGet, path, nil)
	gReq.URL.Path = path
	gW := httptest.NewRecorder()
	th.GetTodo(gW, gReq)

	if gW.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", gW.Code)
	}
}
