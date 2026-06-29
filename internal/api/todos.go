package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
)

// TodoHandler handles HTTP requests for todo resources scoped under a user.
type TodoHandler struct {
	userStore store.UserStore
	todoStore store.TodoStore
}

// NewTodoHandler creates a TodoHandler with the given user and todo stores.
func NewTodoHandler(us store.UserStore, ts store.TodoStore) *TodoHandler {
	return &TodoHandler{userStore: us, todoStore: ts}
}

type todoRequest struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Priority    model.Priority `json:"priority"`
	Status      model.Status   `json:"status"`
	DueAt       *time.Time    `json:"due_at,omitempty"`
}

// Create handles POST /users/{userID}/todos — validates the parent user exists, then creates a new todo.
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if _, err := h.userStore.GetByID(r.Context(), userID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	var req todoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	now := time.Now().UTC()
	t := &model.Todo{
		ID:          ulid.Make().String(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		DueAt:       req.DueAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := h.todoStore.Create(r.Context(), t); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

// List handles GET /users/{userID}/todos — returns all todos for the given user, never null.
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if _, err := h.userStore.GetByID(r.Context(), userID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	todos, err := h.todoStore.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if todos == nil {
		todos = []*model.Todo{}
	}
	writeJSON(w, http.StatusOK, todos)
}

// GetByID handles GET /users/{userID}/todos/{id} — returns 404 if the todo exists but belongs to a different user.
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id := chi.URLParam(r, "id")
	t, err := h.todoStore.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if t.UserID != userID {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// Update handles PUT /users/{userID}/todos/{id} — replaces all mutable fields; returns 404 if the todo belongs to a different user.
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id := chi.URLParam(r, "id")
	if _, err := h.userStore.GetByID(r.Context(), userID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	existing, err := h.todoStore.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if existing.UserID != userID {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	var req todoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	existing.Title = req.Title
	existing.Description = req.Description
	existing.Priority = req.Priority
	existing.Status = req.Status
	existing.DueAt = req.DueAt
	existing.UpdatedAt = time.Now().UTC()
	if err := h.todoStore.Update(r.Context(), existing); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, existing)
}

// Delete handles DELETE /users/{userID}/todos/{id} — returns 404 if the todo belongs to a different user, 204 on success.
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id := chi.URLParam(r, "id")
	existing, err := h.todoStore.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if existing.UserID != userID {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	if err := h.todoStore.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
