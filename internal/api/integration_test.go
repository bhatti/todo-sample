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

func newTestServer() *httptest.Server {
	us := store.NewMemoryUserStore()
	ts := store.NewMemoryTodoStore()
	return httptest.NewServer(NewRouter(us, ts))
}

func do(t *testing.T, srv *httptest.Server, method, path, body string) (int, map[string]any) {
	t.Helper()
	var req *http.Request
	var err error
	if body != "" {
		req, err = http.NewRequest(method, srv.URL+path, bytes.NewBufferString(body))
	} else {
		req, err = http.NewRequest(method, srv.URL+path, nil)
	}
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()
	var m map[string]any
	json.NewDecoder(resp.Body).Decode(&m)
	return resp.StatusCode, m
}

func doList(t *testing.T, srv *httptest.Server, path string) (int, []any) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, srv.URL+path, nil)
	if err != nil {
		t.Fatalf("doList: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("doList: %v", err)
	}
	defer resp.Body.Close()
	var list []any
	json.NewDecoder(resp.Body).Decode(&list)
	return resp.StatusCode, list
}

func TestIntegration_Users(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	// POST /users
	code, body := do(t, srv, http.MethodPost, "/users", `{"username":"alice","email":"alice@example.com"}`)
	if code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %v", code, body)
	}
	id := body["id"].(string)

	// GET /users/{id}
	code, body = do(t, srv, http.MethodGet, "/users/"+id, "")
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}
	if body["username"] != "alice" {
		t.Errorf("expected alice, got %v", body["username"])
	}

	// GET /users
	listCode, list := doList(t, srv, "/users")
	if listCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", listCode)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 user, got %d", len(list))
	}

	// PUT /users/{id}
	code, body = do(t, srv, http.MethodPut, "/users/"+id, `{"username":"alice2","email":"alice2@example.com"}`)
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %v", code, body)
	}
	if body["username"] != "alice2" {
		t.Errorf("expected alice2, got %v", body["username"])
	}

	// DELETE /users/{id}
	delReq, _ := http.NewRequest(http.MethodDelete, srv.URL+"/users/"+id, nil)
	delResp, _ := http.DefaultClient.Do(delReq)
	delResp.Body.Close()
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}

	// GET /users/{id} after delete — should 404.
	code, _ = do(t, srv, http.MethodGet, "/users/"+id, "")
	if code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", code)
	}
}

func TestIntegration_Todos(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	// Create a user first.
	code, body := do(t, srv, http.MethodPost, "/users", `{"username":"bob","email":"bob@example.com"}`)
	if code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %v", code, body)
	}
	uid := body["id"].(string)

	basePath := fmt.Sprintf("/users/%s/todos", uid)

	// POST /users/{id}/todos
	code, body = do(t, srv, http.MethodPost, basePath, `{"title":"Buy milk","priority":1}`)
	if code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %v", code, body)
	}
	tid := body["id"].(string)
	if body["status"] != "Pending" {
		t.Errorf("expected Pending status, got %v", body["status"])
	}

	// GET /users/{id}/todos
	listCode, list := doList(t, srv, basePath)
	if listCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", listCode)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 todo, got %d", len(list))
	}

	// GET /users/{id}/todos/{todoID}
	code, body = do(t, srv, http.MethodGet, basePath+"/"+tid, "")
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %v", code, body)
	}
	if body["title"] != "Buy milk" {
		t.Errorf("expected 'Buy milk', got %v", body["title"])
	}

	// PUT /users/{id}/todos/{todoID}
	code, body = do(t, srv, http.MethodPut, basePath+"/"+tid, `{"title":"Buy milk 2","priority":2,"status":2}`)
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %v", code, body)
	}
	if body["title"] != "Buy milk 2" {
		t.Errorf("expected 'Buy milk 2', got %v", body["title"])
	}

	// DELETE /users/{id}/todos/{todoID}
	delReq, _ := http.NewRequest(http.MethodDelete, srv.URL+basePath+"/"+tid, nil)
	delResp, _ := http.DefaultClient.Do(delReq)
	delResp.Body.Close()
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}

	// GET /users/{id}/todos/{todoID} after delete — should 404.
	code, _ = do(t, srv, http.MethodGet, basePath+"/"+tid, "")
	if code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", code)
	}
}
