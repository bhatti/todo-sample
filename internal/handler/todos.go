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

// TodoHandler handles HTTP requests for todo resources scoped to a user.
type TodoHandler struct {
	users store.UserStore
	todos store.TodoStore
}

// NewTodoHandler creates a TodoHandler backed by the given stores.
func NewTodoHandler(users store.UserStore, todos store.TodoStore) *TodoHandler {
	return &TodoHandler{users: users, todos: todos}
}

// CreateTodoRequest is the request body for creating a todo.
type CreateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	DueAt       *time.Time `json:"due_at,omitempty"`
}

// UpdateTodoRequest is the request body for updating a todo.
type UpdateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	Status      int        `json:"status"`
	DueAt       *time.Time `json:"due_at,omitempty"`
}

// Create handles POST /users/{userID}/todos.
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if _, err := h.users.GetUser(userID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	now := time.Now().UTC()
	t := model.Todo{
		ID:          ulid.Make().String(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    model.Priority(req.Priority),
		Status:      model.StatusPending,
		DueAt:       req.DueAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.todos.CreateTodo(t); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toTodoResponse(t))
}

// List handles GET /users/{userID}/todos.
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if _, err := h.users.GetUser(userID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	todos, err := h.todos.ListTodos(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toTodoResponses(todos))
}

// Get handles GET /users/{userID}/todos/{todoID}.
func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	todoID := r.PathValue("todoID")

	t, err := h.todos.GetTodo(userID, todoID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toTodoResponse(t))
}

// Update handles PUT /users/{userID}/todos/{todoID}.
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	todoID := r.PathValue("todoID")

	existing, err := h.todos.GetTodo(userID, todoID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	existing.Title = req.Title
	existing.Description = req.Description
	existing.Priority = model.Priority(req.Priority)
	existing.Status = model.Status(req.Status)
	existing.DueAt = req.DueAt
	existing.UpdatedAt = time.Now().UTC()

	if err := h.todos.UpdateTodo(existing); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toTodoResponse(existing))
}

// Delete handles DELETE /users/{userID}/todos/{todoID}.
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	todoID := r.PathValue("todoID")

	if err := h.todos.DeleteTodo(userID, todoID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
