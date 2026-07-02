package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/todo/internal/handler"
	"github.com/user/todo/internal/router"
	"github.com/user/todo/internal/store"
)

// setup returns an HTTP handler wired with fresh in-memory stores.
func setup() http.Handler {
	us := store.NewMemoryUserStore()
	ts := store.NewMemoryTodoStore()
	uh := handler.NewUserHandler(us)
	th := handler.NewTodoHandler(us, ts)
	return router.New(uh, th)
}

func body(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&m); err != nil {
		t.Fatalf("decode body: %v (body: %s)", err, rr.Body.String())
	}
	return m
}

func bodySlice(t *testing.T, rr *httptest.ResponseRecorder) []any {
	t.Helper()
	var s []any
	if err := json.NewDecoder(rr.Body).Decode(&s); err != nil {
		t.Fatalf("decode slice body: %v (body: %s)", err, rr.Body.String())
	}
	return s
}

// ---- User handler tests ----

func TestCreateUser_OK(t *testing.T) {
	h := setup()
	payload := `{"username":"alice","email":"alice@example.com"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d (body: %s)", rr.Code, rr.Body.String())
	}
	m := body(t, rr)
	if m["id"] == "" || m["id"] == nil {
		t.Error("expected non-empty id")
	}
	if m["username"] != "alice" {
		t.Errorf("expected username alice, got %v", m["username"])
	}
}

func TestCreateUser_BadRequest(t *testing.T) {
	h := setup()
	payload := `{"username":"","email":"alice@example.com"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
	m := body(t, rr)
	if _, ok := m["error"]; !ok {
		t.Error("expected error key in response")
	}
}

func TestGetUser_OK(t *testing.T) {
	h := setup()
	// Create user first.
	payload := `{"username":"alice","email":"alice@example.com"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	h.ServeHTTP(rr, req)
	created := body(t, rr)
	id := created["id"].(string)

	// Now get it.
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
	h.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d (body: %s)", rr2.Code, rr2.Body.String())
	}
	m := body(t, rr2)
	if m["id"] != id {
		t.Errorf("expected id %s, got %v", id, m["id"])
	}
}

func TestGetUser_NotFound(t *testing.T) {
	h := setup()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestListUsers_OK(t *testing.T) {
	h := setup()
	for _, p := range []string{
		`{"username":"alice","email":"alice@example.com"}`,
		`{"username":"bob","email":"bob@example.com"}`,
	} {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(p))
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("create user failed: %d", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	s := bodySlice(t, rr)
	if len(s) != 2 {
		t.Errorf("expected 2 users, got %d", len(s))
	}
}

func TestUpdateUser_OK(t *testing.T) {
	h := setup()
	payload := `{"username":"alice","email":"alice@example.com"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	h.ServeHTTP(rr, req)
	created := body(t, rr)
	id := created["id"].(string)

	update := `{"username":"bob","email":"bob@example.com"}`
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPut, "/users/"+id, strings.NewReader(update))
	h.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d (body: %s)", rr2.Code, rr2.Body.String())
	}
	m := body(t, rr2)
	if m["username"] != "bob" {
		t.Errorf("expected username bob, got %v", m["username"])
	}
}

func TestDeleteUser_OK(t *testing.T) {
	h := setup()
	payload := `{"username":"alice","email":"alice@example.com"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	h.ServeHTTP(rr, req)
	created := body(t, rr)
	id := created["id"].(string)

	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodDelete, "/users/"+id, nil)
	h.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rr2.Code)
	}
}

// ---- Todo handler tests ----

// createTestUser creates a user and returns their ID.
func createTestUser(t *testing.T, h http.Handler) string {
	t.Helper()
	payload := `{"username":"alice","email":"alice@example.com"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create user failed: %d %s", rr.Code, rr.Body.String())
	}
	m := body(t, rr)
	return m["id"].(string)
}

func TestCreateTodo_OK(t *testing.T) {
	h := setup()
	uid := createTestUser(t, h)

	payload := `{"title":"Buy milk","priority":"low"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s/todos", uid), strings.NewReader(payload))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d (body: %s)", rr.Code, rr.Body.String())
	}
	m := body(t, rr)
	if m["id"] == nil || m["id"] == "" {
		t.Error("expected non-empty id")
	}
	if m["user_id"] != uid {
		t.Errorf("expected user_id %s, got %v", uid, m["user_id"])
	}
	if m["status"] != "pending" {
		t.Errorf("expected status pending, got %v", m["status"])
	}
}

func TestCreateTodo_BadRequest(t *testing.T) {
	h := setup()
	uid := createTestUser(t, h)

	payload := `{"title":""}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s/todos", uid), strings.NewReader(payload))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
	m := body(t, rr)
	if _, ok := m["error"]; !ok {
		t.Error("expected error key")
	}
}

func TestCreateTodo_UserNotFound(t *testing.T) {
	h := setup()
	payload := `{"title":"Buy milk","priority":"low"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/unknown/todos", strings.NewReader(payload))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestListTodos_OK(t *testing.T) {
	h := setup()
	uid := createTestUser(t, h)

	for _, p := range []string{
		`{"title":"Task 1"}`,
		`{"title":"Task 2"}`,
	} {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s/todos", uid), strings.NewReader(p))
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("create todo failed: %d %s", rr.Code, rr.Body.String())
		}
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s/todos", uid), nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	s := bodySlice(t, rr)
	if len(s) != 2 {
		t.Errorf("expected 2 todos, got %d", len(s))
	}
}

func TestGetTodo_OK(t *testing.T) {
	h := setup()
	uid := createTestUser(t, h)

	payload := `{"title":"Buy milk"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s/todos", uid), strings.NewReader(payload))
	h.ServeHTTP(rr, req)
	created := body(t, rr)
	tid := created["id"].(string)

	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s/todos/%s", uid, tid), nil)
	h.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d (body: %s)", rr2.Code, rr2.Body.String())
	}
	m := body(t, rr2)
	if m["id"] != tid {
		t.Errorf("expected id %s, got %v", tid, m["id"])
	}
}

func TestGetTodo_NotFound(t *testing.T) {
	h := setup()
	uid := createTestUser(t, h)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s/todos/bad", uid), nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUpdateTodo_Status(t *testing.T) {
	h := setup()
	uid := createTestUser(t, h)

	payload := `{"title":"Buy milk"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s/todos", uid), strings.NewReader(payload))
	h.ServeHTTP(rr, req)
	created := body(t, rr)
	tid := created["id"].(string)

	update := `{"status":"done"}`
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s/todos/%s", uid, tid), strings.NewReader(update))
	h.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d (body: %s)", rr2.Code, rr2.Body.String())
	}
	m := body(t, rr2)
	if m["status"] != "done" {
		t.Errorf("expected status done, got %v", m["status"])
	}
}

func TestDeleteTodo_OK(t *testing.T) {
	h := setup()
	uid := createTestUser(t, h)

	payload := `{"title":"Buy milk"}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s/todos", uid), strings.NewReader(payload))
	h.ServeHTTP(rr, req)
	created := body(t, rr)
	tid := created["id"].(string)

	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s/todos/%s", uid, tid), nil)
	h.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rr2.Code)
	}
}

