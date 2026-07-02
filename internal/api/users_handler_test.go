package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/todo/internal/store"
)

func newTestUsersHandler() (*UsersHandler, *store.MemoryUserStore) {
	s := store.NewMemoryUserStore()
	return NewUsersHandler(s), s
}

func createUserViaHandler(t *testing.T, h *UsersHandler, username, email string) string {
	t.Helper()
	body := `{"username":"` + username + `","email":"` + email + `"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.CreateUser(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("createUserViaHandler: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	return resp["id"].(string)
}

func TestCreateUser(t *testing.T) {
	h, _ := newTestUsersHandler()
	body := `{"username":"alice","email":"alice@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["id"] == "" || resp["id"] == nil {
		t.Error("expected id field in response")
	}
}

func TestCreateUser_MissingUsername(t *testing.T) {
	h, _ := newTestUsersHandler()
	body := `{"username":"","email":"alice@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] == nil {
		t.Error("expected error field in response")
	}
}

func TestGetUser(t *testing.T) {
	h, _ := newTestUsersHandler()
	id := createUserViaHandler(t, h, "alice", "alice@example.com")

	req := httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
	req.URL.Path = "/users/" + id
	w := httptest.NewRecorder()
	h.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "alice" {
		t.Errorf("expected username alice, got %v", resp["username"])
	}
}

func TestGetUser_NotFound(t *testing.T) {
	h, _ := newTestUsersHandler()

	req := httptest.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	req.URL.Path = "/users/nonexistent"
	w := httptest.NewRecorder()
	h.GetUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	h, _ := newTestUsersHandler()
	id := createUserViaHandler(t, h, "alice", "alice@example.com")

	uBody := `{"username":"alice2","email":"alice2@example.com"}`
	uReq := httptest.NewRequest(http.MethodPut, "/users/"+id, bytes.NewBufferString(uBody))
	uReq.URL.Path = "/users/" + id
	uW := httptest.NewRecorder()
	h.UpdateUser(uW, uReq)

	if uW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", uW.Code, uW.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(uW.Body).Decode(&resp)
	if resp["username"] != "alice2" {
		t.Errorf("expected username alice2, got %v", resp["username"])
	}
}

func TestDeleteUser(t *testing.T) {
	h, _ := newTestUsersHandler()
	id := createUserViaHandler(t, h, "alice", "alice@example.com")

	dReq := httptest.NewRequest(http.MethodDelete, "/users/"+id, nil)
	dReq.URL.Path = "/users/" + id
	dW := httptest.NewRecorder()
	h.DeleteUser(dW, dReq)

	if dW.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", dW.Code, dW.Body.String())
	}

	// GET after delete should 404.
	gReq := httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
	gReq.URL.Path = "/users/" + id
	gW := httptest.NewRecorder()
	h.GetUser(gW, gReq)

	if gW.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", gW.Code)
	}
}
