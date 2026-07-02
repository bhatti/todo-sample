// Package handler provides HTTP handlers for the todo REST API.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
)

// UserHandler handles HTTP requests for user resources.
type UserHandler struct {
	users store.UserStore
}

// NewUserHandler creates a UserHandler backed by the given UserStore.
func NewUserHandler(users store.UserStore) *UserHandler {
	return &UserHandler{users: users}
}

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Create handles POST /users.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	now := time.Now().UTC()
	u := model.User{
		ID:        ulid.Make().String(),
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.users.CreateUser(u); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(u))
}

// Get handles GET /users/{userID}.
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userID")
	u, err := h.users.GetUser(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(u))
}

// List handles GET /users.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.ListUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toUserResponses(users))
}

// Delete handles DELETE /users/{userID}.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userID")
	if err := h.users.DeleteUser(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
