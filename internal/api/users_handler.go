package api

import (
	"net/http"
	"strings"

	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
)

// UsersHandler handles HTTP requests for user resources.
type UsersHandler struct {
	users store.UserStore
}

// NewUsersHandler returns a UsersHandler backed by the given store.
func NewUsersHandler(users store.UserStore) *UsersHandler {
	return &UsersHandler{users: users}
}

// createUserRequest is the request body for POST /users.
type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// updateUserRequest is the request body for PUT /users/{id}.
type updateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// userResponse is the JSON representation of a user.
type userResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toUserResponse(u *model.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// CreateUser handles POST /users.
func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Username == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}

	created, err := h.users.Create(&model.User{Username: req.Username, Email: req.Email})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toUserResponse(created))
}

// ListUsers handles GET /users.
func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]userResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserResponse(u))
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetUser handles GET /users/{id}.
func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/users/")
	u, err := h.users.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(u))
}

// UpdateUser handles PUT /users/{id}.
func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/users/")
	var req updateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Username == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}
	updated, err := h.users.Update(&model.User{ID: id, Username: req.Username, Email: req.Email})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(updated))
}

// DeleteUser handles DELETE /users/{id}.
func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/users/")
	if err := h.users.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// extractID parses a path segment after the given prefix.
// e.g. extractID("/users/abc123", "/users/") -> "abc123"
func extractID(path, prefix string) string {
	trimmed := strings.TrimPrefix(path, prefix)
	// Remove any trailing path segments.
	if idx := strings.Index(trimmed, "/"); idx != -1 {
		return trimmed[:idx]
	}
	return trimmed
}
